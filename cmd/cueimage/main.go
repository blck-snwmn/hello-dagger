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

	dir := client.Host().Directory(".").WithoutDirectory(".github").WithoutDirectory(".git")

	container := client.Container().From("cuelang/cue:0.4.3")
	container = container.WithMountedDirectory("/cue", dir).WithWorkdir("/cue")
	_, err = container.WithExec([]string{"vet", "sample.yaml", "check.cue"}).Sync(ctx)
	if err != nil {
		return err
	}
	fmt.Println("end")
	return nil
}
