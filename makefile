src_dir=$(shell pwd)
seetaface_dir=$(src_dir)/SeetaFace6Warp

SeetaFace6Warp=SeetaFace6Warp
mac_out=../SeetaFace6Warp/seeta/lib/dylib64

.PHONY:lib

lib:
	echo $(seetaface_dir)/build/include/
	cd SeetaFace6Warp && g++ -std=c++11 *.cpp -fPIC -shared -o $(mac_out)/libSeetaFace6Warp.dylib \
		-I$(seetaface_dir) -I$(seetaface_dir)/seeta \
		-L$(seetaface_dir)/seeta/lib/dylib64 \
		-lSeetaFaceTracking600 -lSeetaFaceDetector600 -lSeetaFaceLandmarker600 -lSeetaFaceRecognizer610 -lSeetaQualityAssessor300

#输出后手动执行
env:
	export CGO_LDFLAGS="-L$(seetaface_dir)/seeta/lib/dylib64 -lSeetaFaceTracking600 -lSeetaQualityAssessor300 -ltennis"
	export DYLD_LIBRARY_PATH="$(seetaface_dir)/seeta/lib/dylib64"

bin:
	cd SeetaFace6Warp/seeta/lib/dylib64 && go build -o main ../../../../test/main.go

# 手动执行 \
cd SeetaFace6Warp/seeta/lib/dylib64 && ./main ../../models ../../../../test/duo6.jpeg
#



#========== MSVC ==========
out_path=../SeetaFace6Warp/bin
seetaface_dir_win=$(sehll cd)
seetaface6_lib_path=D:\code\seetaface6-master\seetaface6-master\build-2\lib\x64
seetaface6_inc_path=$(seetaface_dir_win)

dll:
	cd seetaFace6Warp && cl /EHsc /LD *.cpp /I".\seeta" /Fo$(out_path)\ /link /LIBPATH:".\seeta\lib\x64" \
	SeetaAgePredictor600.lib SeetaEyeStateDetector200.lib SeetaFaceAntiSpoofingX600.lib SeetaFaceDetector600.lib SeetaFaceLandmarker600.lib SeetaFaceRecognizer610.lib \
	SeetaFaceTracking600.lib SeetaGenderPredictor600.lib SeetaMaskDetector200.lib SeetaPoseEstimation600.lib SeetaQualityAssessor300.lib \
	/implib:$(out_path)/SeetaFace6Warp.lib /out:$(out_path)/SeetaFace6Warp.dll


run:
	cd seetaFace6Warp\bin && go build -o main.exe ..\..\test\main.go && main.exe ..\seeta\models ..\..\test\duo6.jpeg

#查看dll导出的函数
echo_export:
	dumpbin /EXPORTS FaceTracker_warp.dll

	# cl /EHsc /LD Seetaface6CGO.cpp /I"./include" /link /LIBPATH:"D:\code\seetaface6-master\seetaface6-master\build-2\lib\x64" SeetaFaceTracking600.lib SeetaFaceDetector600.lib SeetaFaceLandmarker600.lib SeetaFaceRecognizer610.lib SeetaEyeStateDetector200.lib SeetaGenderPredictor600.lib SeetaMaskDetector200.lib SeetaAgePredictor600.lib SeetaQualityAssessor300.lib SeetaFaceAntiSpoofingX600.lib
