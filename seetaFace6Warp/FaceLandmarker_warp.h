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

    typedef struct facelandmarker
    {
        void *cls;
    } facelandmarker;

    DLLEXPORT facelandmarker *faceLandmarker_new(char *model);
    DLLEXPORT void facelandmarker_free(facelandmarker *fl);
    DLLEXPORT int facelandmarker_number(facelandmarker *fl);
    DLLEXPORT void facelandmarker_mark(facelandmarker *fl, const SeetaImageData image, const SeetaRect face, SeetaPointF *points);
    DLLEXPORT void facelandmarker_mark_mask(facelandmarker *fl, const SeetaImageData image, const SeetaRect face, SeetaPointF *points, int32_t *mask);
#ifdef __cplusplus
}
#endif