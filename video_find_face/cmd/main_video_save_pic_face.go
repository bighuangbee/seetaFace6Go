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
	"sync/atomic"
	"time"
	"video-find-face"
)

var face *video_find_face.Face

func init() {
	os.Mkdir(video_find_face.Output, 0755)

	//min := image.Point{0, 500}
	var targetRect = image.Rectangle{
		//Min: min,
		//Max: image.Point{min.X + 2844, min.Y + 1600},
	}

	face = video_find_face.NewFace("../../seetaFace6Warp/seeta/models", targetRect)
}

//var window *gocv.Window

func main() {
	// 视频窗口
	//window = gocv.NewWindow("人脸检测")
	//defer window.Close()

	videoPath := flag.String("videoPath", "", "视频地址或本地视频目录, rtsp://或./video")
	picturePath := flag.String("picPath", "", "抓拍图目录")
	flag.Parse()

	log.Println("videoPath:", *videoPath, "picturePath:", *picturePath)

	videoList := []string{}

	var videoBasePath string

	if isVideo(*videoPath) {
		if strings.HasPrefix(*videoPath, "rtsp") {
			videoBasePath = video_find_face.ExtractIP(*videoPath)
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
			videoBasePath = video_find_face.GetPathName(*videoPath)

			//抓拍图匹配录像文件
			if picturePath != nil && *picturePath != "" {
				pictures, err := video_find_face.GetFilesName(*picturePath)
				if err != nil {
					log.Println("GetFilesName,", err)
				}
				for _, picture := range pictures {
					matchVideo, err := video_find_face.FindMatchingVideo(picture, *videoPath)
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
			videoFiles, err := video_find_face.GetFilesName(*videoPath)
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

	for _, v := range videoList {

		fmt.Println(videoBasePath)

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
	face.VideoName = videoPath

	// 人脸检测框颜色
	var borderColor = color.RGBA{0, 255, 0, 0}

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

	totalFrames := int(videoCapture.Get(gocv.VideoCaptureFrameCount))
	fps := videoCapture.Get(gocv.VideoCaptureFPS)
	face.VideoFPS = fps
	fmt.Printf("视频文件: %s, 帧率: %.2f fps, 帧数: %d\n", videoPath, fps, totalFrames)

	frame := gocv.NewMat()
	defer frame.Close()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	frameCount := int32(0)
	processingCount := 0

	go func() {
		timeCount := 0
		for range ticker.C {
			timeCount++
			log.Printf("处理效率, 第%d秒, FPS:%d\n", timeCount, processingCount)
			processingCount = 0
		}
	}()

	for {
		atomic.AddInt32(&frameCount, 1)
		processingCount++

		if ok := videoCapture.Read(&frame); !ok {
			fmt.Println("视频结束或无法读取", frameCount)

			face.Seeta.ResetTracker()
			break
		}
		if frame.Empty() {
			continue
		}

		face.Process(&video_find_face.Frame{
			Mat:   &frame,
			Count: int(atomic.LoadInt32(&frameCount)),
		})

		gocv.Rectangle(&frame, face.TargetRect, borderColor, 2)

		//window.IMShow(frame)
		//if window.WaitKey(33) >= 0 {
		//	break
		//}

	}

	return nil
}
