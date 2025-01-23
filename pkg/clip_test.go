package clip_test

import (
	"fmt"

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

func ExampleVideos_UnmarshalText() {
	vids := `

	vids.UnmarshalText([]byte(in))
	fmt.Println(vid)
	// Output:
	// fox 100 Un2WjSjriEA.webm 122.8 22
}
