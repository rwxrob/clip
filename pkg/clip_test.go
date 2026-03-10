package clip_test

import (
	"bytes"
	"fmt"
	"strings"

	clip "github.com/rwxrob/clip/pkg"
)

func ExampleVideo_UnmarshalText() {
	in := `fox 100 Un2WjSjriEA.webm 122.8 22`

	vid := new(clip.Video)
	vid.UnmarshalText([]byte(in))
	fmt.Println(vid)
	// Output:
	// fox 100 Un2WjSjriEA.webm 122.8 22
}

func ExampleConvert() {
	input := strings.NewReader(
		"pleaseclap 100 OUXvrWeQU0g.webm,21.6,4\n",
	)

	var out bytes.Buffer

	_ = clip.Convert(input, &out)

	fmt.Print(out.String())

	// Output:
	// pleaseclap OUXvrWeQU0g 100 21.6 4
}

func ExampleClip_MarshalText() {
	c := clip.Clip{
		Name:     "pleaseclap",
		ID:       "OUXvrWeQU0g",
		Volume:   100,
		Start:    21.6,
		Duration: 4,
	}

	b, _ := c.MarshalText()
	fmt.Println(string(b))

	// Output:
	// pleaseclap OUXvrWeQU0g 100 21.6 4
}

func ExampleClip_UnmarshalText() {
	var c clip.Clip

	_ = c.UnmarshalText([]byte("pleaseclap OUXvrWeQU0g 100 21.6 4"))

	fmt.Println(c.Name)
	fmt.Println(c.ID)
	fmt.Println(c.Volume)
	fmt.Println(c.Start)
	fmt.Println(c.Duration)

	// Output:
	// pleaseclap
	// OUXvrWeQU0g
	// 100
	// 21.6
	// 4
}

func ExampleClip_String() {
	c := clip.Clip{
		Name:     "pleaseclap",
		ID:       "OUXvrWeQU0g",
		Volume:   100,
		Start:    21.6,
		Duration: 4,
	}

	fmt.Println(c.String())

	// Output:
	// pleaseclap OUXvrWeQU0g 100 21.6 4
}
