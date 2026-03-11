package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"
	"github.com/rwxrob/bonzai/yt"
	clip "github.com/rwxrob/clip/pkg"
)

func main() {
	cmd.Exec()
}

var cmd = &bonzai.Cmd{
	Name: `clip`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{help.Cmd, _edit, dir, data, add, list, play, download, cache, convert},
	Def:  play,
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

var _edit = &bonzai.Cmd{
	Name:  "edit",
	Short: "edit the CLIP_DATA file",
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) != 0 {
			return fmt.Errorf("takes no arguments")
		}

		path, ok := os.LookupEnv("CLIP_DATA")
		if !ok || path == "" {
			return fmt.Errorf("CLIP_DATA must be set")
		}

		return edit.Files(path)
	},
}

var dir = &bonzai.Cmd{
	Name: `dir`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		dir, ok := os.LookupEnv("CLIP_DIR")
		if !ok || dir == "" {
			return fmt.Errorf("CLIP_DIR must be set")
		}
		fmt.Println(dir)
		return nil
	},
}

var data = &bonzai.Cmd{
	Name: `data`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		data, ok := os.LookupEnv("CLIP_DATA")
		if !ok || data == "" {
			return fmt.Errorf("CLIP_DATA must be set")
		}
		fmt.Println(data)
		return nil
	},
}

var add = &bonzai.Cmd{
	Name:    "add",
	Short:   "download a clip and append it to CLIP_DATA",
	Usage:   "add <name> <youtube-url-or-id>",
	NumArgs: 2,
	Do: func(x *bonzai.Cmd, args ...string) error {

		dataPath := os.Getenv("CLIP_DATA")
		if dataPath == "" {
			return fmt.Errorf("CLIP_DATA must be set")
		}

		dir := os.Getenv("CLIP_DIR")
		if dir == "" {
			return fmt.Errorf("CLIP_DIR must be set")
		}

		name := args[0]
		id := yt.ExtractID(args[1])
		url := yt.NormalizeURL(args[1])

		opts := yt.DownloadOptions{
			URL:        url,
			OutputDir:  dir,
			OutputName: id,
		}

		fmt.Printf("%s: %s -> %s\n", name, url, filepath.Join(dir, id))

		if _, err := yt.Download(opts); err != nil {
			return err
		}

		f, err := os.OpenFile(dataPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		line := fmt.Sprintf("%s %s %d %g %g\n", name, id, 100, 0.0, 9001.0)

		if _, err := f.WriteString(line); err != nil {
			return err
		}

		return nil
	},
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

var download = &bonzai.Cmd{
	Name:    "download",
	Short:   "download a single clip into CLIP_DIR",
	Usage:   "download <youtube-id>|<name>",
	NumArgs: 1,
	Do: func(_ *bonzai.Cmd, args ...string) error {

		dir, ok := os.LookupEnv("CLIP_DIR")
		if !ok || dir == "" {
			return fmt.Errorf("CLIP_DIR must be set")
		}

		id := args[0]

		cl, _ := clip.Find(id, os.Getenv("CLIP_DATA"))
		if cl != nil {
			id = cl.ID
		}

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}

		out := filepath.Join(dir, id)

		if _, err := os.Stat(out); err == nil {
			return nil // already cached
		}

		_, err := yt.Download(yt.DownloadOptions{
			URL:        id,
			OutputDir:  dir,
			OutputName: id,
		})
		if err != nil {
			return fmt.Errorf("download %s: %w", id, err)
		}

		return nil
	},
}

var list = &bonzai.Cmd{
	Name:  "list",
	Short: "print all unique clip names",
	Usage: "list",
	Do: func(x *bonzai.Cmd, args ...string) error {

		if len(args) != 0 {
			return fmt.Errorf("takes no arguments")
		}

		path := os.Getenv("CLIP_DATA")
		if path == "" {
			return fmt.Errorf("CLIP_DATA must be set")
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		data, err := clip.Load(f)
		if err != nil {
			return err
		}

		seen := make(map[string]struct{})

		for _, clip := range *data {
			seen[clip.Name] = struct{}{}
		}

		names := make([]string, 0, len(seen))
		for n := range seen {
			names = append(names, n)
		}

		sort.Strings(names)

		for i, n := range names {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(n)
		}

		fmt.Println()

		return nil
	},
}

var play = &bonzai.Cmd{
	Name:    "play",
	Short:   "play a random clip matching a name",
	Usage:   "play <name>",
	NumArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {

		dataPath := os.Getenv("CLIP_DATA")
		if dataPath == "" {
			return fmt.Errorf("CLIP_DATA must be set")
		}

		dir := os.Getenv("CLIP_DIR")
		if dir == "" {
			return fmt.Errorf("CLIP_DIR must be set")
		}

		f, err := os.Open(dataPath)
		if err != nil {
			return err
		}
		defer f.Close()

		data, err := clip.Load(f)
		if err != nil {
			return err
		}

		name := args[0]

		var matches []*clip.Clip
		for _, clip := range *data {
			if clip.Name == name {
				matches = append(matches, clip)
			}
		}

		if len(matches) == 0 {
			return fmt.Errorf("clip not found: %s", name)
		}

		clip := matches[rand.IntN(len(matches))]
		path := filepath.Join(dir, clip.ID)

		cmd := exec.Command(
			"mpv",
			"--fs",
			fmt.Sprintf("--start=%g", clip.Start),
			fmt.Sprintf("--length=%g", clip.Duration),
			fmt.Sprintf("--volume=%d", clip.Volume),
			path,
		)

		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = os.Stderr

		return cmd.Run()
	},
}
