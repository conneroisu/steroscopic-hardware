#include "range_code.h"

int main()
{
    uint8_t data[256];
    uint8_t coded[256];

    for(int i = 0; i < 256; ++i)
    {
        data[i] = i;
    }

    range_code(data, coded, 256);

    return 0;
}