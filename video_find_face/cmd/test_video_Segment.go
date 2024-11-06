package main

import (
	"fmt"
	"os"
	video_find_face "video-find-face"
)

func main() {
	videoPath := os.Args[1]    // 输入视频路径
	outputPath := "output.mp4" // 输出视频路径
	startTime := 216.0         // 开始时间（秒）
	duration := 240.0          // 持续时间（秒）

	err := video_find_face.ExtractVideoSegment(videoPath, outputPath, startTime, duration, 1600)
	if err != nil {
		fmt.Println("出错:", err)
	} else {
		fmt.Println("视频片段已保存到:", outputPath)
	}
}
