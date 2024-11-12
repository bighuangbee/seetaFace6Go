package video_find_face

import (
	"face_recognize/recognize/face_rec"
	"gocv.io/x/gocv"
	"image"
	"log"
	"seetaFace6go"
	"sync"
	"video-find-face/rec_gpu"
	"video-find-face/seetaFace"
)

type Face struct {
	Seeta        *seetaFace.SeetaFace
	RecognizeGpu face_rec.IFaceFeature
	TargetRect   image.Rectangle

	VideoInfo *VideoInfo

	TrackState TrackState

	//已跟踪到人脸的视频帧
	trackedBuffer chan *Frame

	//最佳图像
	bestImage *Frame

	//视频截取
	VideoWriter   VideoWriter
	muVideoWriter sync.RWMutex

	//视频帧缓存，用于视频截取
	FrameBuffer []*Frame
}

type Frame struct {
	Mat        *gocv.Mat
	Count      int //帧计数
	CountStart float64
	Score      float32
}

type TrackState struct {
	//连续多少帧没检测到人脸
	EmptyCount    int
	MaxEmptyCount int
	Tracking      bool
}

func (frame *Frame) ToSeetaImage(targetRect image.Rectangle) (seetaImg *seetaFace6go.SeetaImageData) {
	return seetaFace.ToSeetaImage(*frame.Mat, targetRect)
}

func NewFace(sFaceModel string, targetRect image.Rectangle) *Face {
	var FaceFeature face_rec.IFaceFeature
	var err error
	FaceFeature, err = rec_gpu.New("/root/face_recognize/recognize/libs/face_gpu/models")
	if err != nil {
		log.Fatal(err)
	}

	sFace := seetaFace.NewSeetaFace(sFaceModel, targetRect)

	face := &Face{
		RecognizeGpu:  FaceFeature,
		Seeta:         sFace,
		TargetRect:    targetRect,
		trackedBuffer: make(chan *Frame, 3),
		TrackState: TrackState{
			MaxEmptyCount: 15,
		},
	}

	return face
}
