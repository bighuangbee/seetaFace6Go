package video_find_face

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"path/filepath"
)

func ExtractVideoSegment(videoPath, outputPath string, start, end int) error {
	videoCapture, err := gocv.VideoCaptureFile(videoPath)
	if err != nil {
		return fmt.Errorf("无法打开视频文件: %v", err)
	}
	defer videoCapture.Close()

	fps := videoCapture.Get(gocv.VideoCaptureFPS)
	totalFrames := int(videoCapture.Get(gocv.VideoCaptureFrameCount))
	log.Printf("截取视频, 名称: %s, 帧率: %.2f fps, 总帧数: %d, 开始帧: %d, 结束帧: %d\n", filepath.Base(videoPath), fps, totalFrames, start, end)

	if start >= totalFrames {
		return fmt.Errorf("开始帧超出范围")
	}
	if start < 0 {
		start = 0
	}
	if end > totalFrames {
		end = totalFrames
	}

	writer, err := gocv.VideoWriterFile(outputPath, "mp4v",
		fps,
		int(videoCapture.Get(gocv.VideoCaptureFrameWidth)),
		int(videoCapture.Get(gocv.VideoCaptureFrameHeight)), true)
	if err != nil {
		return fmt.Errorf("无法创建输出视频文件: %v", err)
	}
	defer writer.Close()

	videoCapture.Set(gocv.VideoCapturePosFrames, float64(start))

	frame := gocv.NewMat()
	defer frame.Close()

	for frameIndex := start; frameIndex < end; frameIndex++ {
		if ok := videoCapture.Read(&frame); !ok || frame.Empty() {
			fmt.Println("视频读取结束或无法读取")
			break
		}
		writer.Write(frame)
		log.Printf("写入帧: %d\n", frameIndex) // 打印写入的帧数
	}

	log.Println("视频截取完成")
	return nil
}
func ExtractVideoSegment2(videoPath, outputPath string, startFrame, endFrame int) error {
	videoCapture, err := gocv.VideoCaptureFile(videoPath)
	if err != nil {
		return fmt.Errorf("无法打开视频文件: %v", err)
	}
	defer videoCapture.Close()

	// 获取视频的帧率和总帧数
	fps := videoCapture.Get(gocv.VideoCaptureFPS)
	totalFrames := int(videoCapture.Get(gocv.VideoCaptureFrameCount))
	log.Printf("视频名称: %s, 帧率: %.2f fps, 总帧数: %d, 开始帧: %d, 结束帧: %d\n", filepath.Base(videoPath), fps, totalFrames, startFrame, endFrame)

	// 确保结束帧不超过总帧数
	if endFrame > totalFrames {
		endFrame = totalFrames
	}

	// 创建 VideoWriter 对象
	writer, err := gocv.VideoWriterFile(outputPath, "mp4v",
		fps,
		int(videoCapture.Get(gocv.VideoCaptureFrameWidth)),
		int(videoCapture.Get(gocv.VideoCaptureFrameHeight)), true)
	if err != nil {
		return fmt.Errorf("无法创建输出视频文件: %v", err)
	}
	defer writer.Close()

	// 跳转到开始帧
	videoCapture.Set(gocv.VideoCapturePosFrames, float64(startFrame))

	// 逐帧读取并写入输出视频
	frame := gocv.NewMat()
	defer frame.Close()

	for frameIndex := startFrame; frameIndex < endFrame; frameIndex++ {
		if ok := videoCapture.Read(&frame); !ok || frame.Empty() {
			fmt.Println("视频读取结束或无法读取")
			break
		}
		writer.Write(frame) // 写入当前帧
	}

	return nil
}
