src_dir=$(shell pwd)
seetaface_dir=$(src_dir)/../SeetaFace6Open

SeetaFace6Warp=SeetaFace6Warp
out_path=../lib

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
seetaface6_inc_path=$(seetaface_dir_win)\seetaFace6Warp

dll:
	echo $(seetaface_dir_win)
	cd seetaFace6Warp && cl /EHsc /LD *.cpp /I$(seetaface6_inc_path) /link /LIBPATH:$(seetaface6_lib_path) \
	SeetaFaceTracking600.lib SeetaFaceLandmarker600.lib SeetaFaceDetector600.lib SeetaFaceRecognizer610.lib \
	/out:$(out_path)/$(SeetaFace6Warp).dll /implib:$(out_path)/$(SeetaFace6Warp).lib



run:
	cd test && go build -o ../lib/test.exe . && ..\lib\test.exe -path1 duo6.jpeg -path2 duo5.jpeg

#查看dll导出的函数
echo_export:
	dumpbin /EXPORTS FaceTracker_warp.dll

	# cl /EHsc /LD Seetaface6CGO.cpp /I"./include" /link /LIBPATH:"D:\code\seetaface6-master\seetaface6-master\build-2\lib\x64" SeetaFaceTracking600.lib SeetaFaceDetector600.lib SeetaFaceLandmarker600.lib SeetaFaceRecognizer610.lib SeetaEyeStateDetector200.lib SeetaGenderPredictor600.lib SeetaMaskDetector200.lib SeetaAgePredictor600.lib SeetaQualityAssessor300.lib SeetaFaceAntiSpoofingX600.lib
