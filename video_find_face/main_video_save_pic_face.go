package main

import (
	"flag"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"video-find-face/common"
)

/*
-Wl,-rpath=/root/face_recognize/recognize/libs/face_gpu/sdk/lib/ -lhiar_cluster
export CGO_LDFLAGS="-Wl,-rpath=/hiar_face/seetaFace6Go/seetaFace6Warp/seeta/lib/linux_x64/ -L/hiar_face/seetaFace6Go/seetaFace6Warp/seeta/lib/linux_x64/ -lSeetaFace6Warp -lSeetaEyeStateDetector200 -lSeetaFaceAntiSpoofingX600  -lSeetaFaceDetector600  -lSeetaFaceLandmarker600  -lSeetaFaceRecognizer610  -lSeetaFaceTracking600  -lSeetaGenderPredictor600   -lSeetaPoseEstimation600 -lSeetaQualityAssessor300"
export LD_LIBRARY_PATH=/hiar_face/seetaFace6Go/seetaFace6Warp/seeta/lib/linux_x64/:$LD_LIBRARY_PATH
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_recognize/recognize/libs/face_gpu/sdk/lib/

GNU 9.4.0

export CGO_CXXFLAGS="-I/usr/local/include/opencv4"
export CGO_CFLAGS="-I/usr/local/include/opencv4"
export CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_imgproc -lopencv_highgui -lopencv_videoio -lopencv_imgcodecs -lopencv_objdetect -lopencv_features2d -lopencv_video -lopencv_dnn -lopencv_calib3d"

*/

//go run . "rtsp://admin:Ab123456.@192.168.1.108:554/cam/realmonitor?channel=1&subtype=0"

// 人脸检测框颜色
var borderColor = color.RGBA{0, 255, 0, 0}

var face *common.Face

func init() {
	os.Mkdir(common.Output, 0755)

	var targetRect = image.Rectangle{
		Min: image.Point{0, 600},
		Max: image.Point{1600, 2160},
	}
	face = common.NewFace("../seetaFace6Warp/seeta/models", targetRect)
}

var videoBasePath string

var window *gocv.Window

func main() {
	// 视频窗口
	window = gocv.NewWindow("人脸检测")
	defer window.Close()

	videoPath := flag.String("videoPath", "", "视频地址或本地视频目录, rtsp://或./video")
	picturePath := flag.String("picPath", "", "抓拍图目录")
	flag.Parse()

	fmt.Println(*videoPath, *picturePath)

	videoList := []string{}

	if isVideo(*videoPath) {
		if strings.HasPrefix(*videoPath, "rtsp") {
			videoBasePath = common.ExtractIP(*videoPath)
		} else {
			videoBasePath = filepath.Base(*videoPath)
		}
		videoList = append(videoList, *videoPath)
	} else {
		info, err := os.Stat(*videoPath)
		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() {
			videoBasePath = filepath.Base(*videoPath)

			//抓拍图匹配录像文件
			if picturePath != nil && *picturePath != "" {
				pictures, err := common.GetFilesName(*picturePath)
				if err != nil {
					log.Println("GetFilesName,", err)
				}
				for _, picture := range pictures {
					matchVideo, err := common.FindMatchingVideo(picture, *videoPath)
					if err != nil {
						log.Println("抓拍图匹配视频,", err)
					}

					if matchVideo == "" {
						fmt.Println("抓拍图匹配部不到视频，删除图片, picture:", picture, "matchVideo:", matchVideo)
						if err := os.Remove(filepath.Join(*picturePath, picture)); err != nil {
							log.Println(err)
						}
					}
					fmt.Println("抓拍图匹配视频, picture:", picture, "matchVideo:", matchVideo)
				}
			}

			//获取视频文件
			videoFiles, err := common.GetFilesName(*videoPath)
			if err != nil {
				log.Fatal(err)
			}

			for _, v := range videoFiles {
				if isVideo(v) {
					videoList = append(videoList, filepath.Join(*videoPath, v))
				}
			}

		}
	}

	common.Output = filepath.Join(common.Output, videoBasePath+"_"+time.Now().Format("01021504"))
	os.MkdirAll(common.Output, 0755)

	for _, v := range videoList {
		face.VideoName = v
		err := videoRecognize(v)
		if err != nil {
			log.Println("videoRecognize", err)
		}
	}
}

func isVideo(videoPath string) bool {
	return strings.HasPrefix(videoPath, "rtsp") || strings.HasSuffix(videoPath, ".mp4") || strings.HasSuffix(videoPath, ".dav")
}

func videoRecognize(videoPath string) error {
	var videoCapture *gocv.VideoCapture
	var err error

	if isVideo(videoPath) {
		videoCapture, err = gocv.VideoCaptureFile(videoPath)
	} else {
		videoCapture, err = gocv.VideoCaptureFile(videoPath)
	}

	if err != nil {
		return err
	}

	defer videoCapture.Close()

	fps := videoCapture.Get(gocv.VideoCaptureFPS)
	fmt.Printf("视频文件: %s, 帧率: %.2f fps\n", videoPath, fps)

	frame := gocv.NewMat()
	defer frame.Close()

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

		//if frameCount%2 == 0 {
		//	continue
		//}

		face.ProcessSaveBestImage(&common.Frame{
			Mat:   &frame,
			Count: frameCount,
		})

		gocv.Rectangle(&frame, face.TargetRect, borderColor, 2)

		// 在窗口中显示帧
		window.IMShow(frame)

		// 等待键盘输入，同时保持窗口更新
		if window.WaitKey(1) >= 0 {
			break
		}
	}

	return nil
}
