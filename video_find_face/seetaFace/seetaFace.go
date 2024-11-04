package seetaFace

import (
	"image"
	"log"
	"seetaFace6go"
)

type SeetaFace struct {
	Detector     *seetaFace6go.FaceDetector
	Landmarker   *seetaFace6go.FaceLandmarker
	Recognizer   *seetaFace6go.FaceRecognizer
	QualityCheck *seetaFace6go.QualityCheck
	Tracker      *seetaFace6go.FaceTracker

	modelPath     string
	Width, Height int
}

type DetectInfo struct {
	Confidence float32
	Clarity    float32
	Brightness float32
	Integrity  float32
	FaceInfo   *seetaFace6go.SeetaFaceInfo
}

func NewSeetaFace(modelPath string) *SeetaFace {
	seetaFace6go.InitModelPath(modelPath)

	// 人脸检测器
	fd := seetaFace6go.NewFaceDetector()

	// 人脸特征定位器
	// 使用5点信息模型
	fl := seetaFace6go.NewFaceLandmarker(seetaFace6go.ModelType_light)

	// 人脸特征提取器
	fr := seetaFace6go.NewFaceRecognizer(seetaFace6go.ModelType_light)

	// 质量评估器
	qr := seetaFace6go.NewQualityCheck()

	return &SeetaFace{
		Detector:     fd,
		Landmarker:   fl,
		Recognizer:   fr,
		QualityCheck: qr,
	}
}

func (face *SeetaFace) NewTracker(width, height int, targetRect image.Rectangle) {
	if face.Tracker == nil {
		log.Println("NewTracker", width, height, targetRect)

		if !targetRect.Empty() {
			width = targetRect.Size().X
			height = targetRect.Size().Y
		}
		face.Tracker = seetaFace6go.NewFaceTracker(width, height)
		face.Tracker.SetVideoStable(true)
		face.Tracker.SetInterval(1)
		face.Tracker.SetThreads(1) //mac: 4
		face.Tracker.SetMinFaceSize(60)
	}
}

func (face *SeetaFace) ResetTracker() {
	face.Tracker = nil
}
