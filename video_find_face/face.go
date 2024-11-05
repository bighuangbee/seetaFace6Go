package video_find_face

import (
	"face_recognize/recognize/face_rec"
	"gocv.io/x/gocv"
	"image"
	"seetaFace6go"
	"video-find-face/seetaFace"
)

type Face struct {
	Seeta       *seetaFace.SeetaFace
	FaceFeature face_rec.IFaceFeature
	TargetRect  image.Rectangle

	VideoInfo *VideoInfo

	TrackState TrackState

	frames chan *Frame

	bestImage *Frame
}

type Frame struct {
	Mat        *gocv.Mat
	Count      int //帧计数
	CountStart int
	Score      float32
}

type TrackState struct {
	//连续多少帧没检测到人脸
	EmptyCount int
	Tracking   bool
}

func (frame *Frame) ToSeetaImage(targetRect image.Rectangle) (seetaImg *seetaFace6go.SeetaImageData) {
	var frameRegion = *frame.Mat
	if !targetRect.Empty() {
		frameRegion = frame.Mat.Region(targetRect)
		//defer frameRegion.Close()
	}

	//img, _ := frameRegion.ToImage()
	//return seetaFace6go.NewSeetaImageDataFromImage(img)

	imageData := seetaFace6go.NewSeetaImageData(frameRegion.Cols(), frameRegion.Rows(), frameRegion.Channels())
	imageData.SetUint8(frameRegion.ToBytes())
	return imageData
}

var Output = "./output"

func NewFace(sFaceModel string, targetRect image.Rectangle) *Face {
	var FaceFeature face_rec.IFaceFeature
	//var err error
	//FaceFeature, err = faceRec.New("/root/face_recognize/Recognize/libs/face_gpu/models")
	//if err != nil {
	//	log.Fatal(err)
	//}

	sFace := seetaFace.NewSeetaFace(sFaceModel)
	sFace.Detector.SetProperty(seetaFace6go.FaceDetector_PROPERTY_MIN_FACE_SIZE, 60)
	sFace.Detector.SetProperty(seetaFace6go.FaceDetector_PROPERTY_NUMBER_THREADS, 4)

	face := &Face{
		FaceFeature: FaceFeature,
		Seeta:       sFace,
		TargetRect:  targetRect,
		frames:      make(chan *Frame, 10),
	}

	return face
}

func (face *Face) Detect(frame *Frame) (infos []*seetaFace.DetectInfo) {
	img := frame.ToSeetaImage(face.TargetRect)
	faces := face.Seeta.Detector.Detect(img)

	if len(faces) > 0 {
		for _, info := range faces {
			pointInfo := face.Seeta.Landmarker.Mark(img, info.Postion)
			brightness := face.Seeta.QualityCheck.CheckBrightness(img, info.Postion, pointInfo)
			clarity := face.Seeta.QualityCheck.CheckClarity(img, info.Postion, pointInfo)
			integrity := face.Seeta.QualityCheck.CheckIntegrity(img, info.Postion, pointInfo)

			//ok, _ := face.Seeta.Recognizer.Extract(img, pointInfo)

			infos = append(infos, &seetaFace.DetectInfo{
				Confidence: info.Score,
				Clarity:    clarity.Score,
				Brightness: brightness.Score,
				Integrity:  integrity.Score,
				FaceInfo:   info,
			})
		}

		//if face.FaceFeature != nil {
		//	go face.RecognizeFrame(frame.Mat, frame.Count, pids)
		//}
	}

	return infos
}
