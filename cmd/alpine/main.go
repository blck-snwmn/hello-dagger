package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

func main() {
	if err := build(); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func build() error {
	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	dir := client.Host().Workdir().WithoutDirectory(".github").WithoutDirectory(".git")

	container := client.Container().From("alpine:3.16.2")

	// https://github.com/cue-lang/cue/releases/download/v0.4.3/cue_v0.4.3_linux_amd64.tar.gz
	cversion := "v0.4.3"
	ctarball := fmt.Sprintf("cue_%s_linux_amd64.tar.gz", cversion)
	path := fmt.Sprintf("https://github.com/cue-lang/cue/releases/download/%s/%s", cversion, ctarball)

	container = container.Exec(dagger.ContainerExecOpts{
		Args: []string{"wget", "-O", "/" + ctarball, path},
	})
	container = container.Exec(dagger.ContainerExecOpts{
		Args: []string{"ls", "-alt", "/"},
	})
	container = container.Exec(dagger.ContainerExecOpts{
		Args: []string{"tar", "zxf", ctarball, "-C", "/usr/local/bin"},
	})
	container = container.Exec(dagger.ContainerExecOpts{
		Args: []string{"cue", "version"},
	})
	container = container.WithMountedDirectory("/cue", dir).WithWorkdir("/cue")
	container = container.Exec(dagger.ContainerExecOpts{
		Args: []string{"cue", "vet", "sample.yaml", "check.cue"},
	})
	out, err := container.Stdout(ctx)
	fmt.Println(out)
	if err != nil {
		return err
	}
	fmt.Println(out)
	fmt.Println("end")
	return nil
}
