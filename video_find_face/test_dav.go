package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	// 打开 DAV 文件
	videoCapture, err := gocv.VideoCaptureFile(os.Args[1])
	if err != nil {
		fmt.Println("无法打开视频文件:", err)
		return
	}
	defer videoCapture.Close()

	// 创建一个窗口来显示视频
	window := gocv.NewWindow("Video")
	defer window.Close()

	// 创建一个 Mat 用于存储视频帧
	frame := gocv.NewMat()
	defer frame.Close()

	// 读取视频帧并显示
	for {
		if ok := videoCapture.Read(&frame); !ok {
			fmt.Println("视频结束或无法读取")
			break
		}
		if frame.Empty() {
			continue
		}

		window.IMShow(frame)

		// 等待键盘输入，1 毫秒后继续
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
