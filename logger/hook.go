package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"
)

type StreamLogHook struct {
	ch            chan string
	handler       func(msg string)
	stopFunc      func()
	streamEnabled bool
	stdio         io.Writer
	lock          *sync.Mutex
}

func NewStreamLogHook(handler func(msg string), stopFunc func()) *StreamLogHook {
	return &StreamLogHook{
		ch:            make(chan string, 4096),
		handler:       handler,
		streamEnabled: true,
		stdio:         bufio.NewWriter(os.Stdout),
		stopFunc:      stopFunc,
		lock:          &sync.Mutex{},
	}
}

func (s *StreamLogHook) Fire(entry *logrus.Entry) error {
	if !s.streamEnabled {
		return nil
	}
	str, _ := entry.String()
	s.ch <- str
	return nil
}

func (s *StreamLogHook) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.streamEnabled = false
	close(s.ch)
	s.stopFunc()
	return nil
}

func (s *StreamLogHook) Send() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("\n--------------------\ncatch panic !!! \nhandle server message error: %v, stack: %s\n--------------------\n", err, debug.Stack())
		}
	}()
	for {
		if !s.streamEnabled {
			return
		}
		msg, ok := <-s.ch
		if !ok {
			return
		}
		s.handler(msg)
	}
}

func (s *StreamLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
