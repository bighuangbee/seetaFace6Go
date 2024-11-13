package main

import (
	"flag"
	"log"
	video_find_face "video-find-face"
)

func main() {
	videoPath := flag.String("videoPath", "", "视频地址或本地视频目录, rtsp://或./video")
	picturePath := flag.String("picPath", "", "抓拍图目录")
	flag.Parse()

	log.Println("videoPath:", *videoPath, "picturePath:", *picturePath)

	video_find_face.VideoRecTrim(*videoPath)
}
