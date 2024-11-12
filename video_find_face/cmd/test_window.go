package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
	video_find_face "video-find-face"
	"video-find-face/seetaFace"
)

func main() {
	videoPath := os.Args[1]
	videoCapture, err := gocv.VideoCaptureFile(videoPath)
	if err != nil {
		fmt.Println("无法打开视频文件:", err)
		return
	}
	defer videoCapture.Close()

	window := gocv.NewWindow("Video")
	defer window.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	min := image.Point{0, 500}
	var targetRect = image.Rectangle{
		Min: min,
		Max: image.Point{min.X + 3840*2/3, min.Y + 2160*2/3},
	}
	//targetRect = image.Rectangle{}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	//帧计数器
	frameCount := int32(0)
	//FPS计数器
	processingCount := int32(0)

	go func() {
		timeCount := 0
		for range ticker.C {
			timeCount++
			log.Printf("DEBUG, 视频名称:%s,第%d秒,FPS:%d,已处理%d帧\n", filepath.Base(videoPath), timeCount, atomic.LoadInt32(&processingCount), atomic.LoadInt32(&frameCount))
			atomic.StoreInt32(&processingCount, 0)
		}
	}()

	face := video_find_face.NewFace("../../seetaFace6Warp/seeta/models", targetRect)

	borderColor := color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 1,
	}

	// 读取视频帧并显示
	for {
		atomic.AddInt32(&frameCount, 1)
		atomic.AddInt32(&processingCount, 1)

		if ok := videoCapture.Read(&frame); !ok {
			fmt.Println("视频结束或无法读取")
			break
		}
		if frame.Empty() {
			continue
		}

		//if frameCount < 35*15 {
		//	continue
		//}

		t1 := time.Now()
		face.Seeta.NewTracker(frame.Cols(), frame.Rows())
		img := seetaFace.ToSeetaImage(frame, targetRect)
		faces := face.Seeta.Tracker.Track(img)
		if len(faces) > 0 {
			fmt.Println("===================== Track faceLen:", len(faces), "time:", time.Since(t1), "=====================")
			for i, info := range faces {
				fmt.Println("### Track info:", i, "PID:", info.PID, "Score:", info.Score, "Postion:", *info.Postion)

				originalX := info.Postion.GetX() + face.TargetRect.Min.X
				originalY := info.Postion.GetY() + face.TargetRect.Min.Y

				gocv.PutText(&frame, fmt.Sprintf("pid:%d, i:%d", info.PID, i), image.Point{
					X: originalX,
					Y: originalY - 20,
				}, gocv.FontHersheyPlain, 3.0, borderColor, 2)

				// 绘制人脸框
				gocv.Rectangle(&frame, image.Rectangle{
					Min: image.Point{originalX, originalY},
					Max: image.Point{originalX + info.Postion.GetWidth(), originalY + info.Postion.GetHeight()},
				}, borderColor, 2)
			}

			gocv.IMWrite(fmt.Sprintf("%d.jpg", frameCount), frame)
		}

		gocv.Rectangle(&frame, targetRect, color.RGBA{0, 255, 0, 0}, 2)

		window.IMShow(frame)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
