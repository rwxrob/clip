package clip

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Video struct {
	Name     string
	Volume   int
	File     string
	Start    float64
	Duration float64
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
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		name := parts[0]
		score := parts[1]

		meta := strings.Split(parts[2], ",")
		if len(meta) < 3 {
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
