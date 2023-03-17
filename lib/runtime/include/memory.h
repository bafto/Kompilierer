/*
    declares functions and macros to work
    with memory (allocation, freeing, etc.)
*/
#ifndef DDP_MEMORY_H
#define DDP_MEMORY_H

#include "common.h"

// used for allocation/reallocation and freeing of memory
// to allocate call reallocate(NULL, 0, size)
// to free call reallocate(ptr, oldsize, 0)
// to reallocate call reallocate(ptr, oldsize, newsize)
void* ddp_reallocate(void* pointer, size_t oldSize, size_t newSize);

// helper macro to allocate a specific amount of objects
#define ALLOCATE(type, count) \
    (type*)ddp_reallocate(NULL, 0, sizeof(type) * (count))

// helper macro to free any type (not arrays though)
#define FREE(type, pointer) ddp_reallocate(pointer, sizeof(type), 0)

// helper macro to calculate the new capacity of an array
#define GROW_CAPACITY(capacity) \
    ((capacity) < 8 ? 8 : (capacity) * 2)

// helper macro to expand the capacity of an array
#define GROW_ARRAY(type, pointer, oldCount, newCount) \
    (type*)ddp_reallocate(pointer, sizeof(type) * (oldCount), \
        sizeof(type) * (newCount))

// helper to free a whole array
#define FREE_ARRAY(type, pointer, oldCount) \
    ddp_reallocate(pointer, sizeof(type) * (oldCount), 0)

#endif // DDP_MEMORY_H