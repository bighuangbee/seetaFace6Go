package common

import (
	"face_recognize/recognize/face_rec"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"log"
	"path/filepath"
	"seetaFace6go"
	"time"
	"video-find-face/seetaFace"
)

type Face struct {
	Seeta       *seetaFace.SeetaFace
	FaceFeature face_rec.IFaceFeature
	TargetRect  image.Rectangle

	VideoName string

	frames chan *Frame

	EmptyCount int
	Tracking   bool

	bestImage *Frame
}

type Frame struct {
	Mat   *gocv.Mat
	Count int
	Score float32
}

func (frame *Frame) ToSeetaImage(targetRect image.Rectangle) (seetaImg *seetaFace6go.SeetaImageData) {
	var frameRegion = *frame.Mat
	if !targetRect.Empty() {
		frameRegion = frame.Mat.Region(targetRect)
		defer frameRegion.Close()
	}

	img, _ := frameRegion.ToImage()
	return seetaFace6go.NewSeetaImageDataFromImage(img)

	fmt.Println(frameRegion.Cols(), frameRegion.Rows(), frameRegion.Channels())
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
	sFace.Detector.SetProperty(seetaFace6go.FaceDetector_PROPERTY_NUMBER_THREADS, 1)

	return &Face{
		FaceFeature: FaceFeature,
		Seeta:       sFace,
		TargetRect:  targetRect,
		frames:      make(chan *Frame),
		//bestImages:  make(map[int]ImageScore),
	}
}

func (face *Face) SetFrame(frame *Frame) {
	mat := gocv.NewMat()
	frame.Mat.CopyTo(&mat)
	face.frames <- &Frame{
		Mat:   &mat,
		Count: frame.Count,
	}
}

func (face *Face) ProcessFrame() {
	for frame := range face.frames {
		face.Detect(frame)
		frame.Mat.Close()
	}
}

func (face *Face) NewTracker(width, height int) {
	if face.Seeta.Tracker == nil {
		if !face.TargetRect.Empty() {
			width = face.TargetRect.Size().X
			height = face.TargetRect.Size().Y
		}
		face.Seeta.Tracker = seetaFace6go.NewFaceTracker(width, height)
		face.Seeta.Tracker.SetVideoStable(true)
		face.Seeta.Tracker.SetInterval(1)
		//sFace.Tracker.SetThreads(4)
		face.Seeta.Tracker.SetMinFaceSize(60)
	}
}

func (face *Face) RecognizeFrame(frame *gocv.Mat, frameCount int, pids []int) {
	t := time.Now()
	mat := gocv.NewMat()
	frame.CopyTo(&mat)
	fe, err := face.Recognize(mat)
	if err != nil {
		log.Println("Recognize error", err)
	}
	defer mat.Close()

	log.Println("Recognize faceLen:", len(fe), "time.Sinde:", time.Since(t).Milliseconds(), "pids", pids)

	needSave := false
	rects := []image.Rectangle{}
	for _, entity := range fe {
		log.Println("Recognize entity", entity.Quality, entity.Rect)
		rects = append(rects, entity.Rect)

		if entity.Quality > 0.6 && !needSave {
			needSave = true
		}
	}

	if needSave {
		picBaseName := filepath.Base(filepath.Base(face.VideoName))
		ok := gocv.IMWrite(filepath.Join(Output, fmt.Sprintf("%s_frame_%d.jpg", picBaseName, frameCount)), mat)
		if !ok {
			log.Println("Write image error")
		}

		for i, face := range fe {
			faceRegion := mat.Region(face.Rect)
			gocv.IMWrite(filepath.Join(Output, fmt.Sprintf("%s_frame_%d_pid_%d_q_%0.2f.jpg", picBaseName, frameCount, i, face.Quality)), faceRegion)
			faceRegion.Close()
		}
	}

}

func (face *Face) Detect(frame *Frame) (infos []*seetaFace.DetectInfo) {
	start := time.Now()
	img := frame.ToSeetaImage(face.TargetRect)
	faces := face.Seeta.Detector.Detect(img)

	//if frameCount%15 == 0 {
	log.Println("faceTrack, frame:", frame.Count, "time:", time.Now().Sub(start).Milliseconds())
	//}

	if len(faces) > 0 {
		pids := []int{}
		for _, info := range faces {
			pointInfo := face.Seeta.Landmarker.Mark(img, info.Postion)
			brightness := face.Seeta.QualityCheck.CheckBrightness(img, info.Postion, pointInfo)
			clarity := face.Seeta.QualityCheck.CheckClarity(img, info.Postion, pointInfo)

			infos = append(infos, &seetaFace.DetectInfo{
				Confidence: info.Score,
				Clarity:    clarity.Score,
				Brightness: brightness.Score,
				FaceInfo:   info,
			})
		}

		if face.FaceFeature != nil {
			go face.RecognizeFrame(frame.Mat, frame.Count, pids)
		}
	}

	return infos
}

func (face *Face) Recognize(frame gocv.Mat) ([]*face_rec.FaceEntity, error) {
	buf, err := gocv.IMEncode(".jpg", frame)
	if err != nil {
		return nil, err
	}

	image, _ := face_rec.ReadImageByFormByte(buf.GetBytes(), "1.jpg")
	faces, err := face.FaceFeature.ExtractFeature(image, face_rec.ExtractAll)
	if err != nil {
		return nil, err
	}
	buf.Close()

	return faces, nil
}
