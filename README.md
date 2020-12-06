# BuildKit Export
> Build and export docker images to OCU for use with containerD.

## Build
Build and export an OCI image with [buildx](https://github.com/docker/buildx#buildx-build-options-path--url---):
```bash
docker buildx build --platform=linux/arm/v7 --tag hello -o type=oci,dest=- . > hello.tar
```
**NOTE**: If you have [buildkit](https://github.com/moby/buildkit) installed, you can use that.

## Import
Copy the image to the pi:
```bash
scp ./hello.tar pi@groot.local:
```

On the pi, import the image with `ctr`:
```bash
sudo ctr -n=example images import hello.tar
```
**NOTE**: We use the `example` namespace, but could just use the `default` one too.

The image should be available as `docker.io/library/hello:latest` after imported.

## Use
Copy the `main.go` file to the pi:
```go
scp ./main.go pi@groot.local:
```

Build it:
```bash
go build main.go
```

And run it:
```bash
sudo ./main
```

Or build it locally and copy it to the pi:
```bash
GOARCH=arm GOOS=linux go build -o main ./main.go
scp ./main pi@groot.local:
sudo ./main
```
