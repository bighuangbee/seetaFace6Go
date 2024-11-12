package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"video-find-face"
)

//var face *video_find_face.Face

func init() {

	//min := image.Point{0, 500}
	//var targetRect = image.Rectangle{
	//	//Min: min,
	//	//Max: image.Point{min.X + 2844, min.Y + 1600},
	//}
	//
	//face = video_find_face.NewFace("../../seetaFace6Warp/seeta/models", targetRect)
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
	fmt.Println(videoBasePath)

	if video_find_face.IsVideo(*videoPath) {
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
				if video_find_face.IsVideo(v) {
					videoList = append(videoList, filepath.Join(*videoPath, v))
				}
			}

		}
	}

	processConcurrency(videoList)
}

func processConcurrency(videos []string) {
	numCPU := runtime.NumCPU()
	parallelism := numCPU / 4

	parallelism = 3

	var wg sync.WaitGroup

	log.Println("=============并行任务数量:", parallelism, "numCPU:", numCPU)

	sem := make(chan struct{}, parallelism)

	min := image.Point{0, 600}
	var targetRect = image.Rectangle{
		Min: min,
		Max: image.Point{min.X + 3840*2/3, min.Y + 2160*2/3},
	}

	//targetRect = image.Rectangle{
	//	Min: image.Point{0, 0},
	//	Max: image.Point{frame.Cols(), frame.Rows()},
	//}

	for _, video := range videos {
		sem <- struct{}{}
		wg.Add(1)

		go func(videoPath string) {
			defer func() { <-sem }()
			if err := video_find_face.VideoTracking(videoPath, targetRect); err != nil {
				log.Println("videoRecognize", err, videoPath)
			}
			wg.Done()
		}(video)
	}

	// 等待所有 goroutine 完成
	wg.Wait()
}
