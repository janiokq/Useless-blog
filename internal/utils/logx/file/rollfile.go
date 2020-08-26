package file

import (
	"bytes"
	"github.com/natefinch/lumberjack"
	"io"
	"os"
	"path/filepath"
	"time"
)

type WriteSyncer interface {
	io.Writer
	Sync() error
}

func NewRollingFile(path, srvname string, maxSize, maxAge int) WriteSyncer {
	if err := os.MkdirAll(path, 0766); err != nil {
		panic(err)
	}
	return newlumberjackWriteSyncer(&lumberjack.Logger{
		Filename:  filepath.Join(path, srvname+".log"),
		MaxSize:   maxSize, //megabytes
		MaxAge:    maxAge,  //days
		LocalTime: true,
		Compress:  false,
	})
}

type lumberjackWriteSyncer struct {
	*lumberjack.Logger
	buf       *bytes.Buffer
	logChan   chan []byte
	closeChan chan interface{}
	maxSize   int
}

func newlumberjackWriteSyncer(l *lumberjack.Logger) *lumberjackWriteSyncer {
	ws := &lumberjackWriteSyncer{
		Logger:    l,
		buf:       bytes.NewBuffer([]byte{}),
		logChan:   make(chan []byte, 5000),
		closeChan: make(chan interface{}),
		maxSize:   1024,
	}
	go ws.run()
	return ws
}

func (l *lumberjackWriteSyncer) run() {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			if l.buf.Len() > 0 {
				l.sync()
			}
		case bs := <-l.logChan:
			_, err := l.buf.Write(bs)
			if err != nil {
				continue
			}
			if l.buf.Len() > l.maxSize {
				l.sync()
			}
		case <-l.closeChan:
			l.sync()
			return
		}
	}
}

func (l *lumberjackWriteSyncer) Stop() {
	close(l.closeChan)
}

func (l *lumberjackWriteSyncer) Write(bs []byte) (int, error) {
	b := make([]byte, len(bs))
	for i, c := range bs {
		b[i] = c
	}
	l.logChan <- b
	return 0, nil
}
func (l *lumberjackWriteSyncer) Sync() error {
	return nil
}

func (l *lumberjackWriteSyncer) sync() error {
	defer l.buf.Reset()
	_, err := l.Logger.Write(l.buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
