package main

import (
	"flag"

	"go.uber.org/zap"

	"github.com/zoomoid/papeer/pkg/runner"
	types "github.com/zoomoid/papeer/pkg/types"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	landscapeArg := flag.Bool("landscape", true, "Wrap the pdf in a landscape format")
	gridArg := flag.Bool("grid", true, "Add a dotted grid with 5mm spacing to the background")
	formatArg := flag.String("format", "a3paper", "Paper format to embed the pdf in")
	flag.Parse()

	fileArg := flag.Arg(0)

	logger.Info("Starting papeer",
		zap.String("file", fileArg),
		zap.Bool("landscape", *landscapeArg),
		zap.Bool("grid", *gridArg),
		zap.String("format", *formatArg),
	)

	options := types.WrapperOptions{
		Format:    *formatArg,
		Landscape: *landscapeArg,
		Grid: types.GridOptions{
			Enabled: *gridArg,
			Opacity: 0.1,
			Color:   "black!30",
			Scale:   1,
		},
		Scale: 0.9,
		Pages: "-",
		Delta: types.DeltaOptions{
			X: "10",
			Y: "10",
		},
		Filename: fileArg,
	}

	resPath, err := runner.Run(&options)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("Wrote pdf to filesystem", zap.String("path", resPath))
}
