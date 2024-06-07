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
	REGISTRY_NAME = "cwyfcr"
	CLUSTER_NAME = "cwyf-prod"
	RG_NAME = "CWYF"
	SUBSCRIPTION = "fcf133f3-8581-445e-b19e-a9ac059b7d36"
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

func Verify() error {
	err := sh.RunV("go", "fmt", "./...")
	if err != nil {
		return err
	}
	return sh.RunV("go", "test", "-timeout", fmt.Sprintf("%dms", MAX_UNIT_TEST_DURATION_MS), "-cover", "-covermode=atomic", "-race", "./...")
}

func Build() error {
	return sh.RunV("go", "build", "-o", "./bin/app", "./cmd/app")
}

func fmtTag(tag string) string {
	return fmt.Sprintf("%s.azurecr.io/danger-dodgers:%s", REGISTRY_NAME, tag)
}

func Docker(tag string) error {
	return sh.RunV("docker", "build", ".", "-f", "Dockerfile", "-t", fmtTag(tag))
}

func ACRLogin() error {
	return sh.RunV("az", "acr", "login", "-n", REGISTRY_NAME)
}

func BuildAndPush(tag string) error {
	err := ACRLogin()
	if err != nil {
		return err
	}

	err = Docker(tag)
	if err != nil {
		return err
	}

	return Push(tag)
}

func Push(tag string) error {
	return sh.RunV("docker", "push", fmtTag(tag))
}

func ConnectACR() error {
	return sh.RunV("az", "aks", "update", "-n", CLUSTER_NAME, "-g", RG_NAME, "--attach-acr", REGISTRY_NAME)
}

func Latest() error {
	return BuildAndPush("latest")
}

func ClusterConnect() error {
	err := sh.RunV("az", "account", "set", "--subscription", SUBSCRIPTION)
	if err != nil {
		return err
	}

	return sh.RunV("az", "aks", "get-credentials", "--resource-group", RG_NAME, "--name", CLUSTER_NAME)
}

func Deploy() error {
	return sh.RunV("helm", "upgrade", "--install", "danger-dodgers", "--create-namespace", "-n", "danger-dodgers", "./charts/environment")
}