#pragma once

#include "CStruct.h"

#ifdef __cplusplus
extern "C"
{
#endif

    #ifdef _WIN32
        #define DLLEXPORT __declspec(dllexport)
    #else
        #define DLLEXPORT
    #endif

    typedef struct facerecognizer
    {
        void *cls;
    } facerecognizer;

    DLLEXPORT facerecognizer *facerecognizer_new(char *model);
    DLLEXPORT void facerecognizer_free(facerecognizer *fr);

    DLLEXPORT void facerecognizer_setProperty(facerecognizer *fr, int property, double value);

    DLLEXPORT double facerecognizer_getProperty(facerecognizer *fr, int property);

    DLLEXPORT int facerecognizer_GetCropFaceWidthV2(facerecognizer *fr);
    DLLEXPORT int facerecognizer_GetCropFaceHeightV2(facerecognizer *fr);
    DLLEXPORT int facerecognizer_GetCropFaceChannelsV2(facerecognizer *fr);

    DLLEXPORT SeetaImageData facerecognizer_CropFaceV2(facerecognizer *fr, const SeetaImageData image, const SeetaPointF *points);
    DLLEXPORT int facerecognizer_ExtractCroppedFace(facerecognizer *fr, const SeetaImageData image, float *features);

    DLLEXPORT int facerecognizer_Extract(facerecognizer *fr, const SeetaImageData image, const SeetaPointF *points, float *features);

    DLLEXPORT int facerecognizer_GetExtractFeatureSize(facerecognizer *fr);

    DLLEXPORT float facerecognizer_CalculateSimilarity(facerecognizer *fr, const float *features1, const float *features2);

#ifdef __cplusplus
}
#endif