package main

import (
	"io"
	"os"
	"strings"
)

func main() {
	fs, _ := os.ReadDir(".")
	out, _ := os.Create("swagger.pb.go")
	_, _ = out.Write([]byte("package pb \n\nconst (\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".json") {
			name := strings.TrimPrefix(f.Name(), "inference_engine.")
			_, _ = out.Write([]byte(strings.TrimSuffix(name, ".json") + " = `"))
			f, _ := os.Open(f.Name())
			_, _ = io.Copy(out, f)
			_, _ = out.Write([]byte("`\n"))
		}
	}
	_, _ = out.Write([]byte(")\n"))
}
