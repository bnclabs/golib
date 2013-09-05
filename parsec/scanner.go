// scanner to parse terminals from input text.
package parsec

import (
    "fmt"
    "regexp"
)

var _ = fmt.Sprintf("keep 'fmt' import during debugging");

var patterns = make(map[string]*regexp.Regexp)

type Scanner struct {
    buf []byte
    cursor int
}

func NewScanner(text []byte) *Scanner {
    return &Scanner{text, 0}
}

func (s Scanner) Match(pattern string) ([]byte, *Scanner) {
    if pattern[0] != '^' {
        panic("match patterns must begin with `^`")
    }
    var regc *regexp.Regexp
    var err error

    regc = patterns[pattern]
    if regc == nil {
        if regc, err = regexp.Compile(pattern); err == nil {
            patterns[pattern] = regc
        } else {
            panic(err.Error())
        }
    }
    if token := regc.Find(s.buf[s.cursor:]); token != nil {
        s.cursor += len(token)
        return token, &s
    }
    return nil, &s
}

func (s Scanner) SkipWS() ([]byte, *Scanner) {
    return s.Match(`^[ \t\r\n]+`)
}

func (s Scanner) Endof() bool {
    if s.cursor >= len(s.buf) {
        return true
    }
    return false
}
