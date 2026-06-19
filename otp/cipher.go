//go:build !solution

package otp

import (
	"io"
)

type CipherEncrypt struct {
	prng         io.Reader
	streamWriter io.Writer
}

func (c *CipherEncrypt) Write(p []byte) (n int, err error) {
	res := make([]byte, len(p))
	_, _ = io.ReadFull(c.prng, res)
	for i := range p {
		res[i] = p[i] ^ res[i]
	}
	return c.streamWriter.Write(res)
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &CipherEncrypt{
		prng:         prng,
		streamWriter: w,
	}
}

type CipherDecrypt struct {
	prng         io.Reader
	streamReader io.Reader
}

func (c *CipherDecrypt) Read(p []byte) (n int, err error) {
	n, err = c.streamReader.Read(p)
	if n == 0 {
		return
	}

	res := make([]byte, n)
	_, _ = io.ReadFull(c.prng, res)
	for i := 0; i < n; i++ {
		p[i] ^= res[i]
	}
	return
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &CipherDecrypt{
		prng:         prng,
		streamReader: r,
	}
}
