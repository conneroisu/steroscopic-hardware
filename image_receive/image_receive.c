#include <termios.h>
#include <stdio.h>
#include <fcntl.h>
#include <string.h>
#include <errno.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdint.h>

#define IMAGE_WIDTH 1920
#define IMAGE_HEIGHT 1080

// Referenced: https://blog.mbedded.ninja/programming/operating-systems/linux/linux-serial-ports-using-c-cpp/

int open_port(const char* port_name)
{
    int port = open(port_name, O_RDWR);

    if(port < 0)
    {
        char str_buf[256];

        sprintf(str_buf, "ERROR: Could not read from %s", port_name);
        perror(str_buf);
    }

    return port;
}

int main()
{

    uint16_t* image = malloc(IMAGE_HEIGHT * IMAGE_WIDTH * sizeof(uint16_t));

    struct termios tty;

    const char PORT_NAME[] = "/dev/ttyS22";

    // Open the port
    int port = open_port(PORT_NAME);
    if(port < 0)
    {
        return port;
    }

    // Get serial configuration.
    int status = tcgetattr(port, &tty);
    if(status < 0)
    {
        char str_buf[256];
        sprintf(str_buf, "ERROR: Failed to get tty attributes from %s", PORT_NAME);
        perror(str_buf);
        return status;
    }

    // Disable echo
    tty.c_lflag &= ~ECHO;

    char write_buf = 1;

    // Request an image
    write(port, &write_buf, 1);

    // Make sure we read back the opcode.
    char read_byte;
    read(port, &read_byte, 1);

    if(read_byte == 1)
    {
        // Each sensor reading is 16 bits.
       for(int i = 0; i < (IMAGE_WIDTH * IMAGE_HEIGHT); ++i)
       {
            read(port, &read_byte, 1);
            image[i] = (uint16_t) read_byte;
            read(port, &read_byte, 1);
            image[i] |= (((uint16_t) read_byte) << 8);
       }
    }

    return 0;
}
