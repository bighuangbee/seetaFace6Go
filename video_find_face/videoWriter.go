package video_find_face

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"os"
	"path/filepath"
)

type VideoWriter struct {
	Writer               *gocv.VideoWriter
	videoname            string
	startFrame, endFrame int
}

func NewVideoWriter(info *VideoInfo, startFrame, endFrame float64) (*VideoWriter, error) {
	savePath := filepath.Join(filepath.Dir(info.Name), "output", filepath.Base(info.Name))
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return nil, err
	}

	videoname := filepath.Join(savePath, fmt.Sprintf("%d_%d.mp4", int(startFrame), int(endFrame)))
	writer, err := gocv.VideoWriterFile(videoname, "mp4v", info.FPS, info.Width, info.Height, true)
	if err != nil {
		return nil, err
	}

	videoWriter := VideoWriter{
		Writer:     writer,
		videoname:  videoname,
		startFrame: int(startFrame),
	}

	return &videoWriter, nil
}

func (face *Face) ResetVideoWriter(endFrame int) {

	log.Printf("保存截取视频, 名称: %s, 帧率: %.2f fps, 总帧数: %0.1f, 开始帧: %d, 结束帧: %d\n ResetVideoWriter",
		filepath.Base(face.VideoWriter.videoname), face.VideoInfo.FPS, face.VideoInfo.TotalFrame, face.VideoWriter.startFrame, endFrame)

	face.VideoWriter.Writer.Close()
	face.VideoWriter = nil
}

func (face *Face) StartVideoWriter(startFrame float64) error {
	if face.VideoWriter == nil {
		w, err := NewVideoWriter(face.VideoInfo, startFrame, 0)
		if err != nil {
			return err
		}

		face.VideoWriter = w

		//在前追加
		for _, frame := range face.getFramesBuffer() {
			face.VideoWrite(frame)
			frame.Mat.Close()
		}
	}
	return nil
}

func (face *Face) VideoWrite(frame *Frame) {
	if face.VideoWriter != nil && face.VideoWriter.Writer != nil && face.VideoWriter.Writer.IsOpened() {
		face.VideoWriter.Writer.Write(*frame.Mat)
	}
}
