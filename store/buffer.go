package store

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	buffer = list.New()
	writer io.Writer
	mu     sync.Mutex
)

// StartBuffering takes io.Reader and reads its content storing it to internal buffer
// It allows only `maxBufferSize` elements in the buffer. It drops the logs in FIFO manner.
func StartBuffering(inputReader io.Reader, uiWriter io.Writer, maxBufferSize int) {
	writer = uiWriter
	scanner := bufio.NewScanner(inputReader)

	// We prepend NL character before the logs lines and we don't want to have first empty line.
	firstLine := true

	for scanner.Scan() {
		if buffer.Len() >= maxBufferSize {
			buffer.Remove(buffer.Back())
		}

		line := scanner.Text()
		writeLine(line, firstLine)
		firstLine = false
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

// writeLine handles syncing with UI goroutine.
// It stores the line to internal buffer and writes it to UI if it matches current filter.
func writeLine(line string, firstLine bool) {
	mu.Lock()
	defer mu.Unlock()

	buffer.PushFront(line)

	result, err := isLineMatching(line)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error filtering line:", err)
		return
	}

	if result {

		if firstLine {
			_, err = fmt.Fprint(writer, line)
			firstLine = false
		} else {
			_, err = fmt.Fprint(writer, "\n"+line)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "error writing to ui:", err)
		}
	}
}
