package rfw

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestGetOutdatedPath(t *testing.T) {
	basepath := "test"
	tm, err := time.Parse("2006-01-02 15:04:05 UTC", "2017-12-12 18:40:00 UTC")
	if err != nil {
		t.Fatalf("Parse time failed: %v\n", err)
	}
	paths := []string{"test-20171212", "test-20171211", "test-20171113", "test-20171210", "test-20171207"}
	outdated := getOutdatedPath(basepath, paths, tm, 1)
	actual := []string{"test-20171113", "test-20171207", "test-20171210", "test-20171211", "test-20171212"}
	sort.Strings(outdated)
	if !reflect.DeepEqual(actual[:3], outdated) {
		t.Fatalf("expected: %v    get: %v\n", actual[:3], outdated)
	}
	outdated = getOutdatedPath(basepath, paths, tm, 10)
	if !reflect.DeepEqual(actual[:1], outdated) {
		t.Fatalf("expected: %v    get: %v\n", actual[:1], outdated)
	}
	outdated = getOutdatedPath(basepath, paths, tm, 100)
	if !reflect.DeepEqual([]string{}, outdated) {
		t.Fatalf("expected: %v    get: %v\n", []string{}, outdated)
	}
}
