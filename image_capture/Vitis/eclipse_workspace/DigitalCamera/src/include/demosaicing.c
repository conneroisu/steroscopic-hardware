#include "demosaicing.h"

int location_valid(t_location location)
{
    return (location.row >= 0 && location.row < HEIGHT) && (location.col >= 0 && location.col < WIDTH);
}

void get_neighbors(t_location location, t_neighbors* neighbors)
{
    // Reset neighbors
    neighbors->count = 0;
    for(int i = 0; i < MAX_NEIGHBOR_COUNT; ++i)
    {
        neighbors->neighbors[i] = (t_location) {.row = 0, .col = 0};
    }

    // 9 possible neighbors.
    for(int i = -1; i <= 1; ++i)
    {
        for(int j = -1; j <= 1; ++j)
        {
            // Skip the provided location.
            if(!i && !j)
            {
                continue;
            }
            t_location loc = {.row = location.row + i, .col = location.col + j};

            if(location_valid(loc))
            {
                neighbors->neighbors[neighbors->count] = loc;
                neighbors->count++;
            }
        }
    }
}

t_colors get_filter_color(t_location location)
{
    // 4 Cases:
    // Even row, even column = BLUE
    // Even row, odd column = GREEN
    // Odd row, even column = GREEN
    // Odd row, odd column = RED
    return location.row % 2 ? location.col % 2 ? RED : GREEN : location.col % 2 ? GREEN : BLUE;
}

void run_demosaicing(uint16_t* intensities, uint16_t* yuv_out)
{
    t_location location;
    t_neighbors neighbors;
    int red_count, blue_count, green_count, red_total, blue_total, green_total, index = 0;
    uint8_t intensity = 0;

    // Iterate through all the intensities.
    for(int i = 0; i < HEIGHT; ++i)
    {
    	// Take RGB data only from left pixel currently.
    	// 32-bits = 2 pixels of data.
        for(int j = 0; j < WIDTH; j+=2)
        {
            index = (i * WIDTH) + j;
            red_count = 0;
            green_count = 0;
            blue_count = 0;
            red_total = 0;
            green_total = 0;
            blue_total = 0;

            location = (t_location){.row = i, .col = j};

            // Get neighbors:
            get_neighbors(location, &neighbors);

            for(int n = 0; n < neighbors.count; ++n)
            {
                t_location n_loc = neighbors.neighbors[n];
                int n_index = (WIDTH * n_loc.row) + n_loc.col;
                intensity = intensities[n_index] & 0xFF;

                switch(get_filter_color(n_loc))
                {
                    case RED:
                    {
                        red_total += intensity;
                        red_count++;
                        break;
                    }

                    case GREEN:
                    {
                        green_total += intensity;
                        green_count++;
                        break;
                    }

                    case BLUE:
                    {
                        blue_total += intensity;
                        blue_count++;
                        break;
                    }
                }
            }

            intensity = intensities[index] & 0xFF;
            // Write current color
            switch(get_filter_color(location))
            {
                case RED:
                {
                    red_total += intensity;
                    red_count++;
                    break;
                }

                case GREEN:
                {
                    green_total += intensity;
                    green_count++;
                    break;
                }

                case BLUE:
                {
                    blue_total += intensity;
                    blue_count++;
                    break;
                }
            }

            t_color_24_bit rgb = (t_color_24_bit) {.red = !red_count ? 0 : red_total / red_count, .green = !green_count ? 0 : green_total / green_count, .blue = !blue_count ? 0 : blue_total / blue_count};

            uint32_t yuv = rgb_to_yuv(rgb);

            // Record colors
            yuv_out[index] = (uint16_t) (yuv & 0xFFFF);
            yuv_out[index + 1] = (uint16_t) ((yuv & 0xFFFF0000) >> 16);
        }
    }
}

// 32-bit output
// Ordering: 0  1  2  3
//           U0 Y0 V0 Y1
uint32_t rgb_to_yuv(t_color_24_bit rgb)
{
	uint32_t result = 0;

	// Convert to YUV 4:2:2
	uint8_t u, v = 0;
	uint16_t y = 0;

	y = (((uint16_t) rgb.red) * 0.183f) + (((uint16_t) rgb.green) * 0.614f) + (((uint16_t) rgb.blue) * 0.062f) + 16;
	v = (rgb.red * -0.101f) + (rgb.green * -0.338f) + (rgb.blue * 0.439f) + 128;
	u = (rgb.red * 0.439f) + (rgb.green * -0.399f) + (rgb.blue * -0.04f) + 128;

	// Set from YUV values.
	result = (uint32_t)((u << 8) | (y & 0xFF)) | ((uint32_t)((v << 8) | ((y & 0xFF00) >> 8)) << 16);

	return result;
}

