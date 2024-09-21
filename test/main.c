#include <stdio.h>
#include <stdlib.h>
#include "CStruct.h"
#include "CFaceInfo.h"

// 包含生成的人脸检测器接口
#include "facedetector.h"

int main(int argc, char **argv) {
    if (argc < 2) {
        printf("Usage: %s <image_file>\n", argv[0]);
        return 1;
    }

    // 加载模型路径
    char *modelPath = "path/to/your/model"; // 替换成实际模型文件路径

    // 初始化人脸检测器
    facedetector *fd = faceDetector_new(modelPath);
    if (fd == NULL) {
        printf("Failed to initialize face detector\n");
        return 1;
    }

    // 加载图像（模拟操作）
    // 这里假设你有一个函数将图像加载到 SeetaImageData 结构体
    SeetaImageData image;
    // 模拟加载图像, 需要自行实现图像加载逻辑
    // loadImage(argv[1], &image);

    // 检测人脸
    SeetaFaceInfoArray faces = facedetector_detect(fd, image);
    if (faces.size > 0) {
        printf("Detected %d face(s)\n", faces.size);
        for (int i = 0; i < faces.size; ++i) {
            SeetaRect face = faces.data[i].pos;
            printf("Face %d: [x=%d, y=%d, width=%d, height=%d]\n", 
                   i + 1, face.x, face.y, face.width, face.height);
        }
    } else {
        printf("No faces detected\n");
    }

    // 设置属性（例如设置检测器的阈值）
    facedetector_setProperty(fd, 0, 0.8); // 示例，具体属性请根据库文档调整

    // 获取属性
    double threshold = facedetector_getProperty(fd, 0);
    printf("Face detector threshold: %f\n", threshold);

    // 释放人脸检测器资源
    facedetector_free(fd);

    return 0;
}
