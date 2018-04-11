package glitch

//import (
//    //"math"
//)

func convUint8(r, g, b, a uint32) (uint8, uint8, uint8, uint8) {
    return uint8(r / 0x101), uint8(g / 0x101), uint8(b / 0x101), uint8(a / 0x101)
}

func threeRule(x, max int) uint8 {
    return uint8(((255 * x) / max) & 255)
}

//func normSin(x int) uint8 {
//    return uint8(255 * math.Sin(float64(x)))
//}
//
//func normCos(x int) uint8 {
//    return uint8(255 * math.Cos(float64(x)))
//}
