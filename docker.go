package grm

import (
	"fmt"
	"strings"
)

func dockerInternal(rule string, version string) {
	parts := strings.Split(rule, "-")
	if len(parts) < 2 {
		panic("invalid rule")
	}
	arch := parts[len(parts)-1]
	pkgName := strings.Join(parts[:len(parts)-1], "-")

	imageTag := fmt.Sprintf("%s-%s", pkgName, arch)
	Run(nil, "docker", "build",
		"--platform", archToDockerPlatform(arch),
		"--tag", imageTag,
		"--build-arg", fmt.Sprintf("TAG=%s", version),
		"-f", fmt.Sprintf("cmd/%s/Dockerfile.%s", pkgName, arch),
		".",
	)

	src := fmt.Sprintf("/build/%s_%s_%s.tar.gz", pkgName, version, arch)
	dst := fmt.Sprintf("%s/builds/%s", Opts.SharedFolder, arch)

	DockerCopy(
		arch,
		imageTag,
		src,
		dst,
	)
}

func DockerCopy(
	arch string,
	imageTag string,
	src string,
	dst string,
) {
	container := output(nil, "docker", "container", "create", "--platform", archToDockerPlatform(arch), imageTag)

	Run(nil, "docker", "cp",
		fmt.Sprintf("%s:%s", container, src),
		dst,
	)
}

func Docker(rule string) {
	dockerInternal(rule, getTag())
}

func DockerWithCustomVersion(version string) func(rule string) {
	return func(rule string) {
		dockerInternal(rule, version)
	}
}

func archToDockerPlatform(arch string) string {
	switch arch {
	case "arm64":
		return "linux/arm64/v8"
	case "armhf":
		return "linux/arm/v7"
	case "armv6":
		return "linux/arm/v6"
	default:
		panic("unknown arch")
	}
}
