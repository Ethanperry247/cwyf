//go:build mage

package main

import (
	"fmt"
	"time"

	"github.com/magefile/mage/sh"
)

const (
	MAX_UNIT_TEST_DURATION_MS = 3000
	BENCHMARK_ITERATIONS = 10
)

func Start() error {
	return sh.RunV("go", "run", "./cmd/app")
}

func Setup() error {
	return sh.RunV("go", "install", "golang.org/x/perf/cmd/benchstat")
}

func Test() error {
	return sh.RunV("go", "test", "-timeout", fmt.Sprintf("%dms", MAX_UNIT_TEST_DURATION_MS), "-cover", fmt.Sprintf("--coverprofile=mage_output/coverage/coverage-%s.out", time.Now().Format(time.RFC3339Nano)), "-covermode=atomic", "-race", "./...")
}

func Bench() error {
	return sh.RunV("go", "test", "-v", "./...", "-bench=.", "-run=xxx", "-benchmem", fmt.Sprintf("-count=%d", BENCHMARK_ITERATIONS))
}