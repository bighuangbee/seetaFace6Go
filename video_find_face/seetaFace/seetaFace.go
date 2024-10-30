package seetaFace

import (
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
	Quality  *seetaFace6go.QualityResult
	Clarity  *seetaFace6go.QualityResult
	FaceInfo *seetaFace6go.SeetaFaceInfo
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
