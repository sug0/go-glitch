package glitch

import (
    "bytes"
    "image"
    "image/gif"
)

// pretty shit hack to get a Paletted image out of an NRGBA image
func imgNRGBAToPaletted(data *image.NRGBA) (*image.Paletted, error) {
    buf := new(bytes.Buffer)

    if err := gif.Encode(buf, data, nil); err != nil {
        return nil, err
    }

    decoded, err := gif.Decode(buf)
    if err != nil {
        return nil, err
    }

    return decoded.(*image.Paletted), nil
}

func convUint8(r, g, b, a uint32) (uint8, uint8, uint8, uint8) {
    return uint8(r / 0x101), uint8(g / 0x101), uint8(b / 0x101), uint8(a / 0x101)
}

func threeRule(x, max int) uint8 {
    return uint8(((255 * x) / max) & 255)
}

func max(vals ...uint8) (m uint8) {
    for _,v := range vals {
        if v > m {
            m = v
        }
    }
    return
}

func min(vals ...uint8) (m uint8) {
    m = 255
    for _,v := range vals {
        if v < m {
            m = v
        }
    }
    return
}
