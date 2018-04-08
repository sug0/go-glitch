package main

import (
    "os"
    "fmt"
    "flag"
    "time"
    "strings"
    "math/rand"
    "image"
    "image/png"
    _ "image/jpeg"
    _ "image/gif"

    "github.com/cheggaaa/pb"
    "github.com/sugoiuguu/go-glitch"
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

    img, _, err := image.Decode(r)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
        os.Exit(1)
    }

    rand.Seed(time.Now().UnixNano())
    pixsize := img.Bounds().Dx() * img.Bounds().Dy()

    fmt.Fprintf(os.Stderr, "Glitching %d pixels %d time(s)...\n", pixsize, len(exprs))

    for _,expr := range exprs {
        bar := pb.New(pixsize).SetMaxWidth(80)
        bar.Output = os.Stderr
        bar.Start()

        img, err = expr.JumblePixelsMonitor(img, func(_, _ int) { bar.Increment() })
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
            os.Exit(1)
        }

        bar.Finish()
    }

    png.Encode(f, img)
}

func compileExprs(exprs string) ([]glitch.Expression, error) {
    var e []glitch.Expression

    for _,exprS := range strings.Split(exprs, ",") {
        expr, err := glitch.CompileExpression(exprS)
        if err != nil {
            return nil, err
        }
        e = append(e, expr)
    }

    return e, nil
}
