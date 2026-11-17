package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/cio"
	"github.com/containerd/containerd/v2/pkg/namespaces"
	"github.com/containerd/containerd/v2/pkg/oci"
)

func main() {
	address := flag.String("address", "/run/containerd/containerd.sock", "Path to containerd socket")
	namespace := flag.String("namespace", "example", "Containerd namespace to use")
	imageRef := flag.String("image", "docker.io/library/hello:latest", "Image reference to run")
	containerID := flag.String("container-id", "hello", "Container ID")
	snapshotID := flag.String("snapshot-id", "hello-snapshot", "Snapshot ID")
	timeout := flag.Duration("timeout", 10*time.Second, "Maximum execution timeout")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	if err := run(ctx, *address, *namespace, *imageRef, *containerID, *snapshotID); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, address, namespace, imageRef, containerID, snapshotID string) error {
	client, err := containerd.New(address)
	if err != nil {
		return fmt.Errorf("failed to connect to containerd at %q: %w", address, err)
	}
	defer client.Close()

	ctx = namespaces.WithNamespace(ctx, namespace)

	image, err := client.GetImage(ctx, imageRef)
	if err != nil {
		return fmt.Errorf("failed to get image %q: %w", imageRef, err)
	}

	container, err := client.NewContainer(
		ctx,
		containerID,
		containerd.WithImage(image),
		containerd.WithNewSnapshot(snapshotID, image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return fmt.Errorf("failed to create container %q: %w", containerID, err)
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	defer task.Delete(ctx)

	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait on task: %w", err)
	}

	if err := task.Start(ctx); err != nil {
		return fmt.Errorf("failed to start task: %w", err)
	}

	select {
	case status := <-exitStatusC:
		code, _, err := status.Result()
		if err != nil {
			return fmt.Errorf("task exit result error: %w", err)
		}
		fmt.Printf("hello exited with status: %d\n", code)
		if code != 0 {
			return fmt.Errorf("task exited with non-zero status: %d", code)
		}
	case <-ctx.Done():
		_ = task.Kill(context.Background(), syscall.SIGKILL)
		return ctx.Err()
	}

	return nil
}
