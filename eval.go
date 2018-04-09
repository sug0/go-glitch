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

    for _,tok := range expr.toks {
        if tok == "" {
            continue
        }
        if oper := operMap[tok[0]]; oper != nil {
            a, b := stk[len(stk)-2], stk[len(stk)-1]
            stk = append(stk[:len(stk)-2], sum{oper.f(a.r, b.r),
                                               oper.f(a.g, b.g),
                                               oper.f(a.b, b.b)})
        } else if tok == "c" {
            stk = append(stk, sum{box[4].R, box[4].G, box[4].B})
        } else if tok == "R" {
            stk = append(stk, sum{box[4].R, 0, 0})
        } else if tok == "G" {
            stk = append(stk, sum{0, box[4].G, 0})
        } else if tok == "B" {
            stk = append(stk, sum{0, 0, box[4].B})
        } else if tok == "Y" {
            y := uint8(float64(box[4].R)*0.299) +
                 uint8(float64(box[4].G)*0.587) +
                 uint8(float64(box[4].B)*0.0722)
            stk = append(stk, sum{y, y, y})
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
