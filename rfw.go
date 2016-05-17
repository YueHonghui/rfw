package rfw

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type Rfw struct {
	lock     sync.RWMutex
	Basepath string
	LastTime time.Time
	OutFile  *os.File
}

func generatePath(basepath string, t time.Time) string {
	return fmt.Sprintf("%s-%4d%02d%02d", basepath, t.Year(), t.Month(), t.Day())
}

func New(basepath string) (*Rfw, error) {
	t := time.Now()
	path := generatePath(basepath, t)
	r, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	return &Rfw{Basepath: basepath, LastTime: t, OutFile: r}, nil
}

func (w *Rfw) Write(p []byte) (int, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	if w.OutFile == nil {
		return 0, errors.New(fmt.Sprintf("Rfw is closed. Basepath=%s", w.Basepath))
	}
	t := time.Now()
	if t.Day() != w.LastTime.Day() {
		w.lock.RUnlock()
		w.lock.Lock()
		if t.Day() != w.LastTime.Day() {
			path := generatePath(w.Basepath, t)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
			if err != nil {
				w.lock.Unlock()
				w.lock.RLock()
				return 0, err
			}
			w.OutFile.Close()
			w.OutFile = f
		}
		w.lock.Unlock()
		w.lock.RLock()
	}
	n, err := w.OutFile.Write(p)
	return n, err
}

func (w *Rfw) Close() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.OutFile != nil {
		w.OutFile.Close()
		w.OutFile = nil
	}
	return
}
