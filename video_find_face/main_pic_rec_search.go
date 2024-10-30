package main

import (
	"face_recognize/recognize/face_rec"
	"flag"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"os"
	"path/filepath"
	"strings"
	"video-find-face/common"
)

func main() {
	var seachResultPath = "./seachResult"
	os.MkdirAll(seachResultPath, 0755)

	regInfos := map[uint32]string{}
	// 视频窗口
	//window = gocv.NewWindow("人脸检测")
	//defer window.Close()

	regPath := flag.String("regPath", "", "注册目录")
	searchPath := flag.String("searchPath", "", "搜索目标")
	flag.Parse()

	var targetRect = image.Rectangle{
		Min: image.Point{0, 600},
		Max: image.Point{1600, 2160},
	}
	var face = common.NewFace("../seetaFace6Warp/seeta/models", targetRect)

	regFiles, _ := common.GetFilesName(*regPath)

	regFilesCount := uint32(0)
	fs := []*face_rec.FaceEntity{}
	for _, filename := range regFiles {
		if strings.Contains(filename, "pid") {
			continue
		}
		mat := gocv.IMRead(filepath.Join(*regPath, filename), gocv.IMReadColor)
		fe, err := face.Recognize(mat)
		if err != nil {
			fmt.Println(err)
			continue
		}
		regFilesCount++

		for i, entity := range fe {
			fe[i].Id = regFilesCount
			regInfos[fe[i].Id] = filename
			fe[i].RegisteFile = face_rec.RegisteFile{
				Filename: filename,
			}
			fmt.Println("reg", filename, i, entity.Quality)
		}
		fs = append(fs, fe...)
	}

	err := face.FaceFeature.SetFeatures(fs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("reg SetFeatures, regFilesLen", regFilesCount, "featureLen", len(fs))

	//搜索
	searchFiles, _ := common.GetFilesName(*searchPath)
	for _, filename := range searchFiles {
		if strings.Contains(filename, "pid") {
			continue
		}
		mat := gocv.IMRead(filepath.Join(*searchPath, filename), gocv.IMReadColor)
		fe, err := face.Recognize(mat)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dir := filepath.Join(seachResultPath, filename)
		os.MkdirAll(dir, 0755)
		err = common.CopyFile(filepath.Join(*searchPath, filename), filepath.Join(dir, "input_"+filename))
		if err != nil {
			fmt.Println(err)
		}

		for _, entity := range fe {
			results, err := face.FaceFeature.CompareFeaturesByRegister(entity, 10, 0.6)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for i, result := range results {
				err = common.CopyFile(filepath.Join(*regPath, regInfos[result.Id]), filepath.Join(dir, fmt.Sprintf("%0.3f", result.Match)+regInfos[result.Id]))
				if err != nil {
					fmt.Println("results ", err)
				}
				fmt.Println("search results, inputname:", filename, "resultIndex:", i, "match:", result.Match, "regFile:", regInfos[result.Id])
			}
		}

	}
}
