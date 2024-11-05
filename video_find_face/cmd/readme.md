export CGO_LDFLAGS="-L/Users/bighuangbee/code/custom/face/seetaFace6Go/SeetaFace6Warp/seeta/lib/dylib64"
export DYLD_LIBRARY_PATH="/Users/bighuangbee/code/custom/face/seetaFace6Go/SeetaFace6Warp/seeta/lib/dylib64"

/*
-Wl,-rpath=/root/face_recognize/recognize/libs/face_gpu/sdk/lib/ -lhiar_cluster
export CGO_LDFLAGS="-Wl,-rpath=/hiar_face/seetaFace6Go/seetaFace6Warp/seeta/lib/linux_x64/ -L/hiar_face/seetaFace6Go/seetaFace6Warp/seeta/lib/linux_x64/ -lSeetaFace6Warp -lSeetaEyeStateDetector200 -lSeetaFaceAntiSpoofingX600  -lSeetaFaceDetector600  -lSeetaFaceLandmarker600  -lSeetaFaceRecognizer610  -lSeetaFaceTracking600  -lSeetaGenderPredictor600   -lSeetaPoseEstimation600 -lSeetaQualityAssessor300"
export LD_LIBRARY_PATH=/hiar_face/seetaFace6Go/seetaFace6Warp/seeta/lib/linux_x64/:$LD_LIBRARY_PATH
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_recognize/recognize/libs/face_gpu/sdk/lib/

GNU 9.4.0

export CGO_CXXFLAGS="-I/usr/local/include/opencv4"
export CGO_CFLAGS="-I/usr/local/include/opencv4"
export CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_imgproc -lopencv_highgui -lopencv_videoio -lopencv_imgcodecs -lopencv_objdetect -lopencv_features2d -lopencv_video -lopencv_dnn -lopencv_calib3d"

*/

//go run . "rtsp://admin:Ab123456.@192.168.1.108:554/cam/realmonitor?channel=1&subtype=0"
