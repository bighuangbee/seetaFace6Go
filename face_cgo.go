package seetaFace6go

/*
#cgo linux CPPFLAGS: -I./include
#cgo linux LDFLAGS: -lSeetaAgePredictor600 -lSeetaEyeStateDetector200 -lSeetaFaceAntiSpoofingX600  -lSeetaFaceDetector600  -lSeetaFaceLandmarker600  -lSeetaFaceRecognizer610  -lSeetaFaceTracking600  -lSeetaGenderPredictor600  -lSeetaMaskDetector200  -lSeetaPoseEstimation600  -lSeetaQualityAssessor300
#cgo linux CXXFLAGS: -std=c++11
#cgo linux LDFLAGS: -Wl,-rpath,\$ORIGIN/lib:\$ORIGIN/libs -Wl,--disable-new-dtags
#cgo linux,arm64 LDFLAGS: -Wl,--no-as-needed -ldl
#cgo windows CFLAGS: -I${SRCDIR}/seetaFace6Warp -I${SRCDIR}/seetaFace6Warp/seeta
#cgo windows LDFLAGS: -L${SRCDIR}/seetaFace6Warp/bin -L${SRCDIR}/seetaFace6Warp/seeta/lib/x64 -lSeetaface6Warp -static
*/
import "C"

// #cgo windows LDFLAGS: -L${SRCDIR}/lib -lSeetaface6Warp -static
