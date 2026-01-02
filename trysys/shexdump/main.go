package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {

	var src io.Reader

	switch len(os.Args) {
	case 1:
		src = os.Stdin
	case 2:
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s: %v", os.Args[1], err)
			os.Exit(1)
		}
		defer f.Close()
		src = f
	default:
		fmt.Fprintf(os.Stderr, "Usage: %s [filename]", os.Args[0])
		os.Exit(1)
	}

	if err := hexdump(os.Stdout, src); err != nil {
		fmt.Fprintf(os.Stderr, "hexdump: %v", err)
		os.Exit(1)
	}
}

func hexdump(w io.Writer, src io.Reader) error {
	r := bufio.NewReader(src)
	for {
		var raw [16]byte

		encoded := make([]byte, 0, 16*3+1+1)

		n, err := io.ReadFull(r, raw[:])
		const hex = "0123456789abcdef"
		if n != 0 {
			for i := range min(n, 8) {
				encoded = append(encoded, hex[raw[i]>>4], hex[raw[i]&0x0f], ' ')
			}

			encoded = append(encoded, ' ')

			for i := 8; i < min(n, 16); i++ {
				encoded = append(encoded, hex[raw[i]>>4], hex[raw[i]&0x0f], ' ')
			}
			encoded[len(encoded)-1] = '\n'

			if _, err := w.Write(encoded); err != nil {
				return err
			}
		}

		if err == io.ErrUnexpectedEOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}
