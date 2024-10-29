package main

import (
	"face_recognize/recognize/face_rec"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"log"
	"path/filepath"
	"seetaFace6go"
	"time"
	"video-find-face/faceRec"
	"video-find-face/seetaFace"
)

type Face struct {
	sFace       *seetaFace.SeetaFace
	FaceFeature face_rec.IFaceFeature
	targetRect  image.Rectangle
}

var videoName string
var faceOutput = "./faceOutput"

func NewFace(sFaceModel string, targetRect image.Rectangle) *Face {
	FaceFeature, err := faceRec.New("/root/face_recognize/recognize/libs/face_gpu/models")
	if err != nil {
		log.Fatal(err)
	}

	sFace := seetaFace.NewSeetaFace(sFaceModel)
	if err != nil {
		log.Fatal(err)
	}

	return &Face{
		FaceFeature: FaceFeature,
		sFace:       sFace,
		targetRect:  targetRect,
	}
}

var frameEmptyCount int
var frameFaces []gocv.Mat

func (face *Face) newTracker(width, height int) {
	if face.sFace.Tracker == nil {
		face.sFace.Tracker = seetaFace6go.NewFaceTracker(face.targetRect.Size().X, face.targetRect.Size().Y)
		//sFace.Tracker = seetaFace6go.NewFaceTracker(frame.Cols(), frame.Rows())
		face.sFace.Tracker.SetVideoStable(true)
		face.sFace.Tracker.SetInterval(1)
		//sFace.Tracker.SetThreads(4)
		face.sFace.Tracker.SetMinFaceSize(60)
	}
}

func (face *Face) recognizeFrame(frame *gocv.Mat, frameCount int, pids []int) {
	t := time.Now()
	mat := gocv.NewMat()
	frame.CopyTo(&mat)
	fe, err := face.recognize(mat)
	if err != nil {
		log.Println("recognize error", err)
	}
	defer mat.Close()

	log.Println("recognize faceLen:", len(fe), "time.Sinde:", time.Since(t).Milliseconds(), "pids", pids)

	needSave := false
	rects := []image.Rectangle{}
	for _, entity := range fe {
		log.Println("recognize entity", entity.Quality, entity.Rect)
		rects = append(rects, entity.Rect)

		if entity.Quality > 0.6 && !needSave {
			needSave = true
		}
	}

	if needSave {
		picBaseName := filepath.Base(videoName)
		ok := gocv.IMWrite(filepath.Join(faceOutput, fmt.Sprintf("%s_frame_%d.jpg", picBaseName, frameCount)), mat)
		fmt.Println("--saveImage", filepath.Join(faceOutput, picBaseName+".jpg"), ok)

		for i, face := range fe {
			faceRegion := mat.Region(face.Rect)
			gocv.IMWrite(filepath.Join(faceOutput, fmt.Sprintf("%s_frame_%d_pid_%d_q_%0.2f.jpg", picBaseName, frameCount, i, face.Quality)), faceRegion)
			faceRegion.Close()
		}
	}

}

func (face *Face) detectFace(frame *gocv.Mat, frameCount int) {
	//fmt.Println("=============")
	start := time.Now()
	face.newTracker(frame.Cols(), frame.Rows())

	faceRegion := frame.Region(face.targetRect)
	frameImage, _ := faceRegion.ToImage()
	faces := face.sFace.Tracker.Track(seetaFace6go.NewSeetaImageDataFromImage(frameImage))

	if frameCount%15 == 0 {
		log.Println("faceTrack, frame:", frameCount, "time:", time.Now().Sub(start).Milliseconds())
	}

	if len(faces) > 0 {
		mat := gocv.NewMat()
		frame.CopyTo(&mat)
		defer mat.Close()
		frameFaces = append(frameFaces, mat)

		pids := []int{}
		for _, face := range faces {
			pids = append(pids, face.PID)
			fmt.Printf("face.Track, Frame_NO: %d, PID: %d, Score: %f, Step: %d, Postion: %v \n", face.Frame_NO, face.PID, face.Score, face.Step, face.Postion)

			// 将人脸框的坐标转换到原图
			//originalX := face.Postion.GetX() + targetRect.Min.X
			//originalY := face.Postion.GetY() + targetRect.Min.Y
			//
			//// 绘制人脸框
			//gocv.Rectangle(frame, image.Rectangle{
			//	Min: image.Point{originalX, originalY},
			//	Max: image.Point{originalX + face.Postion.GetWidth(), originalY + face.Postion.GetHeight()},
			//}, borderColor, 2)
		}

		if face.FaceFeature != nil {
			go face.recognizeFrame(frame, frameCount, pids)
		}
	} else {
		frameEmptyCount++

		//超过x侦没有检测到人脸
		if frameEmptyCount > 50 {
			frameEmptyCount = 0
			frameFaces = []gocv.Mat{}

		}
	}
}

//func recognizeCollegeFrame() {
//	//var recEntityBest []*face_rec.FaceEntity
//	var idCount = uint32(0)
//
//	for i, frame := range frameFaces {
//		t := time.Now()
//		fe, err := recognize(frame)
//		if err != nil {
//			log.Println("recognize error", err)
//			continue
//		}
//
//		if i == 0 {
//			//recEntityBest = fe
//		} else {
//		}
//
//		log.Println("recognize faceLen:", len(fe), "time.Sinde:", time.Since(t).Milliseconds())
//
//		for _, entity := range fe {
//			log.Println("recognize entity", entity.Quality, entity.Rect)
//
//			//标记id
//			idCount++
//			entity.Id = idCount
//			if entity.Quality > 0.6 {
//			}
//		}
//	}
//
//}

func (face *Face) recognize(frame gocv.Mat) ([]*face_rec.FaceEntity, error) {
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

func saveImage(frame gocv.Mat, filename string, frameCount int, rects []image.Rectangle) {

	ok := gocv.IMWrite(filepath.Join(faceOutput, fmt.Sprintf("%s_%d.jpg", filename, frameCount)), frame)
	fmt.Println("--saveImage", filepath.Join(faceOutput, filename+".jpg"), ok)

	for i, info := range rects {
		faceRegion := frame.Region(info)
		gocv.IMWrite(filepath.Join(faceOutput, fmt.Sprintf("%s_frame_%d_pid_%d.jpg", filename, frameCount, i)), faceRegion)
		faceRegion.Close()
	}
}
