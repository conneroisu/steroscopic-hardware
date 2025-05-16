#include <stdio.h>
#include <stdlib.h>
#include <limits.h>

// #define WIDTH  640
// #define HEIGHT 480
// #define MAX_DISP 64
// #define WINDOW_SIZE 15

#define WIDTH   128
#define HEIGHT  128
#define MAX_DISP 64
#define WINDOW_SIZE 15

int abs_diff (int a, int b) {
    if (a >= b) {
        return a - b;
    } else {
        return b - a;
    }
}

// int min_usr (int a, int b) {
//     if (a >= b) {
//         return a; 
//     } else {
//         return b;
//     }
// }

// int max_usr (int a, int b) {
//     if (a <= b) {
//         return a; 
//     } else {
//         return b;
//     }
// }

int main() {

    // ChatGPT code for opening and importing png //
    int image_size = WIDTH * HEIGHT;
    unsigned char* img_L = malloc(image_size);
    unsigned char* img_R = malloc(image_size);
    int* disparity = calloc(image_size, sizeof(int));
    if (!img_L || !img_R || !disparity) {
        perror("Memory allocation failed");
        return 1;
    }

    // FILE* fL = fopen("img_L.raw", "rb");
    // FILE* fR = fopen("img_R.raw", "rb");
    FILE* fL = fopen("mems/img_L_patch_1.raw","rb");
    FILE* fR = fopen("mems/img_R_patch_1.raw","rb");
    if (!fL || !fR) {
        perror("Failed to open input files");
        return 1;
    }

    fread(img_L, 1, image_size, fL);
    fread(img_R, 1, image_size, fR);
    fclose(fL);
    fclose(fR);
    ////////////////////////////////////////////////

    int margin = WINDOW_SIZE / 2;

    for (int y = margin; y < HEIGHT - margin; y++) {
        for (int x = margin; x < WIDTH - margin; x++) {
            int min_sad = INT_MAX;
            int best_disp = 0;

            for (int d = 0; d <= MAX_DISP; d++) {
                int sad = 0;

                for (int dy = -margin; dy <= margin; dy++) {
                    for (int dx = -margin; dx <= margin; dx++) {
                        int xl = x + dx; // img_L x pos
                        int yl = y + dy; // img_L y pos (same for img_R) 
                        int xr = xl - d; // img_R x pos

                        // Skip boundary pixels
                        if (xr < 0) {
                            continue;
                        }

                        // Indexing 1D array: Index = (y * width) + x
                        int idx_L = yl * WIDTH + xl;
                        int idx_R = yl * WIDTH + xr;

                        // Accumulate SAD values in window
                        sad += abs_diff(img_L[idx_L], img_R[idx_R]);
                    }
                }

                // Save smallest SAD value from windows
                if (sad < min_sad) {
                    min_sad = sad;
                    best_disp = d;
                }
            }

            // Save disparity for pixel 
            disparity[y * WIDTH + x] = best_disp;
        }
    }

    // ChatGPT code for writing output to pgm //
    FILE* fout = fopen("disparity.pgm", "wb");
    if (!fout) {
        perror("Failed to write output");
        return 1;
    }
    
    // Write PGM header (P5 = binary grayscale)
    fprintf(fout, "P5\n%d %d\n255\n", WIDTH, HEIGHT);
    
    // Normalize disparity to 0–255
    int max_val = 0;
    for (int i = 0; i < WIDTH * HEIGHT; i++) {
        if (disparity[i] > max_val) max_val = disparity[i];
    }
    
    for (int i = 0; i < WIDTH * HEIGHT; i++) {
        unsigned char val = (unsigned char)(255.0 * disparity[i] / max_val);
        fwrite(&val, 1, 1, fout);
    }
    
    fclose(fout);
    ////////////////////////////////////////////

    // Write out disparity in hex mem format
    FILE* fm = fopen("exp_disp.mem","w");
    if (!fm) {
    perror("fopen exp_disp.mem");
    return 1;
    }
    for (int i = 0; i < WIDTH * HEIGHT; i++) {
    unsigned val = disparity[i];
    fprintf(fm, "%02X\n", val & 0xFF);
    }
    fclose(fm);

    free(img_L);
    free(img_R);
    free(disparity);
    return 0;
}