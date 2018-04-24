#include <stdio.h>
#include <stdlib.h>

#include "glitch.h"

/*
 * $ cc -Wall -O2 -o ffi ffi.c glitch.a -lpthread */

int main(int argc, char *argv[])
{
    FILE *fp, *fpw;
    char *buf;
    Image_t *img;
    int size;

    if (argc < 2)
        return 1;

    fp = fopen(argv[1], "r");
    
    if (!fp)
        return 1;

    fseek(fp, 0, SEEK_END);
    size = ftell(fp);
    rewind(fp);

    buf = malloc(size);

    if (!buf) {
        fclose(fp);
        return 1;
    }

    puts("glitching...");
    fread(buf, size, 1, fp);
    img = jumble_pixels("Y", buf, size);

    if (!img) {
        free(buf);
        fclose(fp);
        return 1;
    }

    fpw = fopen("out.png", "w");

    if (!fpw) {
        free(buf);
        fclose(fp);
        return 1;
    }

    fwrite(img->data, img->size, 1, fpw);

    free(img->data);
    free(img);
    free(buf);
    fclose(fp);
    fclose(fpw);
    puts("done...");

    return 0;
}
