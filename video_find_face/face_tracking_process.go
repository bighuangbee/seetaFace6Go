package video_find_face

import (
	"gocv.io/x/gocv"
	"log"
	"time"
)

func (face *Face) Tracking(frame *Frame) {

	//帧缓存
	face.AddFrameBuffer(frame)
	//截取视频
	face.VideoWrite(frame)

	t := time.Now()

	img := frame.ToSeetaImage(face.TargetRect)
	faces := face.Seeta.Tracker.Track(img)

	if len(faces) > 0 {
		log.Printf("faceTrack, count: %d, faceLen: %d, time: %d \n", frame.Count, len(faces), time.Since(t).Milliseconds())
		if !face.TrackState.Tracking {
			face.TrackState.Tracking = true
			log.Println("人脸跟踪开始")
		}

		if err := face.StartVideoWriter(float64(frame.Count)); err != nil {
			log.Println("StartVideoWriter", err)
		}

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

				//face.StopTracking(frame.Count)
				face.VideoWriterClose(frame.Count)

				//
			}
		}
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
