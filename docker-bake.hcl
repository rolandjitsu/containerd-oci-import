variable "DEBIAN_VERSION" {
  default = "bookworm-slim"
}

variable "PLATFORMS" {
  default = ["linux/arm64", "linux/arm/v7", "linux/amd64"]
}

group "default" {
  targets = ["image"]
}

target "image" {
  dockerfile = "Dockerfile"
  tags       = ["hello:latest"]
  platforms  = PLATFORMS
  output     = ["type=oci,dest=hello.tar"]
  args = {
    DEBIAN_VERSION = "${DEBIAN_VERSION}"
  }
}
