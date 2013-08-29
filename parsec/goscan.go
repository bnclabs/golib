// Uses text/scanner to generate tokens. This scanner can be integrated with
// parsec tool using `Scanner' interface.

package parsec
import (
    "io"; "fmt"; "bytes"; "text/scanner"; "github.com/prataprc/golib"
)

// Type that implements Scanner interface
type Goscan struct {
    // Complete input text
    text []byte
    // Request channel
    req chan<- Interface
    // Response channel
    res <-chan Interface
    // Filename to open
    filename string
    // If enable do debug logging
    debug bool
}

// Create a new instance of Goscan. Subsequently Scanner interface can be used
// on this object.
func NewGoScan(text []byte, options map[string]Interface) *Goscan {
    var res = make( chan Interface )
    var req = make( chan Interface )
    rd := bytes.NewReader( text )
    go doscan( req, res, rd )
    debug := golib.Bool(options["debug"], false)
    return &Goscan{ req: req, res: res, text:text, debug:debug }
}

// Return the plain text from input file.
func (s *Goscan) Text() []byte {
    return s.text
}

func (s *Goscan) Scan() Token {
    var cmd = make([]Interface, 1)
    cmd[0] = "scan"
    s.req<-cmd
    res := (<-s.res).(Token)
    if s.debug {
        fmt.Println(res)
    }
    return res
}

func (s *Goscan) Next() Token {
    var cmd = make([]Interface, 1)
    cmd[0] = "next"
    s.req<- cmd
    res := (<-s.res).(Token)
    //fmt.Println(res)
    return res
}

func (s *Goscan) Peek(offset int) Token {
    var cmd = make([]Interface, 2)
    cmd[0] = "peek"
    cmd[1] = offset
    s.req <- cmd
    res := (<-s.res).(Token)
    //fmt.Println(res)
    return res
}

func (s *Goscan) BookMark() int {
    var cmd = make([]Interface, 2)
    cmd[0] = "bookmark"
    s.req <- cmd
    return (<-s.res).(int)
}

func (s *Goscan) Rewind(offset int) {
    var cmd = make([]Interface, 2)
    cmd[0] = "rewind"
    cmd[1] = offset
    s.req <- cmd
    <-s.res
}


// This tokenizer is using text/scanner package. Make it generic so that
// parsec can be converted to a separate package.
func doscan( req <-chan Interface, res chan<- Interface, src io.Reader ) {
    var s scanner.Scanner
    var curtok = 0

    s.Init(src)
    toks := fullscan(&s)
    for {
        cmd := (<-req).([]Interface)
        switch cmd[0].(string) {
        case "bookmark" :
            res <- curtok
        case "rewind" :
            curtok = cmd[1].(int)
            res <- Token{} // Dummy
        case "scan" :
            res <- toks[curtok]
            curtok += 1
        case "peek" :
            off := cmd[1].(int)
            if off < 0 { panic("Offset to peek should be 0 or more") }
            res <- toks[curtok+off]
        case "next" :
            res <- toks[curtok]
        default :
            fmt.Printf("Unknown command to goscan : %v\n", cmd[0].(string))
        }
    }
}

func fullscan( s *scanner.Scanner ) []Token {
    var toks = make( []Token, 0 )
    for {
        tok := Token {
            Type: scanner.TokenString( s.Scan() ),
            Value: s.TokenText(),
            Pos: s.Pos(),
        }
        toks = append(toks, tok )
        if tok.Type == "EOF" {
            break
        }
    }
    return toks
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
        //fmt.Println("Literal")
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

