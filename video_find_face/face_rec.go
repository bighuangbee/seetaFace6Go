package video_find_face

import (
	"face_recognize/recognize/face_rec"
	"gocv.io/x/gocv"
	"image"
)

func (face *Face) RecognizeFrame(frame *Frame) ([]*face_rec.FaceEntity, error) {
	if face.FaceFeature != nil {
		return face.Recognize(*frame.Mat)
	} else {
		results := make([]*face_rec.FaceEntity, 0)
		infos := face.Seeta.Detect(frame.ToSeetaImage(face.TargetRect))
		for _, info := range infos {
			results = append(results, &face_rec.FaceEntity{
				Quality: info.FaceInfo.Score,
				Rect:    image.Rect(info.FaceInfo.Postion.GetX(), info.FaceInfo.Postion.GetY(), info.FaceInfo.Postion.GetX()+info.FaceInfo.Postion.GetWidth(), info.FaceInfo.Postion.GetY()+info.FaceInfo.Postion.GetHeight()),
			})
		}
		return results, nil
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
