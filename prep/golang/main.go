package main

import (
	"csi-prep-golang/pkg/config"
	"csi-prep-golang/pkg/idx"
	"fmt"
	"io/ioutil"
	"os"
)

var help = `
Usage:
	mirror <URL> <OPT>

Options:
	-o	output directory for the local website mirror
`

func usage(ec int) {
	fmt.Println(help)
	os.Exit(ec)
}

func main() {
	cfg := getConfig()
	_ = os.Mkdir(cfg.OutputDirectory, 0777)
	mirror := idx.New(fileWriter{cfg.OutputDirectory})
	err := mirror.Index(cfg.URL)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%#v\n", mirror)
	fmt.Printf("%#v\n", mirror.Seen)
}

func getConfig() config.Config {
	baseDir := "/Users/cjapel/misc/csi/prep/golang/out/"
	if len(os.Args) < 2 {
		usage(1)
	}
	targetURL := os.Args[1]
	cfg := config.FromFlags()
	if err := cfg.SetURL(targetURL); err != nil {
		fmt.Println("Invalid URL: " + targetURL)
		usage(1)
	}
	if cfg.OutputDirectory == "" {
		cfg.OutputDirectory = fmt.Sprintf(
			"%s/%s/", baseDir, cfg.URL.Hostname(),
		)
	}
	return cfg
}

type fileWriter struct {
	Path string
}

func (w fileWriter) Write(fileName string, data []byte) error {
	fullPath := w.Path + "_" + fileName
	// fmt.Println("Writing to " + fullPath)
	return ioutil.WriteFile(fullPath, data, 0644)
}
