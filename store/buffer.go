package store

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/auhau/gredux"
	"github.com/auhau/loggy/state/actions"
	"github.com/rivo/tview"
	"io"
	"os"
	"sync"
)

var (
	buffer = list.New()
	mu     sync.Mutex
)

// StartBuffering takes io.Reader, reads its content and stores it to internal buffer
// It allows only `maxBufferSize` elements in the buffer. It drops the logs in FIFO manner.
func StartBuffering(inputReader io.Reader, app *tview.Application, stateStore *gredux.Store, maxBufferSize int) {
	scanner := bufio.NewScanner(inputReader)

	for scanner.Scan() {
		line := scanner.Text()
		processLine(line, app, stateStore, maxBufferSize)
	}

	if err := scanner.Err(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, "reading standard input:", err)
		if err != nil {
			panic(err) // If we can't print errors to STDERR we gonna panic
		}
	}
}

// processLine handles syncing and messaging with UI goroutine.
// It stores the line to internal buffer and sent out DropLogLine and AddLogLine actions accordingly.
func processLine(line string, app *tview.Application, stateStore *gredux.Store, maxBufferSize int) {
	mu.Lock()
	defer mu.Unlock()

	if buffer.Len() >= maxBufferSize {
		buffer.Remove(buffer.Back())
		app.QueueUpdateDraw(func() {
			stateStore.Dispatch(actions.DropLogLine(fmt.Sprint(buffer.Back().Value)))
		})
	}

	buffer.PushFront(line)

	app.QueueUpdateDraw(func() {
		stateStore.Dispatch(actions.AddLogLine(line))
	})
}
