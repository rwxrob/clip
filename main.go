package main

import (
	"fmt"
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
	Cmds: []*bonzai.Cmd{help.Cmd, edit, dir, data, add, list, play, cache, convert},
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

var cache = &bonzai.Cmd{
	Name:   `cache`,
	Short:  `cache all clip videos into CLIP_DIR`,
	Usage:  `CLIP_DATA=/path/to/data CLIP_DIR=/path/to/cache {{cmd .Name}}`,
	NoArgs: true,
	Long: `
Loads all clip metadata from CLIP_DATA and caches every referenced
YouTube video into CLIP_DIR.
`,
	Do: func(*bonzai.Cmd, ...string) error {

		dataPath, ok := os.LookupEnv("CLIP_DATA")
		if !ok || dataPath == "" {
			return fmt.Errorf("CLIP_DATA must be set")
		}

		cacheDir, ok := os.LookupEnv("CLIP_DIR")
		if !ok || cacheDir == "" {
			return fmt.Errorf("CLIP_DIR must be set")
		}

		f, err := os.Open(dataPath)
		if err != nil {
			return fmt.Errorf("open CLIP_DATA %q: %w", dataPath, err)
		}
		defer f.Close()

		data, err := clip.Load(f)
		if err != nil {
			return fmt.Errorf("load CLIP_DATA %q: %w", dataPath, err)
		}

		if err := data.Cache(cacheDir); err != nil {
			return fmt.Errorf("cache to CLIP_DIR %q: %w", cacheDir, err)
		}

		return nil
	},
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
