package glitch

import (
    "fmt"
    "unicode"
)

type Expression string

// shunting yard algorithm
func CompileExpression(input string) (exp Expression, err error){
    defer func() {
        if r := recover(); r != nil {
            exp = ""
            err = fmt.Errorf("invalid expression: %s", input)
        }
    }()

    lastWasDigit := false
    output := ""
    opers := []byte{}

    for i := 0; i < len(input); i++ {
        tok := input[i]
        switch {
        case tok == '(':
            if lastWasDigit {
                lastWasDigit = false
                output += " "
            }
            opers = append(opers, tok)
        case tok == ')':
            if lastWasDigit {
                lastWasDigit = false
                output += " "
            }
            for {
                // pop front
                op := opers[len(opers)-1]
                opers = opers[:len(opers)-1]

                if op == '(' {
                    break
                } else {
                    output += string(op) + " "
                }
            }
        case operMap[tok] != nil:
            op := operMap[tok]
            if lastWasDigit {
                lastWasDigit = false
                output += " "
            }
            for {
                if len(opers) == 0 {
                    break
                }

                opTok := opers[len(opers)-1]
                if op2, ok := operMap[opTok]; !ok || !op2.hasPrecedence(op) {
                    break
                }

                // pop front
                opers = opers[:len(opers)-1]
                output += string(opTok) + " "
            }
            opers = append(opers, tok)
        case tok == 'c' || tok == 's' || tok == 'n' || tok == 'r':
            if lastWasDigit {
                lastWasDigit = false
                output += " "
            }
            output += string(tok) + " "
        case unicode.IsDigit(rune(tok)):
            if !lastWasDigit {
                lastWasDigit = true
            }
            output += string(tok)
        }
    }

    return Expression(output + " " + reverse(opers)), nil
}
