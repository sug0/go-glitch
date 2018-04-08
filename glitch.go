package glitch

import (
    "image"
    "image/color"
)

func (expr Expression) JumblePixels(data image.Image) (image.Image, error) {
    return expr.JumblePixelsMonitor(data, nil)
}

func (expr Expression) JumblePixelsMonitor(data image.Image, mon func(int, int)) (image.Image, error) {
    img := image.NewNRGBA(data.Bounds())
    xm, xM := data.Bounds().Min.X, data.Bounds().Max.X
    ym, yM := data.Bounds().Min.Y, data.Bounds().Max.Y

    var err error
    var nr, ng, nb, sr, sg, sb uint8
    var i, pixsize int

    if mon != nil {
        i = 0
        pixsize = data.Bounds().Dx() * data.Bounds().Dy()
        mon(i, pixsize)
    }

    box := make([]sum, 9)

    for x := xm; x < xM; x++ {
        for y := ym; y < yM; y++ {
            k := 0

            for i := x - 1; i <= x + 1; i++ {
                for j := y - 1; j <= y + 1; j++ {
                    r, g, b := convUint8(data.At(i, j).RGBA())
                    box[k] = sum{r: r, g: g, b: b}
                    k++
                }
            }

            nr, ng, nb, err = expr.evalRPN(x, y, xM, yM,
                                           box[4].r, sr, box[4].g, sg, box[4].b, sb, box)
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
                mon(i, pixsize)
                i++
            }
        }
    }

    return image.Image(img), nil
}
