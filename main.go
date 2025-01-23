package main

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	clip "github.com/rwxrob/clip/pkg"
)

func main() {
	cmd.Exec()
}

var cmd = &bonzai.Cmd{
	Name: `clip`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{edit, dir, data, add, list, play},
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
