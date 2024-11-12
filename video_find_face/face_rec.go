package video_find_face

import (
	"face_recognize/recognize/face_rec"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"log"
	"time"
)

func (face *Face) Recognize(frame *Frame) {
	if frame.Mat != nil {
		t := time.Now()

		var infos []*face_rec.FaceEntity
		var err error
		if face.RecognizeGpu != nil {
			fmt.Println("======RecognizeGpu ")
			infos, err = RecognizeGpu(face.RecognizeGpu, *frame.Mat)
		} else {
			results := face.Seeta.Detect(frame.ToSeetaImage(face.TargetRect))
			for _, info := range results {
				infos = append(infos, &face_rec.FaceEntity{
					Quality: info.FaceInfo.Score,
					Rect:    image.Rect(info.FaceInfo.Postion.GetX(), info.FaceInfo.Postion.GetY(), info.FaceInfo.Postion.GetX()+info.FaceInfo.Postion.GetWidth(), info.FaceInfo.Postion.GetY()+info.FaceInfo.Postion.GetHeight()),
				})
			}
		}

		if err != nil {
			log.Println("Recognize", err)
			return
		}

		if len(infos) > 0 {
			for _, info := range infos {
				if frame.Score == 0 {
					frame.Score = info.Quality
				} else {
					if frame.Score < info.Quality {
						frame.Score = info.Quality
					}
				}
			}

			// debug
			if face.bestImage != nil && face.RecognizeGpu != nil {
				fe2, err := RecognizeGpu(face.RecognizeGpu, *face.bestImage.Mat)
				if err != nil {
					log.Println("ExtractFeatureGPU error", err)
				}
				log.Printf("###ExtractFeatureGPU bestImage, count: %d, faceLen: %d, time: %d\n",
					frame.Count, len(fe2), time.Since(t).Milliseconds())

				for _, entity := range infos {
					for _, entiry2 := range fe2 {
						match := face.RecognizeGpu.CompareFeature(entity, entiry2)
						log.Println("=== CompareFeature", match)
					}
				}
			}
			// debug end

			face.SetBestFrame(frame)
		}

		log.Printf("###Recognize, count: %d, faceLen: %d, time: %d, topScore: %0.5f \n",
			frame.Count, len(infos), time.Since(t).Milliseconds(), frame.Score)
	} else {
		//跟踪结束
		face.VideoWriterClose(frame.Count)
		face.ResetBestFrame()
	}
}

func RecognizeGpu(faceRec face_rec.IFaceFeature, frame gocv.Mat) ([]*face_rec.FaceEntity, error) {
	buf, err := gocv.IMEncode(".jpg", frame)
	if err != nil {
		return nil, err
	}

	image, _ := face_rec.ReadImageByFormByte(buf.GetBytes(), "1.jpg")
	faces, err := faceRec.ExtractFeature(image, face_rec.ExtractAll)
	if err != nil {
		return nil, err
	}
	buf.Close()

	return faces, nil
}
