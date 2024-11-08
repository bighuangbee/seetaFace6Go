package video_find_face

import (
	"gocv.io/x/gocv"
	"log"
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
		infos, err := face.RecognizeFrame(frame)
		if err != nil {
			log.Println("RecognizeFrame", err)
			return
		}
		if len(infos) > 0 {
			for _, info := range infos {
				if frame.Score == 0 {
					frame.Score = info.Quality
				} else {
					if frame.Score < info.Quality {
						frame.Score = info.Quality
					}
				}
			}

			if face.bestImage != nil && face.FaceFeature != nil {

				fe2, err := face.Recognize(*face.bestImage.Mat)
				if err != nil {
					log.Println("ExtractFeature", err)
				}
				log.Printf("###ExtractFeature bestImage, count: %d, faceLen: %d, time: %d\n",
					frame.Count, len(fe2), time.Since(t).Milliseconds())

				for _, entity := range infos {
					for _, entiry2 := range fe2 {
						match := face.FaceFeature.CompareFeature(entity, entiry2)
						log.Println("=== CompareFeature", match)
					}

				}
			}

			face.SetBestFrame(frame)
		}

		log.Printf("###Detect, count: %d, faceLen: %d, time: %d, topScore: %0.5f \n",
			frame.Count, len(infos), time.Since(t).Milliseconds(), frame.Score)

		//frame.Mat.Close()
	} else {
		//跟踪结束

		face.VideoWriterClose(frame.Count)

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
