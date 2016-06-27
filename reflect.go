package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Copy from STDIN to STDOUT stopping only with an EOF or if it has been greater than TIMEOUT since the last copy.

func main() {
	CopyWithTimeout(os.Stdout, os.Stdin, 5*time.Second)
}

func CopyWithTimeout(dst io.Writer, src io.Reader, to time.Duration) (written int64, err error) {
	for {
		select {
		case <-time.After(to):
			fmt.Println("timed out")
			os.Exit(0)
		case j := <-read(src):
			if j.C > 0 {
				nw, ew := dst.Write(j.B[0:j.C])
				if nw > 0 {
					written += int64(nw)
				}
				if ew != nil {
					err = ew
					break
				}
				if j.C != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if j.E == io.EOF {
				break
			}
			if j.E != nil {
				err = j.E
				break
			}
		}
	}
}

func read(src io.Reader) chan bce {
	o := make(chan bce, 1)
	go func() {
		buf := make([]byte, 32*1024)
		c, e := src.Read(buf)
		o <- bce{B: buf, C: c, E: e}
	}()
	return o
}

type bce struct {
	B []byte
	C int
	E error
}
