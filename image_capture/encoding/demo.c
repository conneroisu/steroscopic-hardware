#include "range_code.h"
#include <stdio.h>
#include <errno.h>

#define TEST_MODE 1

#define SIZE 2073600U

int main()
{

#if TEST_MODE == 1
    uint8_t* data = calloc(SIZE, sizeof(uint8_t));
    uint8_t* coded = calloc(SIZE, sizeof(uint8_t));

    size_t result = 0;

    // Load the image.
    FILE* image = fopen("image.bin", "r");

    if(image == NULL)
    {
        perror("Could not open image!");
        return errno;
    }

    fread(data, sizeof(uint8_t), SIZE, image);

    if(fclose(image) != 0)
    {
        perror("Could not close image file!");
        return errno;
    }

    result = range_code(data, coded, SIZE, 32768);


    free(data);
    free(coded);
    printf("Result: %ld\n", result);

#elif TEST_MODE == 2
    uint8_t* data = calloc(1000, sizeof(uint8_t));
    uint8_t* coded = calloc(1000, sizeof(uint8_t));

    memset(coded, 0, 1000);

    for(int i = 0; i < 1000; ++i)
    {
        data[i] = i;
    }

    size_t result = range_code(data, coded, 1000, 64);

    free(data);
    free(coded);
    printf("Result: %ld\n", result);

#else
    uint8_t data[16] = { 130, 55, 39, 55, 130, 72, 72, 9, 72, 8, 80, 76, 125, 130, 72, 9 };
    uint8_t coded[16];

    memset(coded, 0, 16);

    size_t result = range_code(data, coded, 16, 32);

    printf("Result: %ld\n", result);
#endif

    return 0;
}