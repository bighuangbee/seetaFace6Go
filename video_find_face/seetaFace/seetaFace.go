package seetaFace

import (
	"gocv.io/x/gocv"
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

	targetRect image.Rectangle
}

type DetectInfo struct {
	Confidence float32
	Clarity    float32
	Brightness float32
	Integrity  float32
	FaceInfo   *seetaFace6go.SeetaFaceInfo
}

const ThreadsCount = 1

func NewSeetaFace(modelPath string, targetRect image.Rectangle) *SeetaFace {
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
		targetRect:   targetRect,
	}
}

func (face *SeetaFace) NewTracker(width, height int) {
	if face.Tracker == nil {
		log.Println("NewTracker", width, height, face.targetRect)

		if !face.targetRect.Empty() {
			width = face.targetRect.Size().X
			height = face.targetRect.Size().Y
		}
		face.Tracker = seetaFace6go.NewFaceTracker(width, height)
		face.Tracker.SetVideoStable(true)
		face.Tracker.SetInterval(1)
		face.Tracker.SetThreads(ThreadsCount) //mac: 4
		face.Tracker.SetMinFaceSize(60)
		face.Tracker.SetThreshold(0.3)
	}
}

func (face *SeetaFace) ResetTracker() {
	face.Tracker = nil
}

func (face *SeetaFace) Detect(img *seetaFace6go.SeetaImageData) (infos []*DetectInfo) {
	faces := face.Detector.Detect(img)

	if len(faces) > 0 {
		for _, info := range faces {
			pointInfo := face.Landmarker.Mark(img, info.Postion)
			brightness := face.QualityCheck.CheckBrightness(img, info.Postion, pointInfo)
			clarity := face.QualityCheck.CheckClarity(img, info.Postion, pointInfo)
			integrity := face.QualityCheck.CheckIntegrity(img, info.Postion, pointInfo)

			//ok, _ := face.Seeta.Recognizer.Extract(img, pointInfo)

			infos = append(infos, &DetectInfo{
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

func ToSeetaImage(mat gocv.Mat, targetRect image.Rectangle) (seetaImg *seetaFace6go.SeetaImageData) {
	var frameRegion = mat
	if !targetRect.Empty() {
		frameRegion = mat.Region(targetRect)
		//defer frameRegion.Close()
	}

	//img, _ := frameRegion.ToImage()
	//return seetaFace6go.NewSeetaImageDataFromImage(img)

	imageData := seetaFace6go.NewSeetaImageData(frameRegion.Cols(), frameRegion.Rows(), frameRegion.Channels())
	imageData.SetUint8(frameRegion.ToBytes())
	return imageData
}
