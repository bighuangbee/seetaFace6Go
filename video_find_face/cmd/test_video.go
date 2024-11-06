package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"time"
)

func checkForStuckOrBuffering(videoPath string) error {
	videoCapture, err := gocv.VideoCaptureFile(videoPath)
	if err != nil {
		return fmt.Errorf("无法打开视频文件: %v", err)
	}
	defer videoCapture.Close()

	// 获取视频的帧率和总帧数
	fps := videoCapture.Get(gocv.VideoCaptureFPS)
	totalFrames := int(videoCapture.Get(gocv.VideoCaptureFrameCount))
	log.Printf("视频名称: %s, 帧率: %.2f fps, 总帧数: %d\n", videoPath, fps, totalFrames)

	// 初始化一些变量
	frame := gocv.NewMat()
	defer frame.Close()

	var lastFrame gocv.Mat
	var lastTimestamp time.Time
	const stuckThreshold = 2 * time.Second // 设置卡住的阈值（例如：2秒没有进展）
	const frameDifferenceThreshold = 0.1   // 设置帧差异阈值，差异小于此值认为帧没有变化

	for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
		if ok := videoCapture.Read(&frame); !ok || frame.Empty() {
			return fmt.Errorf("无法读取视频帧: %d", frameIndex)
		}

		// 如果上一帧存在，计算当前帧和上一帧之间的差异
		if !lastFrame.Empty() {
			// 计算当前帧和上一帧之间的差异
			diff := calculateFrameDifference(frame, lastFrame)
			if diff < frameDifferenceThreshold {
				// 如果帧差异小于阈值，并且时间超过了阈值，认为视频卡住
				if time.Since(lastTimestamp) > stuckThreshold {
					log.Printf("检测到视频卡住: 在第 %d 帧，已经停滞 %s\n", frameIndex, time.Since(lastTimestamp))
					return nil
				}
			} else {
				// 如果帧有明显变化，更新时间戳
				lastTimestamp = time.Now()
			}
		}

		// 更新上一帧
		lastFrame = frame.Clone()

		// 可选：显示当前帧
		// window.IMShow(frame)
		// window.WaitKey(1)
	}

	log.Println("视频播放完毕，没有检测到卡住的情况")
	return nil
}

// 计算两帧之间的差异
func calculateFrameDifference(frame1, frame2 gocv.Mat) float32 {
	// 将帧转换为灰度图
	gray1 := gocv.NewMat()
	gray2 := gocv.NewMat()
	defer gray1.Close()
	defer gray2.Close()

	gocv.CvtColor(frame1, &gray1, gocv.ColorBGRToGray)
	gocv.CvtColor(frame2, &gray2, gocv.ColorBGRToGray)

	// 计算两帧之间的绝对差异
	diff := gocv.NewMat()
	defer diff.Close()
	gocv.AbsDiff(gray1, gray2, &diff)

	// 计算差异图像的总和
	sum := gocv.SumElems(diff)

	// 计算差异值，返回所有通道的和的平均值
	diffValue := float32(sum.Val1+sum.Val2+sum.Val3) / 3.0
	return diffValue
}

func main() {
	videoPath := "your_video.dav"
	err := checkForStuckOrBuffering(videoPath)
	if err != nil {
		log.Fatal(err)
	}
}
