src_dir=$(shell pwd)
seetaface_dir=$(src_dir)/../SeetaFace6Open

SeetaFace6Warp=SeetaFace6Warp
out_path=../SeetaFace6Warp/bin

.PHONY:lib

lib:
	echo $(seetaface_dir)/build/include/
	cd SeetaFace6Warp && g++ -std=c++11 *.cpp -fPIC -shared -o $(out_path)/$(SeetaFace6Warp).so \
		-I$(seetaface_dir)/build/include/ \
		-L$(seetaface_dir)/build/lib64/ \
		-lSeetaFaceTracking600 -lSeetaFaceDetector600 -lSeetaFaceLandmarker600 -lSeetaFaceRecognizer610

#输出后手动执行
env:
	export CGO_LDFLAGS="-L$(src_dir)/libSf6 -L$(seetaface_dir)/build/lib64 -lSeetaFaceTracking600 -ltennis"
	export DYLD_LIBRARY_PATH="$(src_dir)/libSf6:$(seetaface_dir)/build/lib64"

bin:
	cd test && go build -v



#========== MSVC ==========
seetaface_dir_win=$(sehll cd)
seetaface6_lib_path=D:\code\seetaface6-master\seetaface6-master\build-2\lib\x64
seetaface6_inc_path=$(seetaface_dir_win)

dll:
	cd seetaFace6Warp && cl /EHsc /LD *.cpp /I".\seeta" /Fo$(out_path)\ /link /LIBPATH:".\seeta\lib\x64" \
	SeetaAgePredictor600.lib SeetaEyeStateDetector200.lib SeetaFaceAntiSpoofingX600.lib SeetaFaceDetector600.lib SeetaFaceLandmarker600.lib SeetaFaceRecognizer610.lib \
	SeetaFaceTracking600.lib SeetaGenderPredictor600.lib SeetaMaskDetector200.lib SeetaPoseEstimation600.lib SeetaQualityAssessor300.lib \
	/implib:$(out_path)/SeetaFace6Warp.lib /out:$(out_path)/SeetaFace6Warp.dll


run:
	cd seetaFace6Warp\bin && go build -o main.exe ..\..\test\main.go && main.exe ..\..\test\duo6.jpeg

#查看dll导出的函数
echo_export:
	dumpbin /EXPORTS FaceTracker_warp.dll

	# cl /EHsc /LD Seetaface6CGO.cpp /I"./include" /link /LIBPATH:"D:\code\seetaface6-master\seetaface6-master\build-2\lib\x64" SeetaFaceTracking600.lib SeetaFaceDetector600.lib SeetaFaceLandmarker600.lib SeetaFaceRecognizer610.lib SeetaEyeStateDetector200.lib SeetaGenderPredictor600.lib SeetaMaskDetector200.lib SeetaAgePredictor600.lib SeetaQualityAssessor300.lib SeetaFaceAntiSpoofingX600.lib
