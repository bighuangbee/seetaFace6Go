package video_find_face

import (
	"face_recognize/recognize/face_rec"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"log"
	"path/filepath"
	"time"
)

func (face *Face) RecognizeFrame(frame *gocv.Mat, frameCount int, pids []int) {
	t := time.Now()
	mat := gocv.NewMat()
	frame.CopyTo(&mat)
	fe, err := face.Recognize(mat)
	if err != nil {
		log.Println("Recognize error", err)
	}
	defer mat.Close()

	log.Println("Recognize faceLen:", len(fe), "time.Sinde:", time.Since(t).Milliseconds(), "pids", pids)

	needSave := false
	rects := []image.Rectangle{}
	for _, entity := range fe {
		log.Println("Recognize entity", entity.Quality, entity.Rect)
		rects = append(rects, entity.Rect)

		if entity.Quality > 0.6 && !needSave {
			needSave = true
		}
	}

	if needSave {
		picBaseName := filepath.Base(filepath.Base(face.VideoInfo.Name))
		ok := gocv.IMWrite(filepath.Join(Output, fmt.Sprintf("%s_frame_%d.jpg", picBaseName, frameCount)), mat)
		if !ok {
			log.Println("Write image error")
		}

		for i, face := range fe {
			faceRegion := mat.Region(face.Rect)
			gocv.IMWrite(filepath.Join(Output, fmt.Sprintf("%s_frame_%d_pid_%d_q_%0.2f.jpg", picBaseName, frameCount, i, face.Quality)), faceRegion)
			faceRegion.Close()
		}
	}

}

func (face *Face) Recognize(frame gocv.Mat) ([]*face_rec.FaceEntity, error) {
	buf, err := gocv.IMEncode(".jpg", frame)
	if err != nil {
		return nil, err
	}

	image, _ := face_rec.ReadImageByFormByte(buf.GetBytes(), "1.jpg")
	faces, err := face.FaceFeature.ExtractFeature(image, face_rec.ExtractAll)
	if err != nil {
		return nil, err
	}
	buf.Close()

	return faces, nil
}
