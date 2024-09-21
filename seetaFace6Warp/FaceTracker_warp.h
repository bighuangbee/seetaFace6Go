#pragma once

#include "CStruct.h"
#include "CTrackingFaceInfo.h"

#ifdef __cplusplus
extern "C"
{
#endif

    #ifdef _WIN32
        #define DLLEXPORT __declspec(dllexport)
    #else
        #define DLLEXPORT
    #endif

    typedef struct facetracker
    {
        void *cls;
    } facetracker;

    DLLEXPORT facetracker *facetracker_new(char *model, int video_width, int video_height);
    DLLEXPORT void facetracker_free(facetracker *ft);

    DLLEXPORT SeetaTrackingFaceInfoArray facetracker_Track(facetracker *ft, const SeetaImageData image);

    DLLEXPORT void facetracker_SetMinFaceSize(facetracker *ft, int size);

    DLLEXPORT int facetracker_GetMinFaceSize(facetracker *ft);

    DLLEXPORT void facetracker_SetThreshold(facetracker *ft, float thresh);

    DLLEXPORT float facetracker_GetThreshold(facetracker *ft);

    DLLEXPORT void facetracker_SetVideoStable(facetracker *ft, int stable);
    DLLEXPORT int facetracker_GetVideoStable(facetracker *ft);

    DLLEXPORT void facetracker_SetSingleCalculationThreads(facetracker *ft, int num);

    DLLEXPORT void facetracker_SetInterval(facetracker *ft, int interval);
    DLLEXPORT void facetracker_Reset(facetracker *ft);
#ifdef __cplusplus
}
#endif