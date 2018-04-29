package main

// $ go build -buildmode=c-archive -o glitch.a ffi.go

/*
#include <stdlib.h>

typedef struct {
    char *data;
    size_t size;
} Image_t;
*/
import "C"

import (
    "image"
    "image/png"
    _ "image/jpeg"
    "bytes"
    "unsafe"

    "github.com/sugoiuguu/go-glitch"
)

// this does nothing...
func main() {}

//export jumble_pixels
func jumble_pixels(expression, data *C.char, size C.int) *C.Image_t {
    // compile expression
    expr, err := glitch.CompileExpression(C.GoString(expression))
    if err != nil {
        return nil
    }

    // write the raw image data
    // to a bytes.Buffer
    buf := new(bytes.Buffer)
    if _,err := buf.Write(C.GoBytes(unsafe.Pointer(data), size)); err != nil {
        return nil
    }

    // decode image from bytes.Buffer
    img, _, err := image.Decode(buf)
    if err != nil {
        return nil
    }

    // glitch that shit
    glitchedimg, err := expr.JumblePixels(img)
    if err != nil {
        return nil
    }

    // encode to a png
    img = nil
    buf.Reset()
    if err := png.Encode(buf, glitchedimg); err != nil {
        return nil
    }

    // return raw image data
    b := buf.Bytes()
    l := C.size_t(len(b))
    m := C.malloc(l)
    if m == nil {
        return nil
    }

    s := (*C.Image_t)(m)
    s.size = l
    s.data = (*C.char)(C.CBytes(b))

    return s
}
