package parsec
import ("fmt"; "text/scanner")

type ParsecNode interface{}
type Parsec func() Parser
type Parser func(Scanner) ParsecNode
type Nodify func( []ParsecNode ) ParsecNode
type Scanner interface {
    Scan() Token
    Peek(int) Token
    Next() Token
    BookMark() int
    Rewind(int)
    Text() []byte
}
type Token struct {
    Type string
    Value string
    Pos scanner.Position
}
type Terminal struct {
    Name string         // typically contains terminal's token type
    Value string        // value of the terminal
    Tok Token           // Actual token obtained from the scanner
}

func docallback( callb Nodify, n []ParsecNode ) ParsecNode {
    if callb != nil {
        return callb(n)
    } else {
        return n
    }
}

func And( name string, callb Nodify, assert bool, parsecs ...Parsec ) Parsec {
    return func() Parser {
        return func( s Scanner ) ParsecNode {
            //fmt.Println(name)
            var ns = make([]ParsecNode, 0)
            bm := s.BookMark()
            for _, parsec := range parsecs {
                n := parsec()(s)
                if n == nil && assert {
                    panic( fmt.Sprintf("`And` combinator failed for %v \n", name) )
                } else if n == nil {
                    s.Rewind(bm)
                    return docallback(callb, nil)
                }
                ns = append(ns, n)
            }
            return docallback(callb, ns)
        }
    }
}

func OrdChoice(
  name string, callb Nodify, assert bool, parsecs ...Parsec ) Parsec {
    return func() Parser {
        return func(s Scanner) ParsecNode {
            var n ParsecNode
            //fmt.Println(name)
            for _, parsec := range parsecs {
                bm := s.BookMark()
                n = parsec()(s)
                if n != nil {
                    return docallback( callb, []ParsecNode{n} )
                }
                s.Rewind(bm)
            }
            if assert {
                panic(fmt.Sprintf("`OrdChoice` combinator failed for %v \n", name))
            }
            return docallback( callb, nil )
        }
    }
}

func Kleene( name string, callb Nodify, parsecs ...Parsec ) Parsec {
    var opScan, sepScan Parsec
    opScan = parsecs[0]
    if len(parsecs) >= 2 {
        sepScan = parsecs[1]
    }
    return func() Parser {
        return func(s Scanner) ParsecNode {
            var ns = make([]ParsecNode, 0)
            //fmt.Println(name)
            for {
                n := opScan()(s)
                if n == nil {
                    break
                }
                ns = append(ns, n)
                if sepScan != nil  &&  sepScan()(s) == nil {
                    break
                }
            }
            return docallback( callb, ns )
        }
    }
}

func Many( name string, callb Nodify, assert bool, parsecs ...Parsec ) Parsec {
    var opScan, sepScan Parsec
    opScan = parsecs[0]
    if len(parsecs) >= 2 {
        sepScan = parsecs[1]
    }
    return func() Parser {
        return func(s Scanner) ParsecNode {
            var ns = make([]ParsecNode, 0)
            //fmt.Println(name)
            bm := s.BookMark()
            n := opScan()(s)
            if n == nil && assert {
                panic(fmt.Sprintf("`Many` combinator failed for %v \n", name))
            } else if n == nil {
                s.Rewind(bm)
                return docallback( callb, nil )
            } else {
                for {
                    ns = append(ns, n)
                    if sepScan != nil  &&  sepScan()(s) == nil {
                        break
                    }
                    n = opScan()(s)
                    if n == nil {
                        break
                    }
                }
                return docallback( callb, ns )
            }
            return docallback(callb, nil) // Is this return required ??
        }
    }
}

func Maybe( name string, callb Nodify, parsec Parsec ) Parsec {
    return func() Parser {
        return func(s Scanner) ParsecNode {
            //fmt.Println(name)
            n := parsec()(s)
            if n == nil {
                return docallback(callb, nil)
            }
            return docallback( callb, []ParsecNode{n} )
        }
    }
}


// Parsec functions to match special strings.
func Terminalize(matchval string, n string, v string ) Parsec {
    return func() Parser {
        return func(s Scanner) ParsecNode {
            tok := s.Peek(0)
            if matchval == tok.Value {
                s.Scan()
                return &Terminal{Name:n, Value:v, Tok:tok}
            } else {
                return nil
            }
        }
    }
}

// Parsec functions to match `String`, `Char`, `Int`, `Float` literals
func Literal() Parser {
    return func(s Scanner) ParsecNode {
        tok := s.Peek(0)
        if tok.Type == "String" {
            s.Scan()
            return &Terminal{Name:tok.Type, Value:tok.Value, Tok:tok}
        } else if tok.Type == "Char" {
            s.Scan()
            return &Terminal{Name:tok.Type, Value:tok.Value, Tok:tok}
        } else if tok.Type == "Int" {
            s.Scan()
            return &Terminal{Name:tok.Type, Value:tok.Value, Tok:tok}
        } else if tok.Type == "Float" {
            s.Scan()
            return &Terminal{Name:tok.Type, Value:tok.Value, Tok:tok}
        } else {
            return nil
        }
    }
}

// Parsec function to detect end of scanner output.
func End() Parser {
    return func(s Scanner) ParsecNode {
        tok := s.Next()
        if tok.Type == "EOF" {
            return nil
        }
        return &tok
    }
}

