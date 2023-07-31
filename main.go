package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pivotal-cf/jhanda"
	"github.com/pivotal-cf/winfs-injector/tile"
	"github.com/pivotal-cf/winfs-injector/winfsinjector"
)

const usageText = `winfs-injector injects the Windows 2016 root file system into the Windows 2016 Runtime Tile.

Usage: winfs-injector
  --input-tile, -i			path to input tile (example: /path/to/input.pivotal)
  --output-tile, -o			path to output tile (example: /path/to/output.pivotal)
  --preserve-extracted, -p	preserve the files created during the tile extraction process (useful for debugging)
  --registry, -r			path to docker registry (example: /path/to/registry, default: "https://registry.hub.docker.com")
  --help, -h				prints this usage information
`

func main() {
	var arguments struct {
		InputTile         string `short:"i" long:"input-tile"`
		OutputTile        string `short:"o" long:"output-tile"`
		PreserveExtracted bool   `short:"p" long:"preserve-extracted"`
		Registry          string `short:"r" long:"registry" default:"https://registry.hub.docker.com"`
		Help              bool   `short:"h" long:"help"`
	}

	_, err := jhanda.Parse(&arguments, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if arguments.Help {
		printUsage()
		return
	}

	var tileInjector = tile.NewTileInjector()
	var zipper = tile.NewZipper()
	var releaseCreator = winfsinjector.ReleaseCreator{}

	wd, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "tile extraction directory: %s\n", wd)
	if !arguments.PreserveExtracted {
		defer os.RemoveAll(wd)
	}

	app := winfsinjector.NewApplication(releaseCreator, tileInjector, zipper)

	err = app.Run(arguments.InputTile, arguments.OutputTile, arguments.Registry, wd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprint(os.Stdout, usageText)
}
