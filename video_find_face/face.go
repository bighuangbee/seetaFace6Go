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
	CountStart float64
	Score      float32
}

type TrackState struct {
	//连续多少帧没检测到人脸
	EmptyCount int
	Tracking   bool
}

func (frame *Frame) ToSeetaImage(targetRect image.Rectangle) (seetaImg *seetaFace6go.SeetaImageData) {
	return seetaFace.ToSeetaImage(*frame.Mat, targetRect)
}

var Output = "./output"

func NewFace(sFaceModel string, targetRect image.Rectangle) *Face {
	var FaceFeature face_rec.IFaceFeature
	//var err error
	//FaceFeature, err = faceRec.New("/root/face_recognize/Recognize/libs/face_gpu/models")
	//if err != nil {
	//	log.Fatal(err)
	//}

	sFace := seetaFace.NewSeetaFace(sFaceModel, targetRect)
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
