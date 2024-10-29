cgoWarpName=SeetaFace6Warp

src_mac=$(shell pwd)/$(cgoWarpName)
out_mac=$(src_mac)/seeta/lib/dylib64
out_linux=$(src_mac)/seeta/lib/linux_x64

.PHONY:mac_warp

mac_warp:
	cd $(src_mac) && g++ -std=c++11 *.cpp -fPIC -shared -o $(out_mac)/libSeetaFace6Warp.dylib \
	-I$(src_mac) -I$(src_mac)/seeta \
	-L$(out_mac) \
	-lSeetaFaceTracking600 -lSeetaFaceDetector600 -lSeetaFaceLandmarker600 -lSeetaFaceRecognizer610 -lSeetaQualityAssessor300

linux_warp:
	cd $(src_mac) && g++ -std=c++11 *.cpp -fPIC -shared -o $(out_linux)/libSeetaFace6Warp.so \
	-I$(src_mac) -I$(src_mac)/seeta \
	-L$(out_linux) \
	-lSeetaFaceTracking600 -lSeetaFaceDetector600 -lSeetaFaceLandmarker600 -lSeetaFaceRecognizer610 -lSeetaQualityAssessor300


#输出后手动执行
mac_env:
	export CGO_LDFLAGS="-L$(out_mac)"
	export DYLD_LIBRARY_PATH="$(out_mac)"

mac_bin:
	cd $(out_mac) && go build -o mac_bin ../../../../test/main.go
# 手动执行 \
cd SeetaFace6Warp/seeta/lib/dylib64 && ./mac_bin ../../models ../../../../test/duo6.jpeg
#




#========== Windows ==========
src_win=$(shell cd)/$(cgoWarpName)
out_path=$(src_win)/bin
seetaface6_lib_path=D:\code\seetaface6-master\seetaface6-master\build-2\lib\x64

#========== MSVC ========== x86_x64 Cross Tools Command Prompt for VS 2022
win_warp:
	cd $(src_win) && cl /EHsc /LD *.cpp /I".\seeta" /Fo$(out_path)\ /link /LIBPATH:".\seeta\lib\x64" \
	SeetaAgePredictor600.lib SeetaEyeStateDetector200.lib SeetaFaceAntiSpoofingX600.lib SeetaFaceDetector600.lib SeetaFaceLandmarker600.lib SeetaFaceRecognizer610.lib \
	SeetaFaceTracking600.lib SeetaGenderPredictor600.lib SeetaMaskDetector200.lib SeetaPoseEstimation600.lib SeetaQualityAssessor300.lib \
	/implib:$(out_path)/SeetaFace6Warp.lib /out:$(out_path)/SeetaFace6Warp.dll

#========== MinGW ==========
win_run:
	cd $(out_path) && go build -o win-bin.exe ..\..\test\main.go && win-bin.exe ..\seeta\models ..\..\test\duo6.jpeg

#查看dll导出的函数
echo_export:
	dumpbin /EXPORTS FaceTracker_warp.dll

	# cl /EHsc /LD Seetaface6CGO.cpp /I"./include" /link /LIBPATH:"D:\code\seetaface6-master\seetaface6-master\build-2\lib\x64" SeetaFaceTracking600.lib SeetaFaceDetector600.lib SeetaFaceLandmarker600.lib SeetaFaceRecognizer610.lib SeetaEyeStateDetector200.lib SeetaGenderPredictor600.lib SeetaMaskDetector200.lib SeetaAgePredictor600.lib SeetaQualityAssessor300.lib SeetaFaceAntiSpoofingX600.lib
