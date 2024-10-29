package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"
	"test/seetaFace"
	"time"
)

var faceOutput = "./faceOutput"

func init() {
	os.MkdirAll(faceOutput, 0755)

	seetaFace.faceIdMap = make(map[int]struct{})
}

func main() {

	sface := seetaFace.NewSeetaFace(seetaFace.modelPath)

	inputName := strings.ToLower(os.Args[1])

	var videoCapture *gocv.VideoCapture
	var err error
	if strings.HasSuffix(inputName, ".png") || strings.HasSuffix(inputName, ".jpg") {
		sface.ImageProcess(inputName)
		return
	} else if strings.HasSuffix(inputName, ".mp4") {
		videoCapture, err = gocv.VideoCaptureFile(inputName)
	} else {
		videoCapture, err = gocv.VideoCaptureDevice(0)
	}

	if err != nil {
		log.Fatal(err)
	}

	frame := gocv.NewMat()
	defer frame.Close()

	// 人脸检测框颜色
	borderColor := color.RGBA{0, 255, 0, 0}

	// 视频窗口
	window := gocv.NewWindow("人脸检测")
	defer window.Close()

	frameCount := 0
	for {
		frameCount++
		if ok := videoCapture.Read(&frame); !ok {
			fmt.Println("视频结束或无法读取")
			break
		}
		if frame.Empty() {
			continue
		}

		log.Println("=====================================")

		start := time.Now()
		faces, err := sface.Detect(frame)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("frameCount:", frameCount, "timeSince:", time.Now().Sub(start).Milliseconds())

		if len(faces) > 0 {
			faceIdExists := false
			for i, face := range faces {
				gocv.Rectangle(&frame, face.Rects, borderColor, 2)
				log.Println("faceInfo", "index:", i+1, "rect:", face.Rects, "pid:", face.PID, "score:", face.Score)

				if !faceIdExists {
					_, faceIdExists = seetaFace.faceIdMap[face.PID]
					seetaFace.faceIdMap[face.PID] = struct{}{}
				}
			}

			if !faceIdExists {
				saveImage(frame, frameCount, faces)
			}

		}

		window.IMShow(frame)

		if window.WaitKey(1) >= 0 {
			break
		}
		time.Sleep(66 * time.Millisecond)
	}
}

func saveImage(frame gocv.Mat, frameCount int, faceInfos []faceInfo) {
	gocv.IMWrite(filepath.Join(faceOutput, fmt.Sprintf("%d_src.jpg", frameCount)), frame)

	for i, info := range faceInfos {
		faceRegion := frame.Region(info.Rects)
		gocv.IMWrite(filepath.Join(faceOutput, fmt.Sprintf("frame_%d_pid_%d.jpg", frameCount, i)), faceRegion)
		faceRegion.Close()
	}
}
