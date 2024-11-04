package video_find_face

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"os"
	"path/filepath"
	"time"
)

func (face *Face) Process(frame *Frame) {
	//t := time.Now()

	face.Seeta.NewTracker(frame.Mat.Cols(), frame.Mat.Rows(), face.TargetRect)

	img := frame.ToSeetaImage(face.TargetRect)
	faces := face.Seeta.Tracker.Track(img)

	//log.Printf("faceTrack, count: %d, faceLen: %d, time: %d \n", frame.Count, len(faces), time.Since(t).Milliseconds())

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

		face.SetFrame(frame)

		//savePath := ""
		//filename := fmt.Sprintf("%d_%0.3f.jpg", frame.Count, frame.Score)
		//ok := gocv.IMWrite(filepath.Join("", filename), *frame.Mat)
		//if !ok {
		//	log.Println("Write image error", filepath.Join(savePath, filename))
		//} else {
		//	log.Println("savePath: ", filepath.Join(savePath, filename))
		//}

	} else {
		//连续x帧检测不到人脸，认为已经过，重置
		if face.Tracking {
			face.EmptyCount++
			if face.EmptyCount > int(face.VideoFPS) {
				face.EmptyCount = 0
				face.Tracking = false

				frame.Mat = nil
				face.SetFrame(frame)

			}
		}
	}
}

func (face *Face) SetFrame(frame *Frame) {
	if frame.Mat != nil {
		mat := gocv.NewMat()
		frame.Mat.CopyTo(&mat)
		face.frames <- &Frame{
			Mat:   &mat,
			Count: frame.Count,
		}
	} else {
		face.frames <- &Frame{
			Count: frame.Count,
		}
	}
}

func (face *Face) FrameProcess() {
	for frame := range face.frames {
		face.ComputeBestFaces(frame)
	}
}

// 计算有多少个/组游客，确定是否漏检
func (face *Face) ComputeBestFaces(frame *Frame) {
	if frame.Mat != nil {
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
				frame.Count, len(infos), time.Since(t).Milliseconds(), frame.Score)
			face.SetBestFrame(frame)
		}

		frame.Mat.Close()
	} else {
		//跟踪结束信号

		outputStart := face.bestImage.CountStart - (int(face.VideoFPS) * 1)
		outputEnd := frame.Count + (int(face.VideoFPS) * 1)
		outputName := fmt.Sprintf("video_output_%d_%d.mp4", outputStart, outputEnd)

		err := ExtractVideoSegment(face.VideoName,
			outputName, outputStart, outputEnd)
		if err != nil {
			log.Println("ExtractVideoSegment", err)
		}

		face.bestImage.CountStart = 0

		//结果输出目录，output/视频文件名 或 output/录像日期/视频文件名
		filename := fmt.Sprintf("%d_%0.5f.jpg", face.bestImage.Count, face.bestImage.Score)
		savePath := filepath.Join("output", GetPathName(face.VideoName), filepath.Base(face.VideoName))
		os.MkdirAll(savePath, 0755)

		ok := gocv.IMWrite(filepath.Join(savePath, filename), *face.bestImage.Mat)
		if !ok {
			log.Println("Write image error", filepath.Join(savePath, filename))
		} else {
			log.Println("savePath: ", filepath.Join(savePath, filename))
		}

		face.bestImage = nil
	}
}

func (face *Face) SetBestFrame(f *Frame) {
	if face.bestImage == nil {
		face.bestImage = &Frame{
			CountStart: f.Count,
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
