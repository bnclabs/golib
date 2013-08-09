package parsec
import ("fmt"; "text/scanner")

type Parsec func() Parser
type Parser func(Scanner) INode
type Nodify func(Scanner, INode) INode
type Scanner interface {
    Scan() Token
    Peek(int) Token
    Next() Token
    BookMark() int
    Rewind(int)
    Text() string
}
type Token struct {
    Type string
    Value string
    Pos scanner.Position
}
type Terminal struct {
    Name string     // typically contains terminal's token type
    Value string    // value of the terminal
    Tok Token       // Actual token obtained from the scanner
}
type NonTerminal struct {
    Name string     // typically contains terminal's token type
    Value string    // value of the terminal
    Children []INode
}

type INode interface{ // AST functions
    Show(string)
    Repr(prefix string) string
}

var EMPTY = Terminal{Name: "EMPTY", Value:""}

func docallback( callb Nodify, s Scanner, n INode ) INode {
    if callb != nil {
        return callb(s, n)
    } else {
        return n
    }
}

func And( name string, callb Nodify, assert bool, parsecs ...Parsec ) Parsec {
    return func() Parser {
        return func( s Scanner ) INode {
            var ns = make([]INode, 0)
            bm := s.BookMark()
            for _, parsec := range parsecs {
                n := parsec()(s)
                if n == nil && assert {
                    panic( fmt.Sprintf("`And` combinator failed for %v \n", name) )
                } else if n == nil {
                    s.Rewind(bm)
                    return docallback( callb, s, nil)
                }
                ns = append(ns, n)
            }
            return docallback( callb, s, &NonTerminal{Children:ns} )
        }
    }
}

func OrdChoice(
  name string, callb Nodify, assert bool, parsecs ...Parsec ) Parsec {
    return func() Parser {
        return func(s Scanner) INode {
            var n INode
            for _, parsec := range parsecs {
                bm := s.BookMark()
                n = parsec()(s)
                if n != nil {
                    return docallback( callb, s, n )
                }
                s.Rewind(bm)
            }
            if assert {
                panic(fmt.Sprintf("`OrdChoice` combinator failed for %v \n", name))
            }
            return docallback( callb, s, nil )
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
        return func(s Scanner) INode {
            var ns = make([]INode, 0)
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
            return docallback( callb, s, &NonTerminal{Children: ns} )
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
        return func(s Scanner) INode {
            var ns = make([]INode, 0)
            bm := s.BookMark()
            n := opScan()(s)
            if n == nil && assert {
                panic(fmt.Sprintf("`Many` combinator failed for %v \n", name))
            } else if n == nil {
                s.Rewind(bm)
                return docallback( callb, s, nil )
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
                return docallback( callb, s, &NonTerminal{Children: ns} )
            }
            return docallback( callb, s, nil) // Is this return required ??
        }
    }
}

func Maybe( name string, callb Nodify, parsec Parsec ) Parsec {
    return func() Parser {
        return func(s Scanner) INode {
            n := parsec()(s)
            return docallback( callb, s, n )
        }
    }
}


// Parsec functions to match special strings.
func Terminalize(matchval string, n string, v string ) Parsec {
    return func() Parser {
        return func(s Scanner) INode {
            tok := s.Peek(0)
            if matchval == tok.Value {
                s.Scan()
                return &Terminal{Name: n, Value: v, Tok: tok}
            } else {
                return nil
            }
        }
    }
}

// Parsec functions to match `String`, `Char`, `Int`, `Float` literals
func Literal() Parser {
    return func(s Scanner) INode {
        tok := s.Peek(0)
        t := Terminal{Name: tok.Type, Value: tok.Value, Tok: tok}
        if tok.Type == "String" {
            s.Scan()
            return &t
        } else if tok.Type == "Char" {
            s.Scan()
            return &t
        } else if tok.Type == "Int" {
            s.Scan()
            return &t
        } else if tok.Type == "Float" {
            s.Scan()
            return &t
        } else {
            return nil
        }
    }
}

// Parsec function to detect end of scanner output.
func End() Parser {
    return func(s Scanner) INode {
        tok := s.Next()
        if tok.Type == "EOF" {
            return nil
        }
        return &Terminal{Name: tok.Type, Value: tok.Value, Tok: tok }
    }
}

func Error( s Scanner, str string ) {
    panic( fmt.Sprintf( "%v before %v \n", str, s.Next().Pos ))
}

// INode interface for Terminal
func (t *Terminal) Show( prefix string ) {
    fmt.Println( t.Repr(prefix) )
}
func (t *Terminal) Repr( prefix string ) string {
    return fmt.Sprintf(prefix) + fmt.Sprintf("%v : %v ", t.Name, t.Value)
}

// INode interface for NonTerminal
func (t *NonTerminal) Show( prefix string ) {
    fmt.Println( t.Repr(prefix) )
    for _, n := range t.Children {
        n.Show(prefix + "  ")
    }
}
func (t *NonTerminal) Repr( prefix string ) string {
    return fmt.Sprintf(prefix) + fmt.Sprintf("%v : %v \n", t.Name, t.Value)
}
