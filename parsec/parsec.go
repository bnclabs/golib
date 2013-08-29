// Package parsec implements a library of parser-combinators using basic
// recognizers like,
//      And, OrdChoice, Kleene, Many and Maybe.
package parsec
import (
    "fmt"
    "text/scanner"
)

type Interface interface{}
type ParsecNode interface{}             // Can be used to construct AST.
type Parsec func() Parser               // Combinable parsers.
type Parser func(Scanner) ParsecNode    // Lazy evaluation of combinable parsers
type Nodify func( []ParsecNode ) ParsecNode

// lexer tool for parser functions. Make sense only for terminal recognizers.
// Non-terminal recognizers simply pass them down. A default scanner is
// available using golang's text/scanner package.
type Token struct {
    Type string
    Value string
    Pos scanner.Position
}
type Scanner interface {
    // Scan will read the next token from input stream and return a `Token`
    // instance.
    Scan() Token

    // Peek will lookhead for nth token from input stream and return a `Token`
    // instance. This function does not consume input.
    Peek(int) Token

    // Same as Peek(1).
    Next() Token

    // Bookmark returns a token offset. When parsec tool decides to backtrack,
    // this toke offset can be used to rewind back the token stream.
    BookMark() int
    Rewind(int)

    // Return the input text stream as byte-slice.
    Text() []byte
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
            if len(ns) == 0 { 
                return docallback(callb, nil)
            } else {
                return docallback(callb, ns)
            }
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
            if len(ns) == 0 { 
                return docallback(callb, nil)
            } else {
                return docallback(callb, ns)
            }
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

