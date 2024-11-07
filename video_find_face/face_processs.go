package video_find_face

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func (face *Face) Process(frame *Frame) {

	//帧缓存
	face.AddFrameBuffer(frame)
	//截取视频
	face.VideoWrite(frame)

	//t := time.Now()

	img := frame.ToSeetaImage(face.TargetRect)
	faces := face.Seeta.Tracker.Track(img)

	//log.Printf("faceTrack, count: %d, faceLen: %d, time: %d \n", frame.Count, len(faces), time.Since(t).Milliseconds())
	if len(faces) > 0 {
		if !face.TrackState.Tracking {
			face.TrackState.Tracking = true
			log.Println("人脸跟踪开始")
		}

		if err := face.StartVideoWriter(float64(frame.Count)); err != nil {
			log.Println("StartVideoWriter", err)
		}

		face.AddTracked(frame)

		//for i, info := range faces {
		//	// 将人脸框的坐标转换到原图
		//	originalX := info.Postion.GetX() + face.TargetRect.Min.X
		//	originalY := info.Postion.GetY() + face.TargetRect.Min.Y
		//
		//	gocv.PutText(frame.Mat, fmt.Sprintf("pid:%d, i:%d", info.PID, i), image.Point{
		//		X: originalX,
		//		Y: originalY - 20,
		//	}, gocv.FontHersheyPlain, 3.0, borderColor, 5)
		//	// 绘制人脸框
		//	gocv.Rectangle(frame.Mat, image.Rectangle{
		//		Min: image.Point{originalX, originalY},
		//		Max: image.Point{originalX + info.Postion.GetWidth(), originalY + info.Postion.GetHeight()},
		//	}, borderColor, 2)
		//

	} else {
		//连续x帧检测不到人脸，认为已经过，重置
		if face.TrackState.Tracking {
			face.TrackState.EmptyCount++
			if face.TrackState.EmptyCount > face.TrackState.MaxEmptyCount {
				face.TrackState.EmptyCount = 0
				face.TrackState.Tracking = false

				face.StopTrack(frame.Count)
			}
		}
	}
}

func (face *Face) AddTracked(frame *Frame) {
	if frame.Mat != nil {
		mat := gocv.NewMat()
		frame.Mat.CopyTo(&mat)
		face.trackedBuffer <- &Frame{
			Mat:   &mat,
			Count: frame.Count,
		}
	} else {
		face.trackedBuffer <- &Frame{
			Count: frame.Count,
		}
	}
}

func (face *Face) StopTrack(count int) {
	face.AddTracked(&Frame{Count: count})
}

func (face *Face) TrackedProcess(wg *sync.WaitGroup) {
	for frame := range face.trackedBuffer {
		face.FrameDetect(frame)
	}
	wg.Done()
}

func (face *Face) TrackedProcessClose() {
	close(face.trackedBuffer)
}

func (face *Face) FrameDetect(frame *Frame) {

	if frame.Mat != nil {
		t := time.Now()
		infos := face.Seeta.Detect(frame.ToSeetaImage(face.TargetRect))
		if len(infos) > 0 {
			for _, info := range infos {
				if frame.Score == 0 {
					frame.Score = info.FaceInfo.Score
				} else {
					if frame.Score < info.FaceInfo.Score {
						frame.Score = info.FaceInfo.Score
					}
				}
			}

			face.SetBestFrame(frame)

			log.Printf("###Detect, count: %d, faceLen: %d, time: %d, topScore: %0.5f \n",
				frame.Count, len(infos), time.Since(t).Milliseconds(), frame.Score)
		}

		//frame.Mat.Close()
	} else {
		//跟踪结束

		videoname := face.VideoWriterClose(frame.Count)

		if face.bestImage != nil {
			picName := filepath.Join(filepath.Dir(videoname),
				fmt.Sprintf("%s_%0.5f.jpg", strings.ReplaceAll(filepath.Base(videoname), filepath.Ext(videoname), ""), face.bestImage.Score))
			ok := gocv.IMWrite(picName, *face.bestImage.Mat)
			log.Println("照片保存, ok:", ok, picName)
		}

		//output/视频文件名 或 output/录像日期/视频文件名
		//outputName, err := face.VideoInfo.SaveVideo(face.bestImage.CountStart, float64(frame.Count))
		//log.Println("视频片段保存, errInfo:", err, "outputName:", outputName)

		face.ResetBestFrame()
	}
}

func (face *Face) SetBestFrame(f *Frame) {
	if face.bestImage == nil {
		face.bestImage = &Frame{
			CountStart: float64(f.Count),
		}
	}

	if f.Mat != nil {
		if face.bestImage.Score < f.Score {
			mat := gocv.NewMat()
			f.Mat.CopyTo(&mat)
			face.bestImage.Mat = &mat
			face.bestImage.Count = f.Count
			face.bestImage.Score = f.Score
		}
	}
}

func (face *Face) ResetBestFrame() {
	face.bestImage = nil
}

func (face *Face) AddFrameBuffer(frame *Frame) {
	mat := gocv.NewMat()
	frame.Mat.CopyTo(&mat)

	//缓存x秒
	if len(face.FrameBuffer) >= int(face.VideoInfo.FPS*2) {
		face.FrameBuffer[0].Mat.Close()
		face.FrameBuffer = face.FrameBuffer[1:] // 去掉最早的帧
	}
	face.FrameBuffer = append(face.FrameBuffer, &Frame{
		Mat:        &mat,
		Count:      frame.Count,
		CountStart: frame.CountStart,
		Score:      frame.Score,
	})
}

func (face *Face) GetFramesBuffer() []*Frame {
	frames := face.FrameBuffer
	face.FrameBuffer = nil
	return frames
}
