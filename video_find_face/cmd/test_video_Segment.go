package main

import (
	"fmt"
	"os"
)

func main() {
	videoPath := os.Args[1]    // 输入视频路径
	outputPath := "output.mp4" // 输出视频路径
	startTime := 5.0           // 开始时间（秒）
	duration := 10.0           // 持续时间（秒）

	err := ExtractVideoSegment(videoPath, outputPath, startTime, duration)
	if err != nil {
		fmt.Println("出错:", err)
	} else {
		fmt.Println("视频片段已保存到:", outputPath)
	}
}
