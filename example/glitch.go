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
    "bufio"
    "image"
    "image/png"
    "image/gif"
    _ "image/jpeg"

    "github.com/cheggaaa/pb"
    "github.com/sug0/go-glitch"
    "github.com/sug0/go-exit"
)

func main() {
    defer exit.Handler()
    var f, r *os.File

    out := flag.String("o", "out.png", "output file")
    in := flag.String("i", "<none>", "input file")
    exprS := flag.String("e", "c ^ 255", "comma separated list of expressions")
    flag.Parse()

    exprs, err := compileExprs(*exprS)
    if err != nil {
        exit.WithMsg(os.Stderr, 1, "%s: couldn't compile expression %q", os.Args[0], *exprS)
    }

    if *in == "<none>" {
        exit.WithMsg(os.Stderr, 1, "usage: %s -i <png, jpg or gif image>", os.Args[0])
    } else if *in == "-" {
        r = os.Stdin
    } else {
        r2, err := os.Open(*in)
        if err != nil {
            exit.WithMsg(os.Stderr, 1, "%s: %s", os.Args[0], err)
        }
        r = r2
        defer r.Close()
    }

    img, err := decodeImage(r)
    if err != nil {
        exit.WithMsg(os.Stderr, 1, "%s: %s", os.Args[0], err)
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
            exit.WithMsg(os.Stderr, 1, "%s: %s", os.Args[0], err)
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
            exit.WithMsg(os.Stderr, 1, "%s: %s", os.Args[0], err)
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
    br := bufio.NewReader(r)

decodeOtherFormats:
    // on first try is false
    // so it skips this step
    if willDecodeNormal {
        if img, _, err := image.Decode(br); err != nil {
            return nil, err
        } else {
            return img, nil
        }
    }

    // not a valid gif, try decoding
    // another format
    gifmagic, err := br.Peek(6)
    if err != nil {
        willDecodeNormal = true
        goto decodeOtherFormats
    }

    // we found our gif -- decode that shit
    switch {
    case bytes.Equal(gifmagic, []byte("GIF87a")):
        fallthrough
    case bytes.Equal(gifmagic, []byte("GIF89a")):
        img, err := gif.DecodeAll(br)
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
