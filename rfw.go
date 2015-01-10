package rfw

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type Rfw struct {
	Basepath string
	LastTime time.Time
	OutFile  *os.File
}

func generatePath(basepath string, t time.Time) string {
	return fmt.Sprintf("%s-%4d%02d%02d-00", basepath, t.Year(), t.Month(), t.Day())
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
	if w.OutFile == nil {
		return 0, errors.New(fmt.Sprintf("Rfw is closed. Basepath=%s", w.Basepath))
	}
	t := time.Now()
	if t.Day() != w.LastTime.Day() {
		path := generatePath(w.Basepath, t)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			return 0, err
		}
		w.OutFile.Close()
		w.OutFile = f
	}
	n, err := w.OutFile.Write(p)
	return n, err
}

func (w *Rfw) Close() {
	if w.OutFile != nil {
		w.OutFile.Close()
		w.OutFile = nil
	}
}
