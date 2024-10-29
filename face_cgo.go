package seetaFace6go

/*
#cgo darwin CFLAGS: -I./seetaFace6Warp -I./seetaFace6Warp/seeta
#cgo darwin LDFLAGS: -L./seetaFace6Warp/seeta/lib/dylib64 -lSeetaFace6Warp -lSeetaAgePredictor600 -lSeetaEyeStateDetector200 -lSeetaFaceAntiSpoofingX600  -lSeetaFaceDetector600  -lSeetaFaceLandmarker600  -lSeetaFaceRecognizer610  -lSeetaFaceTracking600  -lSeetaGenderPredictor600  -lSeetaMaskDetector200  -lSeetaPoseEstimation600  -lSeetaQualityAssessor300
#cgo windows CFLAGS: -I${SRCDIR}/seetaFace6Warp -I${SRCDIR}/seetaFace6Warp/seeta
#cgo windows LDFLAGS: -L${SRCDIR}/seetaFace6Warp/bin -L${SRCDIR}/seetaFace6Warp/seeta/lib/x64 -lSeetaface6Warp -static
#cgo linux LDFLAGS: -L${SRCDIR}/seetaFace6Warp/seeta/lib/linux_x64 -lSeetaFace6Warp -lSeetaEyeStateDetector200 -lSeetaFaceAntiSpoofingX600  -lSeetaFaceDetector600  -lSeetaFaceLandmarker600  -lSeetaFaceRecognizer610  -lSeetaFaceTracking600  -lSeetaGenderPredictor600   -lSeetaPoseEstimation600 -lSeetaQualityAssessor300
#cgo linux CPPFLAGS: -I${SRCDIR}/seetaFace6Warp -I./seetaFace6Warp/seeta
#cgo linux CXXFLAGS: -std=c++11
#cgo linux,arm64 LDFLAGS: -Wl,--no-as-needed -ldl
*/
import "C"
