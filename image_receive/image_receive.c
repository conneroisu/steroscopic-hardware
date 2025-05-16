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

void clear_uart_fifo(struct termios* tty, int port)
{
    // Configure TTY to not block on reads
    cc_t old_vmin = tty->c_cc[VMIN];

    tty->c_cc[VMIN] = 0;

    int status = tcsetattr(port, TCSANOW, tty);
    if(status < 0)
    {
        perror("Could not apply the tty config!");
        return;
    }

    char* buf = 0;

    while(read(port, &buf, 1))
    {

    }

    // Restore settings
    tty->c_cc[VMIN] = old_vmin;

    status = tcsetattr(port, TCSANOW, tty);
    if(status < 0)
    {
        perror("Could not apply the tty config!");
        return;
    }
}

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

int main(int argc, char** argv)
{
    int test_mode = 0;

    if(argc <= 1)
    {
        printf("Serial port was not provided!\n");
        return 0;
    }
    else if(argc == 3)
    {
        if(!strcmp(argv[2], "TEST"))
        {
            printf("Test mode active!\n");
            test_mode = 1;
        }
    }

    uint8_t* image = calloc(IMAGE_HEIGHT * IMAGE_WIDTH, sizeof(uint8_t));

    struct termios tty;

    char* PORT_NAME = argv[1];

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

    usleep(10000);

    // Clear FIFO:
    clear_uart_fifo(&tty, port);

    usleep(10000);

    if(test_mode)
    {
        printf("Starting test!\n");
        uint8_t write_buf[2] = {0xFF, 0xD8};
        write(port, write_buf, 2);
        
        char read_byte;
        read(port, &read_byte, 1);

        if(read_byte != 1)
        {
            printf("Test failed! Sender did not send a 0x1 after image send request! Acutal: %x\n", read_byte);
            goto exit;
        }
        else
        {
            printf("Test Case: Sender sends 0x1 after receiving start sequence: Passed!\n");
        }

        uint8_t bytes[256];

        // Get 256 bytes
        for(int i = 0; i < 256; ++i)
        {
            read(port, &read_byte, 1);
            bytes[i] = read_byte;
        }

        // Verify that all the bytes are not all zero.
        int all_zero = 1;
        for(int i = 0; i < 256; ++i)
        {
            if(bytes[i] != 0)
            {
                all_zero = 0;
                break;
            }
        }

        if(all_zero)
        {
            printf("Test failed! First 256 from image were all zero! This should never happen! Data is bad!\n");
            goto exit;
        }
        else
        {
            printf("Test Case: First 256 bytes are not all zero: Passed!\n");
        }

        // Send stop signal
        write_buf[1] = 0xD9;
        write(port, write_buf, 2);

        // Wait a second:
        sleep(1);

        // Clear the FIFO:
        clear_uart_fifo(&tty, port);

        usleep(10000);

        // Verify that the data has stopped.
        tty.c_cc[VMIN] = 0;

        int status = tcsetattr(port, TCSANOW, &tty);
        if(status < 0)
        {
            perror("Could not apply the tty config!");
            return errno;
        }

        if(read(port, &read_byte, 1))
        {
            printf("Test failed! A byte was still in the FIFO after a clear and after sending stop sequence, so data is still being sent!\n");
            goto exit;
        }
        else
        {
            printf("Test Case: Image data is no longer sent after sender receives stop sequence: Passed!\n");
        }

        printf("Test passed!\n");
        goto exit;
    }
    else
    {
        printf("Getting picture!\n");
        uint8_t write_buf[2] = {0xFF, 0xD8};

        // Request an image
        write(port, write_buf, 2);

        // Make sure we read back the opcode.
        char read_byte;
        read(port, &read_byte, 1);

        if(read_byte == 1)
        {
            printf("Got reply from MP-2 hardware!\n");
            // Each sensor reading is 16 bits.
        for(int i = 0; i < (IMAGE_HEIGHT * IMAGE_WIDTH); ++i)
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

        write_buf[1] = 0xD9;

        write(port, write_buf, 2);

        usleep(10000);

        clear_uart_fifo(&tty, port);

        usleep(10000);

        printf("Finished getting picture!!\n");
    }

    exit:
    // Close the port
    status = close(port);

    if(status < 0)
    {
        char str_buf[256];
        sprintf(str_buf, "ERROR: Failed to close port %s", PORT_NAME);
        perror(str_buf);
        return status;
    }

    return 0;
}
