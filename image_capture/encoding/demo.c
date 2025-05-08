#include "range_code.h"

int main()
{
    uint8_t data[10] = {130, 55, 39, 55, 130, 72, 72, 9, 72, 8};
    uint8_t coded[10];

    memset(coded, 0, 10);

    for(int i = 0; i < 10; ++i)
    {
        data[i] = i;
    }

    range_code(data, coded, 10);

    return 0;
}