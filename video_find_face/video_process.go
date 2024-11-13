package video_find_face

import (
	"errors"
	"face_recognize/recognize/face_rec"
	"gocv.io/x/gocv"
	"image"
	"log"
	"path/filepath"
	"sync/atomic"
	"time"
)

func VideoTracking(videoPath string, targetRect image.Rectangle) error {

	videoCapture, err := OpenVideo(videoPath)
	if err != nil {
		return err
	}
	defer videoCapture.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	if ok := videoCapture.Read(&frame); !ok {
		return errors.New("视频无法读取")
	}

	face := NewFace(targetRect, nil)
	if !targetRect.Empty() {
		face.Seeta.NewTracker(targetRect.Max.X, targetRect.Max.Y)
	} else {
		face.Seeta.NewTracker(frame.Cols(), frame.Rows())
	}

	face.VideoInfo = &VideoInfo{
		Name:       videoPath,
		FPS:        videoCapture.Get(gocv.VideoCaptureFPS),
		TotalFrame: videoCapture.Get(gocv.VideoCaptureFrameCount),
		Width:      frame.Cols(),
		Height:     frame.Rows(),
	}
	face.VideoWriter.PreFrameCount = 30

	//帧计数器
	frameCount := int32(0)
	//FPS计数器
	processingCount := int32(0)

	go func() {
		timeCount := 0

		videoName := filepath.Base(videoPath)
		if IsVideoStream(videoPath) {
			videoName = ExtractIP(videoPath)
		}

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			timeCount++
			log.Printf("DEBUG, 视频名称:%s,第%d秒,FPS:%d,已处理%d帧\n", videoName, timeCount, atomic.LoadInt32(&processingCount), atomic.LoadInt32(&frameCount))
			atomic.StoreInt32(&processingCount, 0)
		}
	}()

	log.Printf("--------------------------------------------------\n")
	log.Printf("开始处理视频，文件名: %s, 帧率: %.1f fps, 总帧数: %0.1f\n", face.VideoInfo.Name, face.VideoInfo.FPS, face.VideoInfo.TotalFrame)

	for {
		atomic.AddInt32(&frameCount, 1)
		atomic.AddInt32(&processingCount, 1)

		if ok := videoCapture.Read(&frame); !ok {
			log.Println("视频结束或无法读取", frameCount, face.VideoInfo.Name)
			break
		}
		if frame.Empty() {
			continue
		}

		if face.VideoInfo.TotalFrame < 0 {
			face.VideoInfo.TotalFrame = float64(atomic.LoadInt32(&frameCount))
		}

		face.Tracking(&Frame{
			Mat:   &frame,
			Count: int(atomic.LoadInt32(&frameCount)),
		})

		//gocv.Rectangle(&frame, face.TargetRect, color.RGBA{0, 255, 0, 0}, 2)
		//window.IMShow(frame)
		//if window.WaitKey(33) >= 0 {
		//	break
		//}
	}
	return nil
}

//--------

func VideoRecTrim(videoPath string) error {

	//人脸识别初始化
	//读取视频帧

	//识别
	//裁剪视频，去除头尾冗余
	//保存最佳图像

	videoCapture, err := OpenVideo(videoPath)
	if err != nil {
		return err
	}
	defer videoCapture.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	if ok := videoCapture.Read(&frame); !ok {
		return errors.New("视频无法读取")
	}

	var FaceFeature face_rec.IFaceFeature
	//var err error
	//FaceFeature, err = rec_gpu.New("/root/face_recognize/recognize/libs/face_gpu/models")
	//if err != nil {
	//	log.Fatal(err)
	//}

	face := NewFace(image.Rectangle{}, FaceFeature)

	face.VideoInfo = &VideoInfo{
		Name:       videoPath,
		FPS:        videoCapture.Get(gocv.VideoCaptureFPS),
		TotalFrame: videoCapture.Get(gocv.VideoCaptureFrameCount),
		Width:      frame.Cols(),
		Height:     frame.Rows(),
	}
	face.VideoWriter.PreFrameCount = 20

	//帧计数器
	frameCount := int32(0)
	//FPS计数器
	processingCount := int32(0)

	go func() {
		timeCount := 0

		videoName := filepath.Base(videoPath)
		if IsVideoStream(videoPath) {
			videoName = ExtractIP(videoPath)
		}

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			timeCount++
			log.Printf("DEBUG, 视频名称:%s,第%d秒,FPS:%d,已处理%d帧\n", videoName, timeCount, atomic.LoadInt32(&processingCount), atomic.LoadInt32(&frameCount))
			atomic.StoreInt32(&processingCount, 0)
		}
	}()

	log.Printf("--------------------------------------------------\n")
	log.Printf("开始处理视频，文件名: %s, 帧率: %.1f fps, 总帧数: %0.1f\n", face.VideoInfo.Name, face.VideoInfo.FPS, face.VideoInfo.TotalFrame)

	for {
		atomic.AddInt32(&frameCount, 1)
		atomic.AddInt32(&processingCount, 1)

		if ok := videoCapture.Read(&frame); !ok {
			log.Println("视频结束或无法读取", frameCount, face.VideoInfo.Name)
			break
		}
		if frame.Empty() {
			continue
		}

		if face.VideoInfo.TotalFrame < 0 {
			face.VideoInfo.TotalFrame = float64(atomic.LoadInt32(&frameCount))
		}

		face.RecognizeProcess(&Frame{
			Mat:   &frame,
			Count: int(atomic.LoadInt32(&frameCount)),
		})

	}
	return nil
}
