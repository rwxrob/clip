package clip

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwxrob/bonzai/yt"
)

type Video struct {
	Name     string
	Volume   int
	File     string
	Start    float64
	Duration float64
}

type Data []*Clip

func (d *Data) Cache(dir string) error {
	if d == nil {
		return fmt.Errorf("nil data")
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	for i, clip := range *d {

		out := filepath.Join(dir, clip.ID)

		if _, err := os.Stat(out); err == nil {
			continue // already cached
		}

		if i > 0 {
			time.Sleep(3 * time.Second)
		}

		opts := yt.DownloadOptions{
			URL:        clip.ID,
			OutputDir:  dir,
			OutputName: clip.ID,
		}

		fmt.Printf("%s: https://youtube.com/watch?v=%s -> %s\n",
			clip.Name, opts.URL, out,
		)

		_, err := yt.Download(opts)

		if err != nil {
			return fmt.Errorf("cache %s: %w", clip.ID, err)
		}
	}

	return nil
}

func Load(r io.Reader) (*Data, error) {
	var d Data
	scanner := bufio.NewScanner(r)

	line := 0
	for scanner.Scan() {
		line++

		s := strings.TrimSpace(scanner.Text())

		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}

		var c Clip
		if err := c.UnmarshalText([]byte(s)); err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}

		d = append(d, &c)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("line %d: %w", line, err)
	}

	return &d, nil
}

type Clip struct {
	Name     string
	ID       string
	Volume   int
	Start    float64
	Duration float64
}

func (c Clip) String() string {
	b, err := c.MarshalText()
	if err != nil {
		return ""
	}
	return string(b)
}

func (c Clip) MarshalText() ([]byte, error) {
	if c.Name == "" {
		return nil, fmt.Errorf("missing Name")
	}
	if c.ID == "" {
		return nil, fmt.Errorf("missing ID")
	}

	s := fmt.Sprintf("%s %s %d %g %g",
		c.Name,
		c.ID,
		c.Volume,
		c.Start,
		c.Duration,
	)

	return []byte(s), nil
}

func (c *Clip) UnmarshalText(text []byte) error {
	if c == nil {
		return fmt.Errorf("nil Clip")
	}

	fields := strings.Fields(strings.TrimSpace(string(text)))
	if len(fields) != 5 {
		return fmt.Errorf("invalid clip text: %q", string(text))
	}

	vol, err := strconv.Atoi(fields[2])
	if err != nil {
		return fmt.Errorf("invalid volume %q: %w", fields[2], err)
	}

	start, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return fmt.Errorf("invalid start %q: %w", fields[3], err)
	}

	dur, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", fields[4], err)
	}

	c.Name = fields[0]
	c.ID = fields[1]
	c.Volume = vol
	c.Start = start
	c.Duration = dur

	return nil
}

func (v *Video) TextMarshaler() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Video) String() string {
	out := fmt.Sprintf(`%v %v %v %v %v`,
		v.Name, v.Volume, v.File, v.Start, v.Duration,
	)
	return out
}

func (v *Video) UnmarshalText(in []byte) (err error) {
	f := bytes.Fields(in)

	if len(f) != 5 {
		return fmt.Errorf(`five fields required`)
	}

	v.Name = string(f[0])
	v.File = string(f[2])

	v.Volume, err = strconv.Atoi(string(f[1]))
	if err != nil {
		return err
	}

	v.Start, err = strconv.ParseFloat(string(f[3]), 64)
	if err != nil {
		return err
	}

	v.Duration, err = strconv.ParseFloat(string(f[4]), 64)
	if err != nil {
		return err
	}

	return nil
}

type Videos []*Video

func Convert(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			fmt.Fprintln(w, line)
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Fprintln(w, line)
			continue
		}

		name := parts[0]
		score := parts[1]

		meta := strings.Split(parts[2], ",")
		if len(meta) < 3 {
			fmt.Fprintln(w, line)
			continue
		}

		id := meta[0]
		id = strings.TrimSuffix(id, ".webm")
		id = strings.TrimSuffix(id, ".mp4")
		id = strings.TrimSuffix(id, ".mkv")

		duration := meta[1]
		count := meta[2]

		fmt.Fprintf(w, "%s %s %s %s %s\n",
			name,
			id,
			score,
			duration,
			count,
		)
	}

	return scanner.Err()
}

func Find(name, path string) (*Clip, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	line := 0
	for scanner.Scan() {
		line++

		s := strings.TrimSpace(scanner.Text())

		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}

		var c Clip
		if err := c.UnmarshalText([]byte(s)); err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}

		if c.Name == name {
			return &c, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("clip not found: %s", name)
}
