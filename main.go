package main

import (
	"io"
	"os"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
	clip "github.com/rwxrob/clip/pkg"
)

func main() {
	cmd.Exec()
}

var cmd = &bonzai.Cmd{
	Name: `clip`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{help.Cmd, edit, dir, data, add, list, play, download, convert},
}

var convert = &bonzai.Cmd{
	Name:    `convert`,
	Short:   `convert file/stdin of old data to new`,
	MaxArgs: 1,
	Cmds:    []*bonzai.Cmd{help.Cmd},
	Comp:    comp.FileDir,

	Do: func(_ *bonzai.Cmd, args ...string) error {
		var r io.Reader

		switch len(args) {
		case 0:
			r = os.Stdin

		case 1:
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()
			r = f
		}

		if err := clip.Convert(r, os.Stdout); err != nil {
			return err
		}

		return nil
	},
}

var edit = &bonzai.Cmd{
	Name: `edit`,
	Do:   bonzai.Wood,
}

var dir = &bonzai.Cmd{
	Name: `dir`,
	Do:   bonzai.Wood,
}

var data = &bonzai.Cmd{
	Name: `data`,
	Do:   bonzai.Wood,
}

var add = &bonzai.Cmd{
	Name: `add`,
	Do:   bonzai.Wood,
}

var download = &bonzai.Cmd{
	Name: `download`,
	Do:   bonzai.Wood,
}

var list = &bonzai.Cmd{
	Name: `list`,
	Do: func(*bonzai.Cmd, ...string) error {
		// TODO
		var vid clip.Video
		_ = vid
		return nil
	},
}

var play = &bonzai.Cmd{
	Name: `play`,
	Do:   bonzai.Wood,
}
