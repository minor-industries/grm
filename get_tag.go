package grm

import (
	"strings"
	"sync"
)

var getTagOnce sync.Once
var tag_ string

func getTag() string {
	getTagOnce.Do(func() {
		args := []string{"describe", "--tags", "--dirty"}

		if !Opts.AllowNoTag {
			args = append(args, "--exact-match")
		}

		tag_ = output(nil, "git", args...)
		if tag_ == "" {
			panic("no tag")
		}
		if strings.Contains(tag_, "-dirty") {
			if !Opts.AllowDirty {
				panic("dirty")
			}
		}
	})
	return tag_
}
