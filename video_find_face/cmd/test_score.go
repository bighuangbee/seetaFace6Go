package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type ImageScore struct {
	FrameCount int
	Confidence float64
	Brightness float64
	Clarity    float64
	FaceId     int
}

// 比较并返回更好的图像
func betterImage(img1, img2 ImageScore) ImageScore {
	if img1.Confidence > img2.Confidence {
		return img1
	} else if img1.Confidence < img2.Confidence {
		return img2
	}

	if img1.FrameCount > img2.FrameCount {
		return img1
	} else if img1.FrameCount < img2.FrameCount {
		return img2
	}

	if img1.Clarity > img2.Clarity {
		return img1
	} else if img1.Clarity < img2.Clarity {
		return img2
	}

	if img1.Brightness > img2.Brightness {
		return img1
	}
	return img2
}

func main() {
	// 假设我们逐个获取文件，模拟这个过程
	fileList := []string{
		"1_1.00_0.03_0.17_1.jpg",
		"2_1.00_0.03_0.17_1.jpg",
		"3_1.00_0.02_0.11_2.jpg",
		"4_1.00_0.02_0.26_2.jpg",
		"4_0.98_0.02_0.24_1.jpg",
		"5_0.99_0.02_0.17_3.jpg",
		"5_0.93_0.02_0.23_3.jpg",
		"6_0.99_0.02_0.16_2.jpg",
		"6_0.99_0.02_0.30_3.jpg",
		"7_0.98_0.02_0.39_1.jpg",
	}

	// 创建一个映射以存储每个 FaceId 的最佳图像
	bestImages := make(map[int]ImageScore)

	// 逐个处理每个文件
	for _, file := range fileList {
		parts := strings.Split(strings.TrimRight(file, filepath.Ext(file)), "_")
		if len(parts) < 5 {
			continue // 跳过无效数据
		}

		frameCount, _ := strconv.Atoi(parts[0])
		confidence, _ := strconv.ParseFloat(parts[1], 64)
		brightness, _ := strconv.ParseFloat(parts[2], 64)
		clarity, _ := strconv.ParseFloat(parts[3], 64)
		faceId, _ := strconv.Atoi(parts[4]) // 解析 FaceId

		img := ImageScore{
			FrameCount: frameCount,
			Confidence: confidence,
			Brightness: brightness,
			Clarity:    clarity,
			FaceId:     faceId,
		}

		fmt.Println("--faceId", faceId)

		// 获取当前 FaceId 的最佳图像
		if bestImage, found := bestImages[faceId]; found {
			// 更新最佳图像
			bestImages[faceId] = betterImage(bestImage, img)
		} else {
			// 如果没有找到，直接设置为当前图像
			bestImages[faceId] = img
		}
	}

	fmt.Println("===bestImages", len(bestImages))

	// 输出每个 FaceId 的最佳结果
	for faceId, bestImage := range bestImages {
		fmt.Printf("Best Image for FaceId %d - FrameCount: %d, Confidence: %.2f, Brightness: %.2f, Clarity: %.2f\n",
			faceId, bestImage.FrameCount, bestImage.Confidence, bestImage.Brightness, bestImage.Clarity)
	}
}
