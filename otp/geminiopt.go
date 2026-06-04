package otp

import (
	"io"
)

// Буфер на 1024 байта для стековой обработки. 
// Размер оптимален: достаточно велик для снижения числа системных вызовов, 
// и достаточно мал, чтобы оставаться на стеке горутины.
const bufferSize = 1024

type CipherEncryptOpt struct {
	prng         io.Reader
	streamWriter io.Writer
}

func (c *CipherEncryptOpt) Write(p []byte) (n int, err error) {
	var buf [bufferSize]byte

	for len(p) > 0 {
		chunkSize := min(len(p), bufferSize)

		// Читаем маску из prng (обрабатываем ошибку на всякий случай)
		if _, err = io.ReadFull(c.prng, buf[:chunkSize]); err != nil {
			return n, err
		}

		// XOR-им во временный стековый буфер, чтобы не мутировать исходный p
		for i := range chunkSize {
			buf[i] ^= p[i]
		}

		// Пишем зашифрованный кусок наружу
		nw, err := c.streamWriter.Write(buf[:chunkSize])
		n += nw
		if err != nil {
			return n, err
		}

		p = p[chunkSize:]
	}

	return n, nil
}

func NewWriterOpt(w io.Writer, prng io.Reader) io.Writer {
	return &CipherEncryptOpt{
		prng:         prng,
		streamWriter: w,
	}
}

type CipherDecryptOpt struct {
	prng         io.Reader
	streamReader io.Reader
}

func (c *CipherDecryptOpt) Read(p []byte) (n int, err error) {
	n, err = c.streamReader.Read(p)
	if n == 0 {
		return
	}

	var buf [bufferSize]byte
	remaining := n
	start := 0

	// Расшифровываем считанные n байт in-place прямо в p
	for remaining > 0 {
		chunkSize := min(remaining, bufferSize)

		if _, prngErr := io.ReadFull(c.prng, buf[:chunkSize]); prngErr != nil {
			return start, prngErr
		}

		for i := range chunkSize {
			p[start+i] ^= buf[i]
		}

		start += chunkSize
		remaining -= chunkSize
	}

	return n, err
}

func NewReaderOpt(r io.Reader, prng io.Reader) io.Reader {
	return &CipherDecrypt{
		prng:         prng,
		streamReader: r,
	}
}
