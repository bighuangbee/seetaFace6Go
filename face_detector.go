package seetaFace6go

// #include <stdlib.h>
// #include "FaceDetector_warp.h"
import "C"

import (
	"path/filepath"
	"reflect"
	"sort"
	"unsafe"
)

type FaceDetectorProperty int

const (
	FaceDetector_PROPERTY_MIN_FACE_SIZE    FaceDetectorProperty = 0
	FaceDetector_PROPERTY_THRESHOLD        FaceDetectorProperty = 1
	FaceDetector_PROPERTY_MAX_IMAGE_WIDTH  FaceDetectorProperty = 2
	FaceDetector_PROPERTY_MAX_IMAGE_HEIGHT FaceDetectorProperty = 3
	FaceDetector_PROPERTY_NUMBER_THREADS   FaceDetectorProperty = 4
	FaceDetector_PROPERTY_ARM_CPU_MODE     FaceDetectorProperty = 0x101
)

type FaceDetector struct {
	ptr *C.struct_facedetector
}

const (
	_FaceDetector_model = "face_detector.csta"
)

// NewFaceDetector 创建一个人脸检测器
func NewFaceDetector() *FaceDetector {
	cs := C.CString(filepath.Join(_model_base_path, _FaceDetector_model))
	defer C.free(unsafe.Pointer(cs))
	fd := &FaceDetector{
		ptr: C.faceDetector_new(cs),
	}
	fd.SetProperty(FaceDetector_PROPERTY_NUMBER_THREADS, 4)
	return fd
}

func (s *FaceDetector) SetProperty(property FaceDetectorProperty, value float64) {
	C.facedetector_setProperty(s.ptr, C.int(property), C.double(value))
}

func (s *FaceDetector) GetProperty(property FaceDetectorProperty) float64 {
	return float64(C.facedetector_getProperty(s.ptr, C.int(property)))
}

func (s *FaceDetector) Detect(img *SeetaImageData) []*SeetaFaceInfo {
	var result C.struct_SeetaFaceInfoArray = C.facedetector_detect(s.ptr, img.getCStruct())
	var clist []C.struct_SeetaFaceInfo
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&clist))
	arrayLen := int(result.size)
	sliceHeader.Cap = arrayLen
	sliceHeader.Len = arrayLen
	sliceHeader.Data = uintptr(unsafe.Pointer(result.data))

	faceInfoList := make([]*SeetaFaceInfo, arrayLen)
	for i := 0; i < arrayLen; i++ {
		faceInfoList[i] = NewSeetaFaceInfo(clist[i])
	}
	// TODO: c free
	return faceInfoList
}

type _SeetaFaceInfoSlice []*SeetaFaceInfo

func (s _SeetaFaceInfoSlice) Len() int {
	return len(s)
}

func (s _SeetaFaceInfoSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s _SeetaFaceInfoSlice) Less(i, j int) bool {
	return s[i].Postion.GetWidth() > s[j].Postion.GetWidth()
}

func (s *FaceDetector) DetectOrderSize(imageData *SeetaImageData) []*SeetaFaceInfo {
	faces := s.Detect(imageData)
	if len(faces) == 0 {
		return faces
	}
	sort.Sort(_SeetaFaceInfoSlice(faces))
	return faces
}

func (s *FaceDetector) Close() {
	C.facedetector_free(s.ptr)
}
