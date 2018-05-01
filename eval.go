package glitch

import (
    "fmt"
    "math/rand"
    "image/color"
    "strconv"
)

type sum struct {
    r, g, b uint8
}

// save some values that
// might be expensive to recalculate
type sumsave struct {
    v_Y,
    v_e,
    v_b,
    v_H,
    v_L *sum
}

func (expr *Expression) evalRPN(x, y, w, h int,
                                sr, sg, sb uint8,
                                box []color.NRGBA,
) (rr uint8, gr uint8, br uint8, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("error evaluating expression: %s", expr.infix)
            return
        }
    }()

    if box[4].A == 0 {
        return
    }

    stk := make([]sum, 0, len(expr.toks))
    saved := sumsave{}

    for _,tok := range expr.toks {
        if tok == "" {
            continue
        }
        if oper := operMap[[]rune(tok)[0]]; oper != nil {
            a, b := stk[len(stk)-2], stk[len(stk)-1]
            stk = append(stk[:len(stk)-2], sum{oper.f(a.r, b.r),
                                               oper.f(a.g, b.g),
                                               oper.f(a.b, b.b)})
        } else if tok == "c" {
            stk = append(stk, sum{box[4].R, box[4].G, box[4].B})
        } else if tok == "R" {
            stk = append(stk, sum{255, 0, 0})
        } else if tok == "G" {
            stk = append(stk, sum{0, 255, 0})
        } else if tok == "B" {
            stk = append(stk, sum{0, 0, 255})
        } else if tok == "Y" {
            if saved.v_Y == nil {
                y := uint8(float64(box[4].R)*0.299) +
                     uint8(float64(box[4].G)*0.587) +
                     uint8(float64(box[4].B)*0.0722)
                saved.v_Y = &sum{y, y, y}
            }
            stk = append(stk, *saved.v_Y)
        } else if tok == "s" {
            stk = append(stk, sum{sr, sg, sb})
        } else if tok == "x" {
            xu := threeRule(x, w)
            stk = append(stk, sum{xu, xu, xu})
        } else if tok == "y" {
            yu := threeRule(y, h)
            stk = append(stk, sum{yu, yu, yu})
        } else if tok == "r" {
            i, j, k := rand.Int() % 8, rand.Int() % 8, rand.Int() % 8
            stk = append(stk, sum{box[i].R, box[j].G, box[k].B})
        } else if tok == "e" {
            if saved.v_e == nil {
                dr := box[8].R - box[0].R + box[5].R - box[3].R +
                      box[7].R - box[1].R + box[6].R - box[2].R
                dg := box[8].G - box[0].G + box[5].G - box[3].G +
                      box[7].G - box[1].G + box[6].G - box[2].G
                db := box[8].B - box[0].B + box[5].B - box[3].B +
                      box[7].B - box[1].B + box[6].B - box[2].B
                saved.v_e = &sum{dr, dg, db}
            }
            stk = append(stk, *saved.v_e)
        } else if tok == "b" {
            if saved.v_b == nil {
                sr := int(box[0].R) + int(box[1].R) + int(box[2].R) + int(box[3].R) +
                      int(box[4].R) + int(box[5].R) + int(box[6].R) + int(box[7].R) + int(box[8].R)
                sg := int(box[0].G) + int(box[1].G) + int(box[2].G) + int(box[3].G) +
                      int(box[4].G) + int(box[5].G) + int(box[6].G) + int(box[7].G) + int(box[8].G)
                sb := int(box[0].B) + int(box[1].B) + int(box[2].B) + int(box[3].B) +
                      int(box[4].B) + int(box[5].B) + int(box[6].B) + int(box[7].B) + int(box[8].B)
                saved.v_b = &sum{uint8(sr/9), uint8(sg/9), uint8(sb/9)}
            }
            stk = append(stk, *saved.v_b)
        } else if tok == "H" {
            if saved.v_H == nil {
                rM := max(box[0].R, box[1].R, box[2].R, box[3].R,
                          box[4].R, box[5].R, box[6].R, box[7].R, box[8].R)
                gM := max(box[0].G, box[1].G, box[2].G, box[3].G,
                          box[4].G, box[5].G, box[6].G, box[7].G, box[8].G)
                bM := max(box[0].B, box[1].B, box[2].B, box[3].B,
                          box[4].B, box[5].B, box[6].B, box[7].B, box[8].B)
                saved.v_H = &sum{rM, gM, bM}
            }
            stk = append(stk, *saved.v_H)
        } else if tok == "L" {
            if saved.v_L == nil {
                rM := min(box[0].R, box[1].R, box[2].R, box[3].R,
                          box[4].R, box[5].R, box[6].R, box[7].R, box[8].R)
                gM := min(box[0].G, box[1].G, box[2].G, box[3].G,
                          box[4].G, box[5].G, box[6].G, box[7].G, box[8].G)
                bM := min(box[0].B, box[1].B, box[2].B, box[3].B,
                          box[4].B, box[5].B, box[6].B, box[7].B, box[8].B)
                saved.v_L = &sum{rM, gM, bM}
            }
            stk = append(stk, *saved.v_L)
        } else if tok == "N" {
            rn, gn, bn := uint8(rand.Int() % 256), uint8(rand.Int() % 256), uint8(rand.Int() % 256)
            stk = append(stk, sum{rn, gn, bn})
        } else {
            if i, err := strconv.Atoi(tok); err == nil {
                u := uint8(i)
                stk = append(stk, sum{u, u, u})
            }
        }
    }

    v := stk[len(stk)-1]
    stk = nil

    return v.r, v.g, v.b, nil
}
