package glitch

import (
    "image"
    "image/color"
)

func (expr Expression) JumblePixels(data image.Image) (image.Image, error) {
    return expr.JumblePixelsMonitor(data, nil)
}

func (expr Expression) JumblePixelsMonitor(data image.Image, mon func()) (image.Image, error) {
    img := image.NewNRGBA(data.Bounds())
    xm, xM := data.Bounds().Min.X, data.Bounds().Max.X
    ym, yM := data.Bounds().Min.Y, data.Bounds().Max.Y

    var err error
    var nr, ng, nb, sr, sg, sb uint8

    box := make([]color.RGBA, 9)

    for x := xm; x < xM; x++ {
        for y := ym; y < yM; y++ {
            k := 0

            for i := x - 1; i <= x + 1; i++ {
                for j := y - 1; j <= y + 1; j++ {
                    r, g, b := convUint8(data.At(i, j).RGBA())
                    box[k] = color.RGBA{R: r, G: g, B: b, A: 255}
                    k++
                }
            }

            nr, ng, nb, err = expr.evalRPN(box[4].R, sr, box[4].G, sg, box[4].B, sb, box)
            if err != nil {
                return nil, err
            }

            sr = nr
            sg = ng
            sb = nb

            img.Set(x, y, color.NRGBA{
                nr, ng, nb,
                255,
            })

            if mon != nil {
                mon()
            }
        }
    }

    return image.Image(img), nil
}
