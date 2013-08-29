package golib
import (
   "os"; "bufio"
)

func Readlines( filename string ) []string {
    var lines = make( []string, 0 )
    fd, _ := os.Open(filename)
    defer func(){ fd.Close() }()
    scanner := bufio.NewScanner(fd)
    for scanner.Scan() {
        lines = append( lines, string(scanner.Bytes()) )
    }
    return lines
}

