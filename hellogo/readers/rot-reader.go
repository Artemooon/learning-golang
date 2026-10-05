package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (r13Reader *rot13Reader) Read(p []byte) (n int, err error) {
	n, err = r13Reader.r.Read(p)

	for i := 0; i < n; i++ {
		switch {
			case 'A' <= p[i] && p[i] <= 'M':
			p[i] += 13
			case 'N' <= p[i] && p[i] <= 'Z':
			p[i] -= 13
			case 'a' <= p[i] && p[i] <= 'm':
			p[i] += 13
			case 'n' <= p[i] && p[i] <= 'z':
			p[i] -= 13
		}
	}

	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
