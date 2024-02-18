package main

import (
	"trintech/review/tools/protoc-gen-custom/internal"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	opt := protogen.Options{}
	internal.Run(opt)
}
