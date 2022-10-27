package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	dir := client.Directory().
		WithNewFile("hello.txt", dagger.DirectoryWithNewFileOpts{
			Contents: "Hello, world!",
		}).
		WithNewFile("goodbye.txt", dagger.DirectoryWithNewFileOpts{
			Contents: "Goodbye, world!",
		})

	dirID, err := dir.ID(ctx)
	if err != nil {
		panic(err)
	}

	container := client.Container().From("alpine:3.16.2")

	container = container.WithMountedDirectory("/mnt", dirID)

	out, err := container.Exec(dagger.ContainerExecOpts{
		Args: []string{"ls", "/mnt"},
	}).Stdout().Contents(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q", out)

}
