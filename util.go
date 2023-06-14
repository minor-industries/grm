package grm

import (
	"fmt"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var top string

func init() {
	var err error
	top, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func noError(err error) {
	if err != nil {
		panic(err)
	}
}

func Cd(dir string, callback func()) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.Chdir(pwd)
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("cd", dir)
	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	callback()
}

func Run(env []string, exe string, args ...string) {
	fmt.Println(strings.TrimSpace(strings.Join([]string{
		strings.Join(env, " "),
		exe,
		strings.Join(args, " "),
	}, " ")))

	cmd := exec.Command(exe, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	noError(cmd.Run())
}

func RunDocker(container string, cmd ...string) {
	pwd, err := os.Getwd()
	noError(err)

	args := []string{
		"run",
		"-i",
		"--rm",
		"-v", fmt.Sprintf("%s:%s", pwd, pwd),
		"-w", pwd,
		container,
	}

	args = append(args, cmd...)
	Run(nil, "docker", args...)
}

func output(env []string, exe string, args ...string) string {
	fmt.Println(strings.TrimSpace(strings.Join([]string{
		strings.Join(env, " "),
		exe,
		strings.Join(args, " "),
	}, " ")))
	cmd := exec.Command(exe, args...)
	cmd.Env = append(os.Environ(), env...)
	outBytes, err := cmd.CombinedOutput()
	out := strings.TrimSpace(string(outBytes))

	if err != nil {
		noError(fmt.Errorf("%s: %s", err.Error(), out))
	}

	return out
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func FindUnexpected() {
	err := filepath.WalkDir(filepath.Join(Opts.SharedFolder, "builds/"), func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if strings.Contains(path, "-dirty") {
			panic("dirty file: " + path)
		}

		switch filepath.Ext(path) {
		case ".gz", ".minisig", ".deb":
			return nil
		default:
			panic("unexpected file: " + path)
		}
	})
	if err != nil {
		panic(err)
	}
}

func LsUnsigned() []string {
	var unsigned []string
	err := filepath.WalkDir(filepath.Join(Opts.SharedFolder, "builds/"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)

		switch ext {
		case ".gz", ".deb":
		// pass
		default:
			return nil
		}

		if FileExists(path + ".minisig") {
			return nil
		}

		unsigned = append(unsigned, path)

		return nil
	})
	noError(err)

	return unsigned
}

func Steps(steplist ...func(rule string)) func(string) {
	return func(rule string) {
		for _, step := range steplist {
			step(rule)
		}
	}
}
