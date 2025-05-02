#ifndef __DEMOSAICING_H_
#define __DEMOSAICING_H_
#endif

#include "xil_types.h"

#define WIDTH 1920
#define HEIGHT 1080
#define PIXELS 2073600
#define MAX_NEIGHBOR_COUNT 9

typedef enum colors
{
    RED, GREEN, BLUE
} t_colors;

typedef struct color_24_bit
{
    uint8_t red;
    uint8_t green;
    uint8_t blue;
} t_color_24_bit;

typedef struct location
{
    int row;
    int col;
} t_location;

typedef struct neighbors
{
    t_location neighbors[MAX_NEIGHBOR_COUNT];
    uint32_t count;
} t_neighbors;

t_colors get_filter_color(t_location location);
void get_neighbors(t_location location, t_neighbors* neighbors);
void run_demosaicing(uint16_t* intensities, uint16_t* yuv_out);
int location_valid(t_location location);
void write_ppm(const char* file_name, t_color_24_bit* image_colors);
uint32_t rgb_to_yuv(t_color_24_bit rgb);
