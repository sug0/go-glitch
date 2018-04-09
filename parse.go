package glitch

import (
    "fmt"
    "strings"
    "unicode"
)

type Expression struct {
    infix string
    toks  []string
}

func (expr *Expression) String() string {
    return expr.infix
}

// shunting yard algorithm
func CompileExpression(input string) (exp *Expression, err error){
    defer func() {
        if r := recover(); r != nil {
            exp = nil
            err = fmt.Errorf("invalid expression: %s", input)
        }
    }()

    lastWasDigit := false
    output := ""
    opers := make([]byte, 0, len(input))

    for i := 0; i < len(input); i++ {
        tok := input[i]
        switch {
        default:
            return nil, fmt.Errorf("invalid expression: %s", input)
        case isWhitespace(tok):
            continue
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
        case validTok(tok):
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

    return &Expression{
        infix: input,
        toks: strings.Split(output + " " + reverse(opers), " "),
    }, nil
}

func isWhitespace(tok byte) bool {
    return tok == ' ' || tok == '\n' || tok == '\t'
}

func validTok(tok byte) bool {
    return tok == 'c' || tok == 's' || tok == 'Y' ||
           tok == 'r' || tok == 'x' || tok == 'y' ||
           tok == 'N' || tok == 'R' || tok == 'G' ||
           tok == 'B'
}
