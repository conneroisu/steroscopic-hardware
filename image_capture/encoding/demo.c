#include "range_code.h"
#include <stdio.h>
#include <errno.h>

#define TEST_MODE 1

#define SIZE 2073600U
#define CHUNK_SIZE 64800U

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
    
    uint8_t* current_data = data;
    uint8_t* current_coded = coded;

    for(int i = 0; i < (SIZE / CHUNK_SIZE); ++i)
    {
        size_t new_result = range_code(current_data, current_coded, CHUNK_SIZE, 32768);
        current_coded += (new_result / 8);
        current_data += CHUNK_SIZE;

        if(!new_result)
        {
            printf("Error: Adjustment factor of %d is too small!\n", 32768);
        }

        result += new_result;
    }



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

    size_t result = range_code(data, coded, 1000, 32);

    free(data);
    free(coded);
    printf("Result: %ld\n", result);

#else
    uint8_t data[10] = { 130, 55, 39, 55, 130, 72, 72, 9, 72, 8 };
    uint8_t coded[10];

    memset(coded, 0, 10);

    size_t result = range_code(data, coded, 10, 32);

    printf("Result: %ld\n", result);
#endif

    return 0;
}