#ifndef RANGE_CODE_H__
#define RANGE_CODE_H__
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>


// Outputs the number of BITS that make up the coded data
// We are not guaranteed that the coded data will be byte aligned.
size_t range_code(uint8_t* uncoded, uint8_t* coded, size_t size, int adjustment_factor);

#endif