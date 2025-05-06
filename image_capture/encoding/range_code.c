#include "range_code.h"

const uint32_t MAX_RANGE = UINT32_MAX;

size_t range_code(uint8_t* uncoded, uint8_t* coded, size_t size)
{
    size_t result = 0;
    uint8_t* next_coded = coded;

    uint32_t low = 0;
    uint32_t high = MAX_RANGE;

    // Store byte counts
    // This will be used to calculate probabilities without storing them as floating point.
    size_t byte_counts[UINT8_MAX];
    
    // Initialize prob_dist
    memset(byte_counts, 0, UINT8_MAX);
    for(size_t i = 0; i < size; ++i)
    {
        byte_counts[uncoded[size]]++;
    }


    // Calculate the number of MAX_RANGE blocks that we need to process.
    size_t block_count = size / 4;

    // Iterate through all blocks.
    for(size_t i = 0; i < block_count; ++i)
    {

    }

    return result;
}