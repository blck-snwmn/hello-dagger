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

	dir := client.Host().Workdir().Read().WithoutDirectory(".github").WithoutDirectory(".git")
	dirID, err := dir.ID(ctx)
	if err != nil {
		return err
	}
	container := client.Container().From("cuelang/cue:0.4.3")
	container = container.WithMountedDirectory("/cue", dirID).WithWorkdir("/cue")
	_, err = container.Exec(dagger.ContainerExecOpts{
		Args: []string{"vet", "sample.yaml", "check.cue"},
	}).ExitCode(ctx)
	if err != nil {
		return err
	}
	fmt.Println("end")
	return nil
}
