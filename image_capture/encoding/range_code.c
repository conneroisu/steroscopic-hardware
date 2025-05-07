#include "range_code.h"

typedef struct
{
    size_t current_count;
    size_t previous_count_sum;
} counts_t;

typedef struct
{
    uint32_t low;
    uint32_t high;
} range_t;

const uint32_t MAX_RANGE = UINT32_MAX;

range_t get_next_range(uint8_t symbol, range_t current_range, counts_t** counts, size_t size)
{
    range_t result = {.low = 0, .high = 0};

    size_t prev_range_size = current_range.high - current_range.low;

    result.low = (uint32_t) (((*counts)[symbol].previous_count_sum * prev_range_size) / size);
    result.high = ((uint32_t) (((*counts)[symbol].current_count * prev_range_size) / size)) + result.low;

    return result;
}

/* 
    Basic Description:
    Range encoding uses integer ranges whose range size is proportional to the probability
    that the symbol sequences associated with it will occur given the probabilities of each symbol,
    which are determined through counting.

    The ranges are ordered by the symbol order, which is 0 to 255. Every iteration,
    the probability range for the new sequence with the next symbol that we want to encode
    is generated. One might first think that you need the previous range values to calculate the new
    range, but this is not the case. When calculating the new "low" and "high" range values for new ranges,
    the symbol probability is multiplied by the previous sequence's range size. Then, the next sequences
    do the same thing but the "low" part the range is the previous sequences "high" part. However,
    since the probabilities are all multiplied by the previous sequence range size, we can figure out
    the "low" and "high" range values by simply multiplying the previous sequence's range size by the
    sum of prior probabilities and the current symbol probability. This results in the following equation:

    new_low_n = Sum(prob_0 to prob_n-1) * prev_seq_range_size
    new_high_n = new_low_n + prob_n * prev_seq_range_size

    In addition, since we are dealing with a large data set, there will be basically an uncountable number of possible sequences,
    which results in a huge starting range. To avoid this, we limit the starting range to the size of an unsigned 32-bit integer
    and work in blocks of that size, shifting off known digits of the final range values.
*/
size_t range_code(uint8_t* uncoded, uint8_t* coded, size_t size)
{
    size_t result = 0;
    uint8_t* next_coded = coded;

    range_t range = {.low = 0, .high = MAX_RANGE};

    // Store byte counts
    // This will be used to calculate probabilities without storing them as floating point.
    counts_t byte_counts[UINT8_MAX];
    
    // Zeroize the counts
    memset(byte_counts, 0, sizeof(counts_t) * UINT8_MAX);

    // Count the bytes
    for(size_t i = 0; i < size; ++i)
    {
        byte_counts[uncoded[size]].current_count++;
    }

    // Calculate the previous counts
    for(int i = 1; i < UINT8_MAX; ++i)
    {
        byte_counts[i].previous_count_sum = byte_counts[i - 1].previous_count_sum + byte_counts[i - 1].current_count;
    }

    // Calculate the number of MAX_RANGE blocks that we need to process.
    size_t block_count = size / 4;

    // Iterate through all blocks.
    for(size_t i = 0; i < block_count; ++i)
    {

    }

    return result;
}