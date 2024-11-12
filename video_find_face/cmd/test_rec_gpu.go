package main

import (
	"face_recognize/recognize/face_rec"
	"fmt"
	"log"
	"os"
	"video-find-face/rec_gpu"
)

func main() {
	var err error
	FaceFeature, err := rec_gpu.New("/root/face_recognize/recognize/libs/face_gpu/models")
	if err != nil {
		log.Fatal(err)
	}

	//frame := gocv.IMRead(os.Args[1], gocv.IMReadColor)

	//image, _ := face_rec.ReadImageByFormByte(frame.ToBytes(), "1.jpg")

	image, err := face_rec.ReadImageByFile(os.Args[1])
	faces, err := FaceFeature.ExtractFeature(image, face_rec.ExtractAll)
	fmt.Println("----err, ", err, len(faces))
}
