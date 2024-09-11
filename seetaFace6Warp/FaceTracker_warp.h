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

    __declspec(dllexport) facetracker *facetracker_new(char *model, int video_width, int video_height);
    __declspec(dllexport) void facetracker_free(facetracker *ft);

    __declspec(dllexport) SeetaTrackingFaceInfoArray facetracker_Track(facetracker *ft, const SeetaImageData image);

    __declspec(dllexport) void facetracker_SetMinFaceSize(facetracker *ft, int size);

    __declspec(dllexport) int facetracker_GetMinFaceSize(facetracker *ft);

    __declspec(dllexport) void facetracker_SetThreshold(facetracker *ft, float thresh);

    __declspec(dllexport) float facetracker_GetThreshold(facetracker *ft);

    __declspec(dllexport) void facetracker_SetVideoStable(facetracker *ft, int stable);
    __declspec(dllexport) int facetracker_GetVideoStable(facetracker *ft);

    __declspec(dllexport) void facetracker_SetSingleCalculationThreads(facetracker *ft, int num);

    __declspec(dllexport) void facetracker_SetInterval(facetracker *ft, int interval);
    __declspec(dllexport) void facetracker_Reset(facetracker *ft);
#ifdef __cplusplus
}
#endif
