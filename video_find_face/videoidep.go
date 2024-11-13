package video_find_face

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type VideoInfo struct {
	Name          string
	FPS           float64
	TotalFrame    float64
	Width, Height int
}

func (videoInfo *VideoInfo) SaveVideo(startFrame, endFrame float64) (string, error) {
	savePath := filepath.Join(filepath.Dir(videoInfo.Name), "output", filepath.Base(videoInfo.Name))
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return "", err
	}

	videoOutputName := filepath.Join(savePath, fmt.Sprintf("%d_%d.mp4", int(startFrame), int(endFrame)))
	return videoOutputName, ExtractVideoSegment(videoInfo.Name, videoOutputName, startFrame, endFrame, videoInfo.TotalFrame)
}

func ExtractVideoSegment(videoPath, outputPath string, start, end, totalFrame float64) error {
	videoCapture, err := gocv.VideoCaptureFile(videoPath)
	if err != nil {
		return fmt.Errorf("无法打开视频文件: %v", err)
	}
	defer videoCapture.Close()

	fps := videoCapture.Get(gocv.VideoCaptureFPS)
	totalFrames := videoCapture.Get(gocv.VideoCaptureFrameCount)
	if totalFrames <= 0 {
		totalFrames = totalFrame
	}

	if start >= totalFrames {
		return fmt.Errorf("开始帧超出范围")
	}
	if start < 0 {
		start = 0
	}
	if end > totalFrames {
		end = totalFrames
	}

	log.Printf("裁剪视频, 名称: %s, 帧率: %.2f fps, 总帧数: %0.1f, 开始帧: %0.1f, 结束帧: %0.1f\n", videoPath, fps, totalFrames, start, end)

	writer, err := gocv.VideoWriterFile(outputPath, "mp4v",
		fps,
		int(videoCapture.Get(gocv.VideoCaptureFrameWidth)),
		int(videoCapture.Get(gocv.VideoCaptureFrameHeight)), true)
	if err != nil {
		return fmt.Errorf("无法创建输出视频文件: %v", err)
	}
	defer writer.Close()

	videoCapture.Set(gocv.VideoCapturePosFrames, start)

	frame := gocv.NewMat()
	defer frame.Close()

	for frameIndex := start; frameIndex < end; frameIndex++ {
		if ok := videoCapture.Read(&frame); !ok || frame.Empty() {
			fmt.Println("视频读取结束或无法读取")
			break
		}
		writer.Write(frame)

		//currentFrame := int(videoCapture.Get(gocv.VideoCapturePosFrames))
		//fmt.Println("write frameIndex", frameIndex, "currentFrame", currentFrame)
		//
		//gocv.IMWrite(fmt.Sprintf("%d.jpg", frameIndex), frame)
	}
	return nil
}

func IsVideo(videoPath string) bool {
	return strings.HasPrefix(videoPath, "rtsp") || strings.HasSuffix(videoPath, ".mp4") || strings.HasSuffix(videoPath, ".dav")
}

func IsVideoStream(videoPath string) bool {
	return strings.HasPrefix(videoPath, "rtsp")
}

func OpenVideo(videoPath string) (videoCapture *gocv.VideoCapture, err error) {
	if IsVideo(videoPath) {
		videoCapture, err = gocv.VideoCaptureFile(videoPath)
	} else {
		videoCapture, err = gocv.VideoCaptureFile(videoPath)
	}

	if err != nil {
		return nil, err
	}
	return videoCapture, err
}
