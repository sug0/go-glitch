package main

import (
    "io"
    "os"
    "fmt"
    "flag"
    "time"
    "strings"
    "math/rand"
    "bytes"
    "image"
    "image/png"
    "image/gif"
    _ "image/jpeg"

    "github.com/cheggaaa/pb"
    "github.com/sugoiuguu/go-glitch"
)

const (
)

func main() {
    var f, r *os.File

    out := flag.String("o", "out.png", "output file")
    in := flag.String("i", "<none>", "input file")
    exprS := flag.String("e", "c ^ 255", "comma separated list of expressions")
    flag.Parse()

    exprs, err := compileExprs(*exprS)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s: couldn't compile expression %q\n", os.Args[0], *exprS)
        os.Exit(1)
    }

    if *in == "<none>" {
        fmt.Fprintf(os.Stderr, "usage: %s -i <png, jpg or gif image>\n", os.Args[0])
        os.Exit(1)
    } else if *in == "-" {
        r = os.Stdin
    } else {
        r2, err := os.Open(*in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
            os.Exit(1)
        }
        r = r2
        defer r.Close()
    }

    img, err := decodeImage(r)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
        os.Exit(1)
    }

    switch img.(type) {
    case *gif.GIF:
        if *out == "out.png" {
            *out = "out.gif"
        }
    }

    if *out == "-" {
        f = os.Stdout
    } else {
        f2, err := os.Create(*out)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
            os.Exit(1)
        }
        f = f2
        defer f.Close()
    }

    rand.Seed(time.Now().UnixNano())
    fmt.Fprintf(os.Stderr, "Glitching a shit ton of pixels %d time(s)...\n\n", len(exprs))

    for i, expr := range exprs {
        var bar *pb.ProgressBar
        var set bool
        var err error

        mon := func(_,pixsize int) {
            if !set {
                set = true
                bar = pb.New(pixsize).SetMaxWidth(80)
                bar.Output = os.Stderr
                bar.Start()
            }
            bar.Increment()
        }

        fmt.Fprintf(os.Stderr, "[%d]: Using expression: %s\n", i + 1, expr)

        switch img.(type) {
        case *gif.GIF:
            img, err = expr.JumbleGIFPixelsMonitor(img.(*gif.GIF), mon)
        default:
            img, err = expr.JumblePixelsMonitor(img.(image.Image), mon)
        }

        if set { bar.Finish() }
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
            os.Exit(1)
        }
    }

    switch img.(type) {
    case *gif.GIF:
        gif.EncodeAll(f, img.(*gif.GIF))
    default:
        png.Encode(f, img.(image.Image))
    }
}

func compileExprs(exprsS string) ([]*glitch.Expression, error) {
    exprs := strings.Split(exprsS, ",")
    e := make([]*glitch.Expression, 0, len(exprs))

    for _,exprS := range exprs {
        expr, err := glitch.CompileExpression(exprS)
        if err != nil {
            return nil, err
        }
        e = append(e, expr)
    }

    return e, nil
}

func decodeImage(r io.Reader) (interface{}, error) {
    var willDecodeNormal bool
    buf := new(bytes.Buffer)

decodeOtherFormats:
    // on first try is false
    // so it skips this step
    if willDecodeNormal {
        if img, _, err := image.Decode(buf); err != nil {
            return nil, err
        } else {
            return img, nil
        }
    }

    // read r contents into buffer
    if _,err := buf.ReadFrom(r); err != nil {
        return nil, err
    }

    // not a valid gif, try decoding
    // another format
    gifmagic := buf.Bytes()
    if len(gifmagic) < 6 {
        willDecodeNormal = true
        goto decodeOtherFormats
    }
    gifmagic = gifmagic[:6]

    // we found our gif -- decode that shit
    switch {
    case bytes.Equal(gifmagic, []byte("GIF87a")):
        fallthrough
    case bytes.Equal(gifmagic, []byte("GIF89a")):
        img, err := gif.DecodeAll(buf)
        if err != nil {
            return nil, err
        } else {
            return img, nil
        }
    }

    // back to the basics
    willDecodeNormal = true
    goto decodeOtherFormats
}
