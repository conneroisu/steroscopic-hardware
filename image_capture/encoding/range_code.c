#include "range_code.h"

// Constants used for math and data manipulation.
#define MAX_RANGE UINT32_MAX
#define SYMBOL_SIZE ((uint32_t) 256)
#define BITS_PER_SYMBOL ((uint32_t) 8)

// This is the mask used to isolate a symbol, it has not been shifted yet, so ANDing with a value gives the
// least significant symbol.
#define UNSHIFTED_SYMBOL_MASK ((uint64_t)(SYMBOL_SIZE - 1))

// ANDing an integer with this value gives the most significant symbol shifted all the way to the left.
#define MOST_SIG_SYMBOL_MASK ((uint64_t)((SYMBOL_SIZE - 1) << (32 - BITS_PER_SYMBOL)))

// This is the most significant symbol position in a byte. For example, if you have 4 bits per symbol, there
// are two valid symbol locations in the byte: 0 and 1. The constant would be 1 in this case since it is the most
// significant
#define MOST_SIG_SYMBOL_POS_IN_BYTE ((uint64_t)((8 / BITS_PER_SYMBOL) - 1))


typedef struct
{
    size_t current_count;
    size_t previous_count_sum;
} counts_t;

typedef struct
{
    uint64_t low;
    uint64_t high;
} range_t;

int to_positive_int(uint8_t byte)
{
    return (int) ((uint16_t) byte);
}

int max(int a, int b)
{
    return a > b ? a : b;
}


void compute_previous_counts(counts_t* counts)
{
    for(int i = 1; i < SYMBOL_SIZE; ++i)
    {
        counts[i].previous_count_sum = counts[i - 1].previous_count_sum + counts[i - 1].current_count;
    }
}

// prob_reduction is in units of percents
void adjust_probabilities(counts_t* counts, size_t size, uint8_t byte, int prob_reduction)
{
    // So we don't have to worry about negatives later, convert byte to a postiive integer.
    int i_byte = to_positive_int(byte);

    // Calculate how many entries to the left and right of "byte" we will be increasing.
    int left_entries = byte < 8 ? byte : 8;
    int right_entries = (byte + 8) > 255 ? 255 - byte : 8;

    size_t percent_count = ((size / 100) * prob_reduction) / (256 - (left_entries + right_entries));

    // 0.01% is the min
    size_t min_percent_count = size / (10000);

    // Attempt to reduce all the probabilities that are not within 8 of "byte" by "percent_count"
    // However, probabilities that are not zero should not be brought down to zero and
    // we need to keep track of how many counts we reduced so we can increase the counts 8 away from
    // "byte" to keep the probability dist 1.

    size_t reduction = 0;

    // Attempt to reduce the probabilities of the symbols that are NOT 8 away from "byte"
    for(int i = 0; i < SYMBOL_SIZE; ++i)
    {
        if(i < (i_byte - 8) || i >(i_byte + 8))
        {
            size_t count = counts[i].current_count;

            if(count < min_percent_count)
            {
                continue;
            }
            else if(count < percent_count)
            {
                counts[i].current_count = min_percent_count;
                reduction += (count - min_percent_count);
            }
            else
            {
                counts[i].current_count -= percent_count;
                reduction += percent_count;
            }
        }
    }

    size_t count_increase = reduction / (left_entries + right_entries);
    size_t extra = reduction % (left_entries + right_entries);

    for(int i = max(i_byte - 8, 0); (i <= i_byte + 8) && (i < SYMBOL_SIZE); ++i)
    {
        if(i >= 0)
        {
            if(i == i_byte)
            {
                counts[i].current_count += extra;
            }
            else
            {
                counts[i].current_count += count_increase;
            }
        }
    }

    // Recompute previous counts
    compute_previous_counts(counts);
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
    than (SYMBOL_SIZE * adjustment_threshold). This makes it so our range is sufficiently large
    so we don't run out of resolution when doing the probability calculations. adjustement_factor
    scales with the data size and is determined experimentally. Setting a large value (>= 32) never
    hurts, however, the larger the value, the more bytes you encode. DECODER MUST USE THE SAME VALUE!!!
*/
size_t range_code(uint8_t* uncoded, uint8_t* coded, size_t size, int adjustment_threshold)
{
    // Keeps track of how many symbols we have put into a byte to be shifted off.
    int byte_pack_symbol_count = 0;
    uint8_t byte_to_shift = 0;

    // Size is in byte, calculate the number of symbols using size.
    size_t total_symbol_count = (8 / BITS_PER_SYMBOL) * size;

    // Don't record byte if it is a part of the leading zeros.
    int leading_zeros = 1;

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
            int symbol = (byte & (UNSHIFTED_SYMBOL_MASK << (BITS_PER_SYMBOL * j))) >> (BITS_PER_SYMBOL * j);
            symbol_counts[symbol].current_count++;
        }
    }

    // Calculate the previous counts
    compute_previous_counts(symbol_counts);

    // Iterate through all bytes
    for(size_t i = 0; i < size; ++i)
    {

        size_t range_size = (range.high - range.low) + 1;
        // Iterate through each symbol in the byte.
        // Most significant symbols should be first.
        for(int j = (8 / BITS_PER_SYMBOL) - 1; j >= 0; --j)
        {
            // Calculate the next range
            range_size = (range.high - range.low) + 1;

            int symbol = (*next_uncoded & (UNSHIFTED_SYMBOL_MASK << (BITS_PER_SYMBOL * j))) >> (BITS_PER_SYMBOL * j);

            counts_t count = symbol_counts[symbol];

            range.low += (uint64_t) ((count.previous_count_sum * range_size) / total_symbol_count);
            range.high = ((uint64_t) ((count.current_count * range_size) / total_symbol_count)) + range.low;


            // Shift off symbols if we can.
            /*
                Edge Case:
                If the range is less than adjustment_threshold and the most significant symbol on the low and high
                parts of the range don't match, then we need to adjust the range such that they do match so they can
                be shifted off. Otherwise, the range gets too small and the low and high range values will eventually equal, which is no good!
            */
            uint64_t low_symbol_value = range.low & MOST_SIG_SYMBOL_MASK;
            uint64_t high_symbol_value = range.high & MOST_SIG_SYMBOL_MASK;
            range_size = (range.high - range.low) + 1;
            while(low_symbol_value == high_symbol_value || (range_size < adjustment_threshold && low_symbol_value != high_symbol_value))
            {
                if(range_size < adjustment_threshold && low_symbol_value != high_symbol_value)
                {
                    uint64_t new_high = (low_symbol_value + (1 << (32 - BITS_PER_SYMBOL))) - 1;
                    uint64_t delta = range.high - new_high;
                    range.low -= delta;
                    range.high = new_high;
                }
                else
                {
                    uint8_t shifted_symbol_value = (uint8_t) (low_symbol_value >> (32 - BITS_PER_SYMBOL));
                    byte_to_shift |= (shifted_symbol_value << ((MOST_SIG_SYMBOL_POS_IN_BYTE - byte_pack_symbol_count) * BITS_PER_SYMBOL));
                    range.low = (range.low & ~MOST_SIG_SYMBOL_MASK) << BITS_PER_SYMBOL;
                    range.high = (range.high & ~MOST_SIG_SYMBOL_MASK) << BITS_PER_SYMBOL;

                    // If we have packed the byte, shift it off.
                    if(byte_pack_symbol_count == MOST_SIG_SYMBOL_POS_IN_BYTE)
                    {
                        if(leading_zeros && byte_to_shift)
                        {
                            leading_zeros = 0;
                        }

                        *next_coded = byte_to_shift;
                        next_coded++;

                        if(!leading_zeros)
                        {
                            result += 8;
                        }

                        byte_to_shift = 0;
                        byte_pack_symbol_count = 0;
                    }
                    else
                    {
                        byte_pack_symbol_count++;
                    }
                }

                low_symbol_value = range.low & MOST_SIG_SYMBOL_MASK;
                high_symbol_value = range.high & MOST_SIG_SYMBOL_MASK;
                range_size = (range.high - range.low) + 1;

                // If the ranges ever equal, we did not have enough resolution!
                // The adjustment factor must be increased!
                if(range.high == range.low)
                {
                    return 0;
                }
            }
        }
        // Adjust every 128 pixels
        if(!(i % 128))
        {
            adjust_probabilities(symbol_counts, size, *next_uncoded, 10);
        }
        next_uncoded++;

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

    return result / 8;
}