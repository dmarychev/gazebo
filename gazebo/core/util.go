package core

import "github.com/go-gl/gl/v4.6-core/gl"

func CheckError() {
	if err := gl.GetError(); err != 0 {
		panic(err)
	}
}
