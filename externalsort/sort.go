//go:build !solution

package externalsort

import (
	"bytes"
	"container/heap"
	"io"
	"os"
	"sort"
)

type LineRead struct {
	reader io.Reader
	buff   []byte
	eofErr error
}

func (lr *LineRead) ReadLine() (string, error) {
	firstEndLine := bytes.IndexByte(lr.buff, '\n')
	if firstEndLine != -1 {
		line := string(lr.buff[:firstEndLine])
		lr.buff = lr.buff[firstEndLine+1:]
		return line, nil
	}
	if lr.eofErr != nil {
		var line string
		if len(lr.buff) != 0 {
			line = string(lr.buff)
		}
		lr.buff = nil
		return line, io.EOF
	}

	for {
		localBuff := make([]byte, 1024)
		n, err := lr.reader.Read(localBuff)
		if err != nil && err != io.EOF {
			return "", err
		}
		localBuff = localBuff[:n]
		lr.buff = append(lr.buff, localBuff...)
		if err == io.EOF {
			lr.eofErr = io.EOF
			return lr.ReadLine() // рекурсия глубиной 1: buff уже содержит все данные
		}
		if bytes.IndexByte(localBuff, '\n') != -1 {
			return lr.ReadLine() // рекурсия глубиной 1: buff уже содержит все данные
		}
	}
}

func NewReader(r io.Reader) LineReader {
	return &LineRead{
		reader: r,
		buff:   make([]byte, 0, 1024),
	}
}

type LineWrite struct {
	io.Writer
}

func (lw *LineWrite) Write(l string) error {
	_, err := lw.Writer.Write([]byte(l + "\n"))
	return err
}

func NewWriter(w io.Writer) LineWriter {
	return &LineWrite{
		Writer: w,
	}
}

func Merge(w LineWriter, readers ...LineReader) error {
	pq := make(PriorityQueue, 0, len(readers))
	heap.Init(&pq)
	for i := range readers {
		if err := pushIfNotEmpty(&pq, readers[i], i); err != nil {
			return err
		}
	}

	for len(pq) > 0 {
		item := heap.Pop(&pq).(*Item)
		err := w.Write(item.value)
		if err != nil {
			return err
		}
		if err := pushIfNotEmpty(&pq, readers[item.id], item.id); err != nil {
			return err
		}
	}
	return nil
}

func pushIfNotEmpty(pq *PriorityQueue, reader LineReader, id int) error {
	line, err := reader.ReadLine()
	if err != nil && err != io.EOF {
		return err
	}
	if err == io.EOF && len(line) == 0 {
		return nil
	}
	heap.Push(pq, &Item{id: id, value: line})
	return nil
}

func Sort(w io.Writer, in ...string) error {
	files := make([]*os.File, 0, len(in))
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()
	readers := make([]LineReader, 0, len(in))

	for _, path := range in {
		if err := sortFile(path); err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		files = append(files, f)
		readers = append(readers, NewReader(f))
	}

	return Merge(NewWriter(w), readers...)
}

func sortFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	reader := NewReader(f)
	var lines []string
	for {
		line, readErr := reader.ReadLine()
		if readErr == io.EOF {
			if len(line) != 0 {
				lines = append(lines, line)
			}
			break
		}
		if readErr != nil {
			return readErr
		}
		lines = append(lines, line)
	}
	f.Close()
	sort.Strings(lines)

	f, err = os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := NewWriter(f)
	for _, line := range lines {
		err = writer.Write(line)
		if err != nil {
			return err
		}
	}
	return nil
}
