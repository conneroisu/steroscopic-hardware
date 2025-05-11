#include "range_code.h"
#define MAX_RANGE UINT32_MAX
#define SYMBOL_SIZE ((uint8_t) 16)
#define BITS_PER_SYMBOL ((uint8_t) 4)
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

// Gets either byte 0, 1, 2, or 3
int find_set_msb(uint32_t value)
{
    for(int i = 3; i >= 0; --i)
    {
        if(value & 0xFF << (i * 8))
        {
            return i;
        }
    }

    return 0;
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

    new_low_n = low_n-1 + Sum(prob_0 to prob_n-1) * prev_seq_range_size
    new_high_n = new_low_n + prob_n * prev_seq_range_size

    In addition, since we are dealing with a large data set, there will be basically an uncountable number of possible sequences,
    which results in a huge starting range. To avoid this, we limit the starting range to the size of an unsigned 32-bit integer
    and work in blocks of that size, shifting off known digits of the final range values.

    The range will be adjusted such that a byte can be shifted off if the range is less
    than (SYMBOL_SIZE * adjustment_factor). This makes it so our range is sufficiently large
    so we don't run out of resolution when doing the probability calculations. adjustement_factor
    scales with the data size and is determined experimentally. Setting a large value (>= 32) never
    hurts, however, the larger the value, the more bytes you encode. DECODER MUST USE THE SAME VALUE!!!
*/
size_t range_code(uint8_t* uncoded, uint8_t* coded, size_t size, int adjustment_factor)
{
    size_t result = 0;
    uint8_t* next_coded = coded;
    uint8_t* next_uncoded = uncoded;

    range_t range = { .low = 0, .high = MAX_RANGE };

    // Store byte counts
    // This will be used to calculate probabilities without storing them as floating point.
    counts_t symbol_counts[SYMBOL_SIZE];

    // Zeroize the counts
    memset(symbol_counts, 0, sizeof(counts_t) * SYMBOL_SIZE);

    // Count the bytes
    for(size_t i = 0; i < size; ++i)
    {
        // Split byte up into symbols.
        for(int j = 0; j < (8 / BITS_PER_SYMBOL); ++j)
        {
            uint8_t byte = uncoded[i];
            int symbol = (byte & ((SYMBOL_SIZE - 1) << (BITS_PER_SYMBOL * j))) >> (BITS_PER_SYMBOL * j);
            symbol_counts[symbol].current_count++;
        }
    }

    // Calculate the previous counts
    for(int i = 1; i < SYMBOL_SIZE; ++i)
    {
        symbol_counts[i].previous_count_sum = symbol_counts[i - 1].previous_count_sum + symbol_counts[i - 1].current_count;
    }

    size_t total_symbol_count = symbol_counts[SYMBOL_SIZE - 1].previous_count_sum + symbol_counts[SYMBOL_SIZE - 1].current_count;

    // Iterate through all bytes
    for(size_t i = 0; i < size; ++i)
    {

        size_t range_size = range.high - range.low;
        // Iterate through each symbol in the byte.
        for(int j = 0; j < (8 / BITS_PER_SYMBOL); ++j)
        {
            int symbol = (*next_uncoded & ((SYMBOL_SIZE - 1) << (BITS_PER_SYMBOL * j))) >> (BITS_PER_SYMBOL * j);

            // Calculate the next range
            range_size = range.high - range.low;

            range.low += (uint32_t) ((symbol_counts[symbol].previous_count_sum * range_size) / total_symbol_count);
            range.high = ((uint32_t) ((symbol_counts[symbol].current_count * range_size) / total_symbol_count)) + range.low;

            // If the ranges ever equal, we did not have enough resolution!
            // The adjustment factor must be increased!
            if(range.high == range.low)
            {
                return 0;
            }
        }
        next_uncoded++;

        range_size = range.high - range.low;

        // Only shift off bytes when we need to
        if(range_size < (SYMBOL_SIZE * adjustment_factor))
        {
            // Emit digits. Max of 3.
            for(int k = 0; k < 3; ++k)
            {

                // See if we can shift off bytes.
                int low_set_msb_loc = find_set_msb(range.low);
                int high_set_msb_loc = find_set_msb(range.high);
                uint32_t low_byte_value = range.low & (0xFF << (low_set_msb_loc * 8));
                uint32_t high_byte_value = range.high & (0xFF << (high_set_msb_loc * 8));

                /*
                    Edge Case:
                    If the range is less than (SYMBOL_COUNT * 32) and the SET low and high MSBs don't match, we need
                    to shift the range such that they do match. Otherwise, we run out of range to get the resolution we need.
                    This is done by shifting BOTH low and high until high is the last possible number with the MSB equal to low's MSB.
                */
                range_size = range.high - range.low;
                if(low_byte_value != high_byte_value)
                {
                    uint32_t new_high = (low_byte_value + (1 << (8 * low_set_msb_loc))) - 1;
                    uint32_t delta = range.high - new_high;
                    range.low -= delta;
                    range.high = new_high;
                }

                if(low_byte_value > range_size && low_set_msb_loc == high_set_msb_loc && low_byte_value == high_byte_value)
                {
                    *next_coded = (uint8_t) (low_byte_value >> (low_set_msb_loc * 8));
                    next_coded++;
                    result += 8;

                    range.low = range.low << ((4 - low_set_msb_loc) * 8);
                    range.high = range.high << ((4 - low_set_msb_loc) * 8);
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
        byte_to_write |= (bit_index << ((uint8_t) low_bit_value & 0x1));
        result++;

        low_bit_value = (low_byte & (0x1 << bit_index)) == 1;
        high_bit_value = (low_byte & (0x1 << bit_index)) == 1;
    }

    return result;
}