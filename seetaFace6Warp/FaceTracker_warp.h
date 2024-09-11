#pragma once

#include "CStruct.h"

#ifdef __cplusplus
extern "C"
{
#endif

    typedef struct SeetaFaceInfo
    {
        SeetaRect pos;
        float score;
    } SeetaFaceInfo;

    typedef struct SeetaTrackingFaceInfo
    {
        SeetaRect pos;
        float score;
        int frame_no;
        int PID;
        int step;
    } SeetaTrackingFaceInfo;

    typedef struct SeetaTrackingFaceInfoArray
    {
        struct SeetaTrackingFaceInfo *data;
        int size;
    } SeetaTrackingFaceInfoArray;

    typedef struct facetracker
    {
        void *cls;
    } facetracker;

    // Conditionally define dllexport for Windows
    #ifdef _WIN32
        #define DLLEXPORT __declspec(dllexport)
    #else
        #define DLLEXPORT
    #endif

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
