package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
	"path/filepath"
	"strings"
	"test/seetaFace"
	"time"
)

type faceInfo struct {
	Rects image.Rectangle
	PID   int
	Score float32
}

type FaceRec interface {
	Detect(frame gocv.Mat) ([]faceInfo, error)
	Recognize(frame gocv.Mat) ([][]float32, error)
	Compare(feature1 []float32, feature2 []float32) (score float32)
}

//======================================

var faceRec FaceRec

var window *gocv.Window

func init() {
	faceRec = seetaFace.NewSeetaFace(seetaFace.modelPath)

	seetaFace.faceIdMap = make(map[int]struct{})
}

func main() {

	//frame := gocv.IMRead("./testData/duo6.jpeg", gocv.IMReadColor)
	//features, err := faceRec.Recognize(frame)
	//
	//fmt.Println(err, len(features))
	//return

	// 视频窗口
	window = gocv.NewWindow("人脸检测")
	//defer window.Close()

	videoDir := "/Users/bighuangbee/Pictures/01古龙峡/2024-10-18/video_001/mp4Output"
	photoDir := "/Users/bighuangbee/Pictures/01古龙峡/2024-10-18/pic_001"

	filesName, err := GetFilesName(photoDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, filename := range filesName {
		fmt.Println("=====================================")

		imageName := filepath.Base(filename)

		if strings.Contains(imageName, "1].jpg") {
			continue
		}

		matchingVideo, err := findMatchingVideo(imageName, videoDir)
		if err != nil {
			fmt.Printf("抓拍图: %s, 视频文件匹配结果: %s\n", imageName, err.Error())
			continue
		}

		fmt.Printf("抓拍图: %s, 视频文件匹配结果: %s\n", imageName, matchingVideo)

		//frame := gocv.IMRead("./testData/duo6.jpeg", gocv.IMReadColor)

		frame := gocv.IMRead(filepath.Join(photoDir, filename), gocv.IMReadColor)
		if frame.Empty() {
			fmt.Println("读取抓拍图失败", err, filename)
			continue
		}
		features, err := faceRec.Recognize(frame)
		if err != nil {
			fmt.Println("faceRec.Recognize", err)
			continue
		}

		fmt.Println("----target featuresLen", len(features))

		err = readVideoAndRecognize(features, filepath.Join(videoDir, matchingVideo))
		if err != nil {
			fmt.Println("readVideoAndRecognize", err)
			continue
		}
	}

}
func readVideoAndRecognize(targetFeatures [][]float32, videoFilename string) error {
	videoCapture, err := gocv.VideoCaptureFile(videoFilename)
	if err != nil {
		return err
	}
	defer videoCapture.Close()

	fps := videoCapture.Get(gocv.VideoCaptureFPS)
	fmt.Printf("视频帧率: %.2f fps\n", fps)

	frame := gocv.NewMat()
	defer frame.Close()

	// 人脸检测框颜色
	borderColor := color.RGBA{0, 255, 0, 0}

	//// 视频窗口
	//window := gocv.NewWindow("人脸检测")
	//defer window.Close()

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

		log.Println("------------------------------------")

		start := time.Now()
		faces, err := faceRec.Detect(frame)
		if err != nil {
			fmt.Println("Recognize", err)
			continue
		}

		log.Println("frameCount:", frameCount, "timeSince:", time.Now().Sub(start).Milliseconds())

		if len(faces) > 0 {
			//faceIdExists := false
			for i, face := range faces {
				gocv.Rectangle(&frame, face.Rects, borderColor, 2)
				log.Println("Process faceInfo", "index:", i+1, "rect:", face.Rects, "pid:", face.PID, "score:", face.Score)

				//if !faceIdExists {
				//	_, faceIdExists = faceIdMap[face.PID]
				//	faceIdMap[face.PID] = struct{}{}
				//}

				features, err := faceRec.Recognize(frame)
				if err != nil {
					fmt.Println("faceRec.Recognize", err)
					continue
				}

				if len(features) > 0 {
					for _, feature := range features {
						for _, target := range targetFeatures {
							score := faceRec.Compare(target, feature)
							fmt.Println("-----score", score)
						}
					}
				} else {
					fmt.Println("frame DetectFaceLen:", len(faces), " featuresLen:", len(features))
				}
			}

			//if !faceIdExists {
			//	saveImage(frame, frameCount, faces)
			//	}

		}

		// 在窗口中显示帧
		window.IMShow(frame)

		// 等待键盘输入，同时保持窗口更新
		if window.WaitKey(1) >= 0 {
			break
		}
	}

	return nil
}
