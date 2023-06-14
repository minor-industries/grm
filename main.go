package grm

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"os"
	"sort"
)

var Opts struct {
	AllowDirty   bool   `long:"allow-dirty"`
	AllowNoTag   bool   `long:"allow-no-tag"`
	SharedFolder string `long:"shared-folder" default:"$HOME/shared"`
}

func Main(rules map[string]func(rule string)) {
	args, err := flags.Parse(&Opts)
	if err != nil {
		panic(errors.Wrap(err, "parse flags"))
	}

	Opts.SharedFolder = os.ExpandEnv(Opts.SharedFolder)

	if len(os.Args) == 1 {
		var ruleList []string
		for rule, _ := range rules {
			ruleList = append(ruleList, rule)
		}
		sort.Strings(ruleList)
		for _, rule := range ruleList {
			fmt.Println(rule)
		}
		return
	}

	for _, target := range args {
		f := rules[target]
		f(target)
	}
}
