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

    uint8_t* image = calloc(IMAGE_HEIGHT * IMAGE_WIDTH, sizeof(uint8_t));

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

    // Setup UART config.
    tty.c_cflag &= ~PARENB;
    tty.c_cflag &= ~CSTOPB;
    tty.c_cflag &= ~CSIZE;
    tty.c_cflag |= CS8;
    tty.c_cflag &= ~CRTSCTS;
    tty.c_cflag |= CREAD | CLOCAL;

    tty.c_lflag &= ~ICANON;
    tty.c_lflag &= ~ECHO;
    tty.c_lflag &= ~ISIG;

    tty.c_iflag |= ~(IXON | IXOFF | IXANY);
    tty.c_iflag &= ~(IGNBRK | BRKINT | PARMRK | ISTRIP | INLCR | IGNCR | ICRNL);

    tty.c_oflag &= ~OPOST;
    tty.c_oflag &= ~ONLCR;

    // Always block
    tty.c_cc[VTIME] = 0;
    tty.c_cc[VMIN] = 1;

    cfsetspeed(&tty, B115200);

    // Apply settings.
    status = tcsetattr(port, TCSANOW, &tty);
    if(status < 0)
    {
        perror("Could not apply the tty config!");
        return status;
    }

    printf("Getting picture!\n");
    char write_buf = 1;

    // Request an image
    write(port, &write_buf, 1);

    // Make sure we read back the opcode.
    char read_byte;
    read(port, &read_byte, 1);

    if(read_byte == 1)
    {
        printf("Got reply from MP-2 hardware!\n");
        // Each sensor reading is 16 bits.
       for(int i = 0; i < (IMAGE_WIDTH * IMAGE_HEIGHT); ++i)
       {
            read(port, &read_byte, 1);
            image[i] = read_byte;
       }

        FILE* output = fopen("image.bin", "w");

        if(output == NULL)
        {
            perror("Could not open output file to write image to!");
            return errno;
        }

        fwrite(image, 1, (IMAGE_WIDTH * IMAGE_HEIGHT), output);
        
        status = fclose(output);

        if(status < 0)
        {
            perror("Could not close output file!");
            return status;
        }
    }
    else
    {
        printf("Got bad response! Actual value: %d\n", read_byte);
    }

    // Close the port
    status = close(port);

    if(status < 0)
    {
        char str_buf[256];
        sprintf(str_buf, "ERROR: Failed to close port %s", PORT_NAME);
        perror(str_buf);
        return status;
    }

    printf("Finished getting picture!!\n");

    return 0;
}
