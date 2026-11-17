# ContainerD OCI Import

> Build and export Docker images as OCI archives for use with containerd on a Raspberry Pi.

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/rolandjitsu/containerd-oci-import/test.yml?branch=main&label=tests&style=flat-square)](https://github.com/rolandjitsu/containerd-oci-import/actions?query=workflow%3ATest)

Companion repository to the blog post: [How to install containerd on Raspberry Pi](https://rolandsdev.blog/posts/how-to-install-containerd-on-raspberry-pi/).

## Prerequisites

- [Docker](https://docs.docker.com/engine) with [Buildx](https://github.com/docker/buildx)
- [Go](https://go.dev/)
- [containerd](https://containerd.io/) (v1.7 LTS or v2.x)

## Build

### Using Docker Buildx Bake (Recommended)

Build and export the multi-platform OCI image to `./hello.tar`:

```bash
docker buildx bake
```

### Using Docker Buildx Directly

For 64-bit Raspberry Pi OS (Pi 3, 4, 5, Zero 2 W):

```bash
docker buildx build --platform=linux/arm64 --tag hello -o type=oci,dest=hello.tar .
```

For 32-bit Raspberry Pi OS (Pi 1, 2, Zero or 32-bit installs):

```bash
docker buildx build --platform=linux/arm/v7 --tag hello -o type=oci,dest=hello.tar .
```

## Import

Copy the exported OCI image to the Pi:

```bash
scp ./hello.tar pi@raspberrypi.local:
```

On the Pi, import the image using `ctr`:

```bash
sudo ctr -n=example images import hello.tar
```

> **Note**: We use the `example` namespace, but `default` can also be used.

The image will be registered as `docker.io/library/hello:latest`.

## Use

### Cross-compile for the Pi

Build the Go binary locally for ARM:

```bash
# For 64-bit ARM (ARMv8 / aarch64)
make build-arm64
scp ./main-arm64 pi@raspberrypi.local:main

# Or for 32-bit ARM (ARMv7)
make build-arm
scp ./main-armv7 pi@raspberrypi.local:main
```

### Run on the Pi

SSH into the Pi and execute:

```bash
sudo ./main
```

You should see:

```text
hello exited with status: 0
```

### Options

The Go binary supports configurable flags:

```bash
sudo ./main -address /run/containerd/containerd.sock -namespace example -image docker.io/library/hello:latest -timeout 10s
```

| Flag            | Default                           | Description                    |
| --------------- | --------------------------------- | ------------------------------ |
| `-address`      | `/run/containerd/containerd.sock` | Path to containerd UNIX socket |
| `-namespace`    | `example`                         | Containerd namespace           |
| `-image`        | `docker.io/library/hello:latest`  | Image reference to run         |
| `-container-id` | `hello`                           | ID of the created container    |
| `-snapshot-id`  | `hello-snapshot`                  | Snapshot ID                    |
| `-timeout`      | `10s`                             | Maximum execution timeout      |

## Development

```bash
make build   # Build local binary
make bake    # Build OCI image with Buildx Bake
make fmt     # Format code
make lint    # Run go vet
make clean   # Remove build artifacts
```
