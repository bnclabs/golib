// Uses text/scanner to generate tokens. This scanner can be integrated with
// parsec tool using `Scanner' interface.

package parsec

import (
    "bytes"
    "fmt"
    "github.com/prataprc/golib"
    "io"
    "regexp"
    "text/scanner"
)

const ( // Internal Commands
    cmdScan = iota
    cmdNext
    cmdPeek
    cmdBookmark
    cmdRewind
)

// Type that implements Scanner interface
type Goscan struct {
    text     []byte                    // Complete input text
    req      chan<- []int              // Request channel
    res      <-chan interface{}        // Response channel
    filename string                    // Filename to open
    debug    bool                      // If enable do debug logging
    patterns map[string]*regexp.Regexp // map of compiled regular expression
}

// Create a new instance of Goscan. Subsequently Scanner interface can be used
// on this object.
func NewGoScan(text []byte, config golib.Config) *Goscan {
    var res = make(chan interface{})
    var req = make(chan []int)
    rd := bytes.NewReader(text)
    go doscan(req, res, rd)
    debug := golib.Bool(config["debug"], false)
    return &Goscan{
        req: req, res: res, text: text, debug: debug,
        patterns: make(map[string]*regexp.Regexp),
    }
}

func (s *Goscan) Text() []byte {
    return s.text
}

func (s *Goscan) Scan() *Token {
    s.req <- []int{cmdScan}
    res := (<-s.res).(*Token)
    if s.debug {
        fmt.Println(res)
    }
    return res
}

func (s *Goscan) Next() *Token {
    s.req <- []int{cmdNext}
    res := (<-s.res).(*Token)
    if s.debug {
        fmt.Println(res)
    }
    return res
}

func (s *Goscan) Peek(offset int) *Token {
    s.req <- []int{cmdPeek, offset}
    res := (<-s.res).(*Token)
    if s.debug {
        fmt.Println(res)
    }
    return res
}

func (s *Goscan) Match(pattern string) *Token {
    var err error
    regc := s.patterns[pattern]
    if regc == nil {
        if regc, err = regexp.Compile(pattern); err == nil {
            s.patterns[pattern] = regc
        } else {
            panic(err.Error())
        }
    }
    tok := s.Peek(0)
    if regc.Match([]byte(tok.Value)) {
        return s.Scan()
    } else {
        return nil
    }
}

func (s *Goscan) Literal(name string) *Token {
    tok := s.Peek(0)
    if tok.Type == name {
        return s.Scan()
    } else {
        return nil
    }
}

func (s *Goscan) BookMark() int {
    s.req <- []int{cmdBookmark}
    return (<-s.res).(int)
}

func (s *Goscan) Rewind(offset int) {
    s.req <- []int{cmdRewind, offset}
    <-s.res
}

// This tokenizer is using text/scanner package. Make it generic so that
// parsec can be converted to a separate package.
func doscan(req <-chan []int, res chan<- interface{}, src io.Reader) {
    var s scanner.Scanner
    var curtok = 0

    s.Init(src)
    toks := fullscan(&s)
    for {
        cmd := <-req
        switch cmd[0] {
        case cmdBookmark:
            res <- curtok
        case cmdRewind:
            curtok = cmd[1]
            res <- &Token{} // Dummy
        case cmdScan:
            res <- &toks[curtok]
            curtok += 1
        case cmdPeek:
            off := cmd[1]
            if off < 0 {
                panic("Offset to peek should be 0 or more")
            }
            res <- &toks[curtok+off]
        case cmdNext:
            res <- &toks[curtok]
        default:
            panic("Unknown command to goscan")
        }
    }
}

func fullscan(s *scanner.Scanner) []Token {
    var toks = make([]Token, 0)
    for {
        tok := Token{
            Type:  scanner.TokenString(s.Scan()),
            Value: s.TokenText(),
            Pos:   s.Pos(),
        }
        toks = append(toks, tok)
        if tok.Type == "EOF" {
            break
        }
    }
    return toks
}
