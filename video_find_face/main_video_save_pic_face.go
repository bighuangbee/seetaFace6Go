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
)

// 人脸检测框颜色
var borderColor = color.RGBA{0, 255, 0, 0}

/*
export CGO_LDFLAGS="-Wl,-rpath=/root/face_recognize/recognize/libs/face_gpu/sdk/lib/ -lhiar_cluster -Wl,-rpath=/hiar_face/seetaFace6Go/seetaFace6Warp/seeta/lib/linux_x64/ -L/hiar_face/seetaFace6Go/seetaFace6Warp/seeta/lib/linux_x64/ -lSeetaFace6Warp -lSeetaEyeStateDetector200 -lSeetaFaceAntiSpoofingX600  -lSeetaFaceDetector600  -lSeetaFaceLandmarker600  -lSeetaFaceRecognizer610  -lSeetaFaceTracking600  -lSeetaGenderPredictor600   -lSeetaPoseEstimation600 -lSeetaQualityAssessor300"
export LD_LIBRARY_PATH=/hiar_face/seetaFace6Go/seetaFace6Warp/seeta/lib/linux_x64/:$LD_LIBRARY_PATH
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_recognize/recognize/libs/face_gpu/sdk/lib/

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_recognize/recognize/libs/face_gpu
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_recognize/recognize/libs/face_gpu/sdk/lib/
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_recognize/recognize/libs/face_gpu/thirdparty/onnxruntime/lib/
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_recognize/recognize/libs/face_gpu/thirdparty/opencv4-ffmpeg/lib/

GNU 9.4.0

export CGO_CXXFLAGS="-I/usr/local/include/opencv4"
export CGO_CFLAGS="-I/usr/local/include/opencv4"
export CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_imgproc -lopencv_highgui -lopencv_videoio -lopencv_imgcodecs -lopencv_objdetect -lopencv_features2d -lopencv_video -lopencv_dnn -lopencv_calib3d"



*/

//go run . "rtsp://admin:Ab123456.@192.168.1.108:554/cam/realmonitor?channel=1&subtype=0"

//var window *gocv.Window

var face *Face

func init() {
	os.Mkdir(faceOutput, 0755)

	var targetRect = image.Rectangle{
		Min: image.Point{0, 600},
		Max: image.Point{1600, 2160},
	}
	face = NewFace("../seetaFace6Warp/seeta/models", targetRect)
}

func main() {
	// 视频窗口
	//window = gocv.NewWindow("人脸检测")
	//defer window.Close()

	videoPath := flag.String("videoPath", "", "视频地址或本地视频目录, rtsp://或./video")
	picturePath := flag.String("picPath", "", "抓拍图目录")
	flag.Parse()

	fmt.Println(*videoPath, *picturePath)

	videoList := []string{}

	if isVideo(*videoPath) {
		videoList = append(videoList, *videoPath)
	} else {
		info, err := os.Stat(*videoPath)
		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() {

			//抓拍图匹配录像文件
			if picturePath != nil && *picturePath != "" {
				pictures, err := GetFilesName(*picturePath)
				if err != nil {
					log.Println("GetFilesName,", err)
				}
				for _, picture := range pictures {
					matchVideo, err := findMatchingVideo(picture, *videoPath)
					if err != nil {
						log.Println("抓拍图匹配视频,", err)
					}

					//if matchVideo == "" {
					//	fmt.Println("-----匹配不到录像文件", picture)
					//	os.Remove(filepath.Join(*picturePath, picture))
					//
					//} else {
					//	fmt.Println("-----匹配录像文件", picture, matchVideo)
					//}
					fmt.Println("抓拍图匹配视频, picture:", picture, "matchVideo:", matchVideo)
				}
			}

			//获取视频文件
			videoFiles, err := GetFilesName(*videoPath)
			if err != nil {
				log.Fatal(err)
			}

			for _, v := range videoFiles {
				if strings.HasSuffix(v, ".mp4") {
					videoList = append(videoList, filepath.Join(*videoPath, v))
				}
			}

		}
	}

	for _, v := range videoList {
		videoName = v
		err := videoRecognize(v)
		if err != nil {
			log.Println("videoRecognize", err)
		}
	}
}

func isVideo(videoPath string) bool {
	return strings.HasPrefix(videoPath, "rtsp") || strings.HasSuffix(videoPath, ".mp4")
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

		face.detectFace(&frame, frameCount)

		gocv.Rectangle(&frame, face.targetRect, borderColor, 2)

		//// 在窗口中显示帧
		//window.IMShow(frame)
		//
		//// 等待键盘输入，同时保持窗口更新
		//if window.WaitKey(1) >= 0 {
		//	break
		//}
	}

	return nil
}

//
//func gocvMatToImage(frame gocv.Mat, filename string) *face_rec.Image {
//	buf, err := gocv.IMEncode(".jpg", frame)
//	if err != nil {
//		return nil
//	}
//
//	image.new
//	seetaFace6go.NewSeetaImageDataFromImage()
//
//	image, err := seetaFace6go.(buf.GetBytes(), filename)
//	if err != nil {
//		return nil
//	}
//	return image
//}
