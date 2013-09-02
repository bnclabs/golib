package parsec

// Parsec functions to match special strings.
func Tokenof(pattern string, n string) Parsec {
    return func() Parser {
        return func(s Scanner) ParsecNode {
            if tok := s.Match(pattern); tok != nil {
                return &Terminal{Name: n, Tok: tok}
            } else {
                return nil
            }
        }
    }
}

// Parsec functions to match `String`, `Char`, `Int`, `Float` literals
func Literalof(name string) Parsec {
    return func() Parser {
        return func(s Scanner) ParsecNode {
            if tok := s.Literal(name); tok != nil {
                return &Terminal{Name: tok.Type, Value: tok.Value, Tok: tok}
            } else {
                return nil
            }
        }
    }
}

// Parsec function to detect end of scanner output.
func End() Parser {
    return func(s Scanner) ParsecNode {
        tok := s.Next()
        if tok.Type == "EOF" {
            return &Terminal{Name: tok.Type, Value: tok.Value, Tok: tok}
        }
        return nil
    }
}

// Parsec function to detect end of scanner output.
func NoEnd() Parser {
    return func(s Scanner) ParsecNode {
        tok := s.Next()
        if tok.Type != "EOF" {
            return &Terminal{Name: tok.Type, Value: tok.Value, Tok: tok}
        }
        return nil
    }
}
