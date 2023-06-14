package grm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func buildSingle(name, version, arch string, env []string) string {
	tmp, err := ioutil.TempDir("", "")
	noError(err)
	defer func() {
		noError(os.RemoveAll(tmp))
	}()

	Cd(filepath.Join("cmd", name), func() {
		Run(env,
			"go",
			"build",
			"-o", filepath.Join(tmp, "bin")+"/",
			".",
		)
	})

	outputFile := fmt.Sprintf("%s/%s_%s_%s.tar.gz", arch, name, version, arch)

	Cd(tmp, func() {
		Run(nil,
			"tar",
			"-czv",
			"-f", filepath.Join(Opts.SharedFolder, "builds", outputFile),
			"bin",
		)
	})

	built := filepath.Join(Opts.SharedFolder, "builds", outputFile)

	return built
}

type target struct {
	name      string
	platforms []builder
}

func build(t target) {
	var allBuilt []string
	for _, p := range t.platforms {
		built := p.Build(t.name)
		if built != "" {
			allBuilt = append(allBuilt, built)
		}
	}
}
