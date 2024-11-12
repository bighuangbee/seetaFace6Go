package rec_gpu

import (
	"face_recognize/recognize/face_rec"
	"face_recognize/recognize/face_rec_gpu"
)

func New(modelPath string) (face_rec.IFaceFeature, error) {
	faceGpu, err := face_rec_gpu.New(modelPath)
	if err != nil {
		return faceGpu, nil
	}

	return faceGpu, nil
}
