package video_find_face

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"path/filepath"
	"time"
)

func (face *Face) Process(frame *Frame) {
	t := time.Now()

	face.NewTracker(frame.Mat.Cols(), frame.Mat.Rows())

	img := frame.ToSeetaImage(face.TargetRect)
	faces := face.Seeta.Tracker.Track(img)

	log.Printf("faceTrack, count: %d, faceLen: %d, time: %d \n", frame.Count, len(faces), time.Since(t).Milliseconds())

	if len(faces) > 0 {

		if !face.Tracking {
			face.Tracking = true
			log.Println("人脸跟踪开始")
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
		//	fmt.Println("==track", info.Score, info.Step)
		//	fmt.Println(fmt.Sprintf("%d_%0.2f_%0.2f_%0.2f.jpg", frameCount, info.FaceInfo.Score, brightness.Score, clarity.Score))
		//	ok := gocv.IMWrite(filepath.Join("output", fmt.Sprintf("%d_%0.2f_%0.2f_%0.2f.jpg", frame.Count, info.Score, brightness.Score, clarity.Score)), frame.Mat)
		//	if !ok {
		//		log.Println("Write image error")
		//	}
		//}

		go face.SetFrame(frame)

	} else {
		if face.Tracking {
			face.EmptyCount++
			if face.EmptyCount > 15 {
				face.EmptyCount = 0
				face.Tracking = false

				face.SetFrame(nil)

			}
		}
	}
}

func (face *Face) SetFrame(frame *Frame) {
	if frame != nil {
		mat := gocv.NewMat()
		frame.Mat.CopyTo(&mat)
		face.frames <- &Frame{
			Mat:   &mat,
			Count: frame.Count,
		}
	} else {
		face.frames <- nil
	}
}

func (face *Face) DetectFrame() {
	for frame := range face.frames {
		if frame != nil {
			t := time.Now()
			infos := face.Detect(frame)
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
				log.Printf("###Detect, count: %d, faceLen: %d, time: %d, topScore: %0.5f \n",
					frame.Count, time.Since(t).Milliseconds(), len(infos), frame.Score)
				face.SetBestFrame(frame)
			}

			frame.Mat.Close()
		} else {
			//跟踪结束信号

			//结果输出目录，output/视频文件名 或 output/录像日期/视频文件名
			filename := fmt.Sprintf("%d_%0.5f.jpg", face.bestImage.Count, face.bestImage.Score)
			savePath := filepath.Join("output", GetPathName(face.VideoName), filepath.Base(face.VideoName))

			ok := gocv.IMWrite(filepath.Join(savePath, filename), *face.bestImage.Mat)
			if !ok {
				log.Println("Write image error", filepath.Join(savePath, filename))
			} else {
				log.Println("savePath: ", filepath.Join(savePath, filename))
			}

			face.bestImage = nil
		}
	}
}

func (face *Face) SetBestFrame(f *Frame) {
	mat := gocv.NewMat()
	f.Mat.CopyTo(&mat)

	ff := Frame{
		Mat:   &mat,
		Count: f.Count,
		Score: f.Score,
	}

	if face.bestImage == nil {
		face.bestImage = &ff
	} else {
		if face.bestImage.Score < f.Score {
			face.bestImage = &ff
		}
	}
}
