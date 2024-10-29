package main

import (
	"fmt"
	"log"
	"seetaFace6go"
	"time"
)

func main() {
	fmt.Println(123)
	seetaFace6go.InitModelPath("../seetaFace6Warp/seeta/models/")
	facetracker_Test()

	// sf6go.NewFacxeTracker(1280, 720)
}

func facetracker_Test() {
	t1 := time.Now()
	log.Println("人脸追踪测试开始:", time.Now())
	imageData, err := seetaFace6go.NewSeetaImageDataFromFile("testData/duo6.jpeg")
	if err != nil {
		log.Panic(err)
	}
	log.Println(imageData.GetWidth(), "*", imageData.GetHeight())
	ft := seetaFace6go.NewFaceTracker(imageData.GetWidth(), imageData.GetHeight())
	ft.SetInterval(10)
	log.Println("MinFaceSize:", ft.GetMinFaceSize())
	log.Println("Threshold:", ft.GetThreshold())
	log.Println("VideoStable:", ft.GetVideoStable())
	defer ft.Close()
	for i := 0; i < 1; i++ {
		log.Println("---------------")
		t := time.Now()
		faces := ft.Track(imageData)
		faceCount := len(faces)
		log.Printf("追踪人脸%v个,耗时:%v", faceCount, time.Since(t))

		for j := 0; j < faceCount; j++ {
			face := faces[j]
			log.Printf("Postion:%v,PID:%v,Frame_NO:%v", face.Postion, face.PID, face.Frame_NO)
		}
	}

	log.Println("人脸追踪测试结束:", time.Since(t1).Milliseconds())
}
