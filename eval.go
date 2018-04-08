package glitch

import (
    "fmt"
    "math/rand"
    "strconv"
    "strings"
)

type sum struct {
    r, g, b uint8
}

func (expr Expression) evalRPN(x, y, w, h int,
                               r, sr, g, sg, b, sb uint8,
                               box []sum,
) (rr uint8, gr uint8, br uint8, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("error evaluating expression: %s", expr)
            return
        }
    }()

    stk := []sum{}
    toks := strings.Split(string(expr), " ")

    for _,tok := range toks {
        if tok == "" {
            continue
        }
        if oper := operMap[tok[0]]; oper != nil {
            a, b := stk[len(stk)-2], stk[len(stk)-1]
            stk = append(stk[:len(stk)-2], sum{oper.f(a.r, b.r),
                                               oper.f(a.g, b.g),
                                               oper.f(a.b, b.b)})
        } else if tok == "c" {
            stk = append(stk, sum{r, g, b})
        } else if tok == "n" {
            y := uint8(float64(r)*0.299) + uint8(float64(g)*0.587) + uint8(float64(b)*0.0722)
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
            stk = append(stk, sum{box[i].r, box[j].g, box[k].b})
        } else {
            if i, err := strconv.Atoi(tok); err == nil {
                u := uint8(i)
                stk = append(stk, sum{u, u, u})
            }
        }
    }

    v := stk[len(stk)-1]

    return v.r, v.g, v.b, nil
}
