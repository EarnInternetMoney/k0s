// +build bin

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Combo struct {
	OS   string
	ARCH string
}

var (
	Basename  = "conntroll"
	Path      = "bin"
	Delimeter = "-"
	Combos    = []Combo{
		Combo{"linux", "amd64"},
		Combo{"linux", "386"},
		Combo{"darwin", "amd64"},
	}
	GlobalEnv = []string{
		"CGO_ENABLED=0",
	}
)

func (c Combo) ReleaseName() string {
	return strings.Join([]string{Basename, c.OS, c.ARCH}, Delimeter)
}

func (c Combo) Env() []string {
	return append(
		GlobalEnv,
		fmt.Sprintf("GOOS=%s", c.OS),
		fmt.Sprintf("GOARCH=%s", c.ARCH),
	)
}

func main() {
	var (
		stripFlag bool
		upxFlag   bool
	)

	flag.BoolVar(&stripFlag, "strip", false, "strip binary")
	flag.BoolVar(&upxFlag, "upx", false, "compress binary with upx")
	flag.Parse()

	for _, c := range Combos {
		var (
			build = exec.Command(
				"go", "build",
				"-trimpath",
				"-v",
				"-o", filepath.Join(Path, c.ReleaseName()),
				".",
			)
			strip = exec.Command(
				"strip",
				filepath.Join(Path, c.ReleaseName()),
			)
			upx = exec.Command(
				"upx",
				filepath.Join(Path, c.ReleaseName()),
			)
		)

		build.Env = append(os.Environ(), c.Env()...)
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr
		log.Println(c.Env(), build)

		if err := build.Run(); err != nil {
			log.Fatalln(err)
		}

		if stripFlag && c.OS == "linux" {
			if err := strip.Run(); err != nil {
				log.Fatalln(err)
			}
		}

		if upxFlag {
			upx.Stdout = os.Stdout
			upx.Stderr = os.Stderr
			if err := upx.Run(); err != nil {
				log.Fatalln(err)
			}
		}
	}
}