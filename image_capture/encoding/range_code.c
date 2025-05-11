#include "range_code.h"
#define MAX_RANGE UINT32_MAX
#define SYMBOL_COUNT 256
#include <stdio.h>

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

    new_low_n = low_n-1 + Sum(prob_0 to prob_n-1) * prev_seq_range_size
    new_high_n = new_low_n + prob_n * prev_seq_range_size

    In addition, since we are dealing with a large data set, there will be basically an uncountable number of possible sequences,
    which results in a huge starting range. To avoid this, we limit the starting range to the size of an unsigned 32-bit integer
    and work in blocks of that size, shifting off known digits of the final range values.
*/
size_t range_code(uint8_t* uncoded, uint8_t* coded, size_t size)
{
    size_t result = 0;
    uint8_t* next_coded = coded;
    uint8_t* next_uncoded = uncoded;

    range_t range = {.low = 0, .high = MAX_RANGE};

    // Store byte counts
    // This will be used to calculate probabilities without storing them as floating point.
    counts_t byte_counts[SYMBOL_COUNT];
    
    // Zeroize the counts
    memset(byte_counts, 0, sizeof(counts_t) * SYMBOL_COUNT);

    // Count the bytes
    for(size_t i = 0; i < size; ++i)
    {
        byte_counts[uncoded[i]].current_count++;
    }

    // Calculate the previous counts
    for(int i = 1; i < SYMBOL_COUNT; ++i)
    {
        byte_counts[i].previous_count_sum = byte_counts[i - 1].previous_count_sum + byte_counts[i - 1].current_count;
    }

    for(size_t i = 0; i < size; ++i)
    {
        
        // Calculate the next range
        size_t prev_range_size = range.high - range.low;

        if(prev_range_size < 1024)
        {
            printf("Oh no!!!\n");
        }

        range.low += (uint32_t) ((byte_counts[*next_uncoded].previous_count_sum * prev_range_size) / size);
        range.high = ((uint32_t) ((byte_counts[*next_uncoded].current_count * prev_range_size) / size)) + range.low;
        next_uncoded++;

        prev_range_size = range.high - range.low;

        // Emit digits. Max of 3.
        for(int j = 0; j < 3; ++j)
        {

            // See if the upper byte is all zeros, if so, shift it off BUT DON'T RECORD THE ZERO!
            // This is just to expand our range whenever we can!
            if(!(range.low & 0xFF000000) && !(range.high & 0xFF000000))
            {
                range.low = range.low << 8;
                range.high = range.high << 8;
            }
            else
            {
                // See if we can shift off bytes.
                uint32_t low_byte_value = range.low & 0xFF000000;
                uint32_t high_byte_value = range.high & 0xFF000000;

                if(low_byte_value > prev_range_size && low_byte_value == high_byte_value)
                {
                    *next_coded = (uint8_t) (low_byte_value >> 24);
                    next_coded++;
                    result += 8;

                    range.low = range.low << 8;
                    range.high = range.high << 8;
                }
            }
        }

    }

    // Emit the remaining bits.
    uint8_t low_byte = (uint8_t) ((range.low & 0xFF000000) >> 24);
    uint8_t high_byte = (uint8_t) ((range.high & 0xFF000000) >> 24);
    int bit_index = 7;
    int low_bit_value = (low_byte & (0x1 << bit_index)) != 0;
    int high_bit_value = (high_byte & (0x1 << bit_index)) != 0;
    uint8_t byte_to_write = 0;
    
    while(high_bit_value == low_bit_value && bit_index >= 0)
    {
        bit_index--;
        byte_to_write |= (bit_index << ((uint8_t)low_bit_value & 0x1));
        result++;

        low_bit_value = (low_byte & (0x1 << bit_index)) == 1;
        high_bit_value = (low_byte & (0x1 << bit_index)) == 1;
    }

    return result;
}