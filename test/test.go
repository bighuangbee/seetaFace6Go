package main

import (
	"log"
	"os"
	"seetaFace6go"
	"time"
)

func main() {
	seetaFace6go.InitModelPath(os.Args[1])
	standard_Test()

}

func standard_Test() {
	log.Println("标准测试开始:", time.Now())
	// 人脸检测器
	fd := seetaFace6go.NewFaceDetector()
	defer fd.Close()
	// 人脸特征定位器
	// 使用5点信息模型
	fl := seetaFace6go.NewFaceLandmarker(seetaFace6go.ModelType_light)
	defer fl.Close()
	// 人脸特征提取器
	fr := seetaFace6go.NewFaceRecognizer(seetaFace6go.ModelType_light)
	defer fr.Close()
	// 活体检测器（全局）
	//fas := seetaFace6go.NewFaceAntiSpoofing_v2()
	//defer fas.Close()
	//// 口罩检测器
	//md := seetaFace6go.NewMaskDetector()
	//defer md.Close()
	//// 质量评估器
	qr := seetaFace6go.NewQualityCheck()
	defer qr.Close()
	// 如果使用默认值，一下参数可以不设置
	qr.SetBrightnessValues(70, 100, 210, 230)
	qr.SetClarityValues(0.1, 0.2)
	qr.SetIntegrityValues(10, 1.5)

	imageName := os.Args[2]
	imageData, err := seetaFace6go.NewSeetaImageDataFromFile(imageName)
	if err != nil {
		log.Panic(err)
	}

	target := []float32{}

	for i := 0; i < 1; i++ {

		start := time.Now()
		begin := start

		postions := fd.Detect(imageData)
		log.Println("检测人脸", len(postions), "个耗时:", time.Since(start))
		// 人脸特征定位器

		for i := 0; i < len(postions); i++ {
			log.Println("---------------------------------------")
			postion := postions[i].Postion
			log.Printf("识别人脸%v,x:%v,y:%v,width:%v,height:%v", i,
				postion.GetX(), postion.GetY(), postion.GetWidth(), postion.GetHeight(),
			)

			start = time.Now()
			//isMask := md.Detect(imageData, postion)
			//log.Println("口罩检测:", isMask, "耗时:", time.Since(start))
			start = time.Now()
			pointInfo := fl.Mark(imageData, postion)
			log.Println("特征定位耗时:", time.Since(start))
			start = time.Now()
			brightness := qr.CheckBrightness(imageData, postion, pointInfo)
			log.Printf("亮度:%v,检测耗时:%v", brightness.Level, time.Since(start))
			start = time.Now()
			clarity := qr.CheckClarity(imageData, postion, pointInfo)
			log.Printf("清晰度:%v,检测耗时:%v", clarity.Level, time.Since(start))
			start = time.Now()
			integrity := qr.CheckIntegrity(imageData, postion, pointInfo)
			log.Printf("完整度:%v,检测耗时:%v", integrity.Level, time.Since(start))
			start = time.Now()
			pose := qr.CheckPose(imageData, postion, pointInfo)
			log.Printf("姿态:%v,可信度:%v,检测耗时:%v", pose.Level, pose.Score, time.Since(start))
			start = time.Now()
			// 组合方法特征提取
			//success, features := fr.Extract(imageData, pointInfo)

			// 单独人脸裁剪
			face := fr.CropFaceV2(imageData, pointInfo)
			// 通过裁剪的人脸获取特征
			success, features := fr.ExtractCroppedFace(face)

			// ok := true
			// // 两种方法获取特征一致性验证
			// for i := 0; i < len(features_crop); i++ {
			// 	if features[i] != features_crop[i] {
			// 		ok = false
			// 	}
			// }
			log.Println("特征提取", success, len(features), "特征提起方法一致性测试:", "耗时:", time.Since(start))

			if i == 0 {
				target = features
			}

			score := fr.CalculateSimilarity(target, features)
			log.Println("score", score)

			start = time.Now()
			//status := fas.Predict(imageData, postion, pointInfo)
			//log.Println("活体检测", status, "耗时:", time.Since(start))

			//从原始图像中裁剪出人脸
			// img := imageData.CutFace(postion) //face.GetImage()
			// outFile, err := os.Create("temp/" + strconv.Itoa(i) + ".jpeg")
			// if err != nil {
			// 	panic(err)
			// }
			// b := bufio.NewWriter(outFile)
			// err = jpeg.Encode(b, img, &jpeg.Options{Quality: 95})
			// if err != nil {
			// 	panic(err)
			// }
			// err = b.Flush()
			// if err != nil {
			// 	panic(err)
			// }
			// outFile.Close()
		}
		log.Println("单帧总耗时:", time.Since(begin))
	}
	log.Println("标准测试结束:", time.Now())
}
