#pragma once

#include "CStruct.h"
#include "CFaceInfo.h"

#ifdef __cplusplus
extern "C"
{
#endif
    // Conditionally define dllexport for Windows
    #ifdef _WIN32
        #define DLLEXPORT __declspec(dllexport)
    #else
        #define DLLEXPORT
    #endif

    typedef struct facedetector
    {
        void *cls;
    } facedetector;

    DLLEXPORT facedetector *faceDetector_new(char *model);

    DLLEXPORT SeetaFaceInfoArray facedetector_detect(facedetector *fd, const SeetaImageData image);

    DLLEXPORT void facedetector_free(facedetector *fd);

    DLLEXPORT void facedetector_setProperty(facedetector *fd, int property, double value);

    DLLEXPORT double facedetector_getProperty(facedetector *fd, int property);

#ifdef __cplusplus
}
#endif
