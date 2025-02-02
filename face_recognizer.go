package seetaFace6go

// #include <stdlib.h>
// #include "FaceRecognizer_warp.h"
import "C"
import (
	"path/filepath"
	"unsafe"
)

type FaceRecognizerProperty int

const (
	FaceRecognizer_PROPERTY_NUMBER_THREADS FaceRecognizerProperty = 4 // 人脸识别其线程数
	FaceRecognizer_PROPERTY_ARM_CPU_MODE   FaceRecognizerProperty = 5
)

const (
	FaceDetector_threshold_default = 0.62 // 默认68点模型对比阈值
	FaceDetector_threshold_light   = 0.55 // 轻量5点模型对比阈值
	FaceDetector_threshold_mask    = 0.48 // 口罩5点模型对比阈值
)

var _FaceRecognizer_model = map[ModelType]string{
	ModelType_default: "face_recognizer.csta",
	ModelType_light:   "face_recognizer_light.csta",
	ModelType_mask:    "face_recognizer_mask.csta",
}

type FaceRecognizer struct {
	ptr         *C.struct_facerecognizer
	FeatureSize int
	FaceType    ModelType
}

// NewFaceRecognizer 创建一个人脸识别器
func NewFaceRecognizer(modelType ModelType) *FaceRecognizer {
	model := filepath.Join(_model_base_path, _FaceRecognizer_model[modelType])
	cs := C.CString(model)
	defer C.free(unsafe.Pointer(cs))
	fr := &FaceRecognizer{
		ptr:      C.facerecognizer_new(cs),
		FaceType: modelType,
	}
	fr.SetProperty(FaceRecognizer_PROPERTY_NUMBER_THREADS, 8)

	fr.FeatureSize = fr.getExtractFeatureSize()
	return fr
}

func (s *FaceRecognizer) SetProperty(property FaceRecognizerProperty, value float64) {
	C.facerecognizer_setProperty(s.ptr, C.int(property), C.double(value))
}

func (s *FaceRecognizer) GetProperty(property FaceRecognizerProperty) float64 {
	return float64(C.facerecognizer_getProperty(s.ptr, C.int(property)))
}

func (s *FaceRecognizer) GetCropFaceWidthV2() int {
	return int(C.facerecognizer_GetCropFaceWidthV2(s.ptr))
}
func (s *FaceRecognizer) GetCropFaceHeightV2() int {
	return int(C.facerecognizer_GetCropFaceHeightV2(s.ptr))
}
func (s *FaceRecognizer) GetCropFaceChannelsV2() int {
	return int(C.facerecognizer_GetCropFaceChannelsV2(s.ptr))
}

// GetExtractFeatureSize 获取当前模型的特征数量
func (s *FaceRecognizer) getExtractFeatureSize() int {
	return int(C.facerecognizer_GetExtractFeatureSize(s.ptr))
}

// CalculateSimilarity 对比两个特征的相似度
func (s *FaceRecognizer) CalculateSimilarity(features1, features2 []float32) float32 {
	count := len(features1)
	cfeatures1 := make([]C.float, count)
	cfeatures2 := make([]C.float, count)
	for i := 0; i < count; i++ {
		cfeatures1[i] = C.float(features1[i])
		cfeatures2[i] = C.float(features2[i])
	}

	return float32(C.facerecognizer_CalculateSimilarity(s.ptr, &cfeatures1[0], &cfeatures2[0]))
}

// CropFaceV2 裁剪人脸，裁剪人脸并做人脸对齐
func (s *FaceRecognizer) CropFaceV2(imageData *SeetaImageData, pointInfo *SeetaPointInfo) *SeetaImageData {
	face := C.facerecognizer_CropFaceV2(s.ptr, imageData.getCStruct(), pointInfo.getCSeetaPointFArray())
	return NewSeetaImageDataFromCStruct(face)
}

// ExtractCroppedFace 对裁剪过的人脸进行特征提取操作
func (s *FaceRecognizer) ExtractCroppedFace(face *SeetaImageData) (bool, []float32) {
	cfeatures := make([]C.float, s.FeatureSize)
	success := int(C.facerecognizer_ExtractCroppedFace(s.ptr, face.getCStruct(), &cfeatures[0])) == 1
	if success {
		features := make([]float32, s.FeatureSize)
		for i := 0; i < s.FeatureSize; i++ {
			features[i] = float32(cfeatures[i])
		}
		return success, features
	}
	return success, nil
}

// Extract 提取人脸特征,从完整图像中提取人脸特征数据
// ！！！pointInfo 必须时标准5点特征定位信息！！！
// 此方法等效CropFaceV2+ExtractCroppedFace
// 返回值 bool代表提取是否成功
// 返回值 []float32为特征数据
func (s *FaceRecognizer) Extract(imageData *SeetaImageData, pointInfo *SeetaPointInfo) (bool, []float32) {
	cfeatures := make([]C.float, s.FeatureSize)
	success := int(C.facerecognizer_Extract(s.ptr, imageData.getCStruct(), pointInfo.getCSeetaPointFArray(), &cfeatures[0])) == 1
	if success {
		features := make([]float32, s.FeatureSize)
		for i := 0; i < s.FeatureSize; i++ {
			features[i] = float32(cfeatures[i])
		}
		return success, features
	}
	return success, nil
}

func (s *FaceRecognizer) Close() {
	C.facerecognizer_free(s.ptr)
}
