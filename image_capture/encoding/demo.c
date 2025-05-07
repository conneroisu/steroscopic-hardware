#include "range_code.h"

int main()
{
    uint8_t data[5] = {128, 45, 96, 4, 239};
    uint8_t coded[5];

    range_code(data, coded, 5);

    return 0;
}