package video_find_face

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type VideoWriter struct {
	Writer               *gocv.VideoWriter
	videoname            string
	startFrame, endFrame int
}

func NewVideoFile(info *VideoInfo, startFrame, endFrame float64) (*gocv.VideoWriter, string, error) {
	savePath := ""
	if IsVideoStream(info.Name) {
		savePath = filepath.Join("output", ExtractIP(info.Name))
	} else {
		savePath = filepath.Join(filepath.Dir(info.Name), "output", filepath.Base(info.Name))
	}
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return nil, "", err
	}

	videoname := filepath.Join(savePath, fmt.Sprintf("%d_%d.mp4", int(startFrame), int(endFrame)))
	writer, err := gocv.VideoWriterFile(videoname, "mp4v", info.FPS, info.Width, info.Height, true)
	if err != nil {
		return nil, "", err
	}

	return writer, videoname, nil
}

func (face *Face) VideoWriterClose(endFrame int) {

	log.Printf("保存截取视频, 名称: %s, 帧率: %.2f fps, 总帧数: %0.1f, 开始帧: %d, 结束帧: %d\n",
		filepath.Base(face.VideoWriter.videoname), face.VideoInfo.FPS, face.VideoInfo.TotalFrame, face.VideoWriter.startFrame, endFrame)

	face.muVideoWriter.Lock()
	face.VideoWriter.Writer.Close()
	face.VideoWriter.Writer = nil
	face.muVideoWriter.Unlock()

	//去尾。头缓存x帧，尾跟踪冗余，5=冗余
	end := endFrame - face.VideoWriter.startFrame + int(face.VideoInfo.FPS*2) - face.TrackState.MaxEmptyCount
	if face.VideoWriter.startFrame < int(face.VideoInfo.FPS*2) {
		end = endFrame - face.VideoWriter.startFrame
	}

	oldName := face.VideoWriter.videoname
	newName := filepath.Join(filepath.Dir(face.VideoWriter.videoname),
		strings.ReplaceAll(filepath.Base(face.VideoWriter.videoname), "_0", fmt.Sprintf("_%d", face.VideoWriter.startFrame+end)))
	go func() {
		ExtractVideoSegment(oldName, newName, 0, float64(end), 0)
		os.Remove(oldName)
	}()

	if face.bestImage != nil {
		picName := filepath.Join(filepath.Dir(face.VideoWriter.videoname),
			fmt.Sprintf("%s_%0.5f.jpg", strings.ReplaceAll(filepath.Base(face.VideoWriter.videoname), filepath.Ext(face.VideoWriter.videoname), ""), face.bestImage.Score))
		ok := gocv.IMWrite(picName, *face.bestImage.Mat)
		log.Println("保存照片, ok:", ok, picName)
	}

	//output/视频文件名 或 output/录像日期/视频文件名
	//outputName, err := face.VideoInfo.SaveVideo(face.bestImage.CountStart, float64(frame.Count))
	//log.Println("视频片段保存, errInfo:", err, "outputName:", outputName)

	face.muVideoWriter.Lock()
	face.muVideoWriter.Unlock()
}

func (face *Face) StartVideoWriter(startFrame float64) (err error) {
	face.muVideoWriter.Lock()
	if face.VideoWriter.Writer == nil {
		face.VideoWriter.Writer, face.VideoWriter.videoname, err = NewVideoFile(face.VideoInfo, startFrame, 0)
		if err != nil {
			face.muVideoWriter.Unlock()
			return err
		}

	}
	face.muVideoWriter.Unlock()

	//在前追加
	for _, frame := range face.GetFramesBuffer() {
		face.VideoWrite(frame)
		frame.Mat.Close()
	}
	return nil
}

func (face *Face) VideoWrite(frame *Frame) {
	if frame.Mat == nil || frame.Mat.Empty() {
		return
	}

	face.muVideoWriter.Lock()
	if face.VideoWriter.Writer != nil {
		mat := gocv.NewMat()
		frame.Mat.CopyTo(&mat)
		face.VideoWriter.Writer.Write(mat)
		mat.Close()
	}
	face.muVideoWriter.Unlock()
}
