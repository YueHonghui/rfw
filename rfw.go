package rfw

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Rfw struct {
	lock               sync.RWMutex
	basepath           string
	lastTime           time.Time
	remainCntOfLogFile int
	outFile            *os.File
}

type RfwOption func(r *Rfw)

func WithCleanUp(remainCnt int) RfwOption {
	return func(r *Rfw) {
		r.remainCntOfLogFile = remainCnt
	}
}

func generatePath(basepath string, t time.Time) string {
	return fmt.Sprintf("%s-%4d%02d%02d", basepath, t.Year(), t.Month(), t.Day())
}

func New(basepath string) (*Rfw, error) {
	return NewWithOptions(basepath)
}

func NewWithOptions(basepath string, opts ...RfwOption) (*Rfw, error) {
	t := time.Now()
	path := generatePath(basepath, t)
	r, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	rt := &Rfw{basepath: basepath, lastTime: t, outFile: r}
	for _, o := range opts {
		o(rt)
	}
	if err = rt.checkClearLogFile(t); err != nil {
		return nil, err
	}
	return rt, nil
}

func getOutdatedPath(basepath string, paths []string, now time.Time, remain int) []string {
	sort.Strings(paths)
	edgepath := generatePath(basepath, now.AddDate(0, 0, 0-remain))
	i := sort.SearchStrings(paths, edgepath)
	if i <= 0 {
		return []string{}
	} else if i > len(paths) {
		return paths
	} else {
		return paths[:i]
	}
}

func (w *Rfw) checkClearLogFile(now time.Time) error {
	if w.remainCntOfLogFile <= 0 {
		return nil
	}
	matches, err := filepath.Glob(w.basepath + "-*")
	if err != nil {
		return err
	}
	torms := getOutdatedPath(w.basepath, matches, now, w.remainCntOfLogFile)
	for _, p := range torms {
		os.Remove(p)
	}
	return nil
}

func (w *Rfw) Write(p []byte) (int, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	if w.outFile == nil {
		return 0, errors.New(fmt.Sprintf("Rfw is closed. Basepath=%s", w.basepath))
	}
	t := time.Now()
	if t.YearDay() != w.lastTime.YearDay() || t.Year() != w.lastTime.Year() {
		needcheck := false
		w.lock.RUnlock()
		w.lock.Lock()
		if t.YearDay() != w.lastTime.YearDay() || t.Year() != w.lastTime.Year() {
			needcheck = true
			path := generatePath(w.basepath, t)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
			if err != nil {
				w.lock.Unlock()
				w.lock.RLock()
				return 0, err
			}
			w.outFile.Close()
			w.outFile = f
			w.lastTime = t
		}
		w.lock.Unlock()
		w.lock.RLock()
		if needcheck {
			w.checkClearLogFile(t)
		}
	}
	return w.outFile.Write(p)
}

func (w *Rfw) Close() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.outFile != nil {
		w.outFile.Close()
		w.outFile = nil
	}
	return
}
