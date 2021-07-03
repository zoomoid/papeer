package main

import (
	"errors"
	"io/fs"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/zoomoid/papeer/pkg/runner"
	"github.com/zoomoid/papeer/pkg/types"
	"go.uber.org/zap"
)

var watchArg bool
var fileArg string
var directoryArg string
var landscapeArg bool
var scaleArg float32
var pagesArg string
var gridArg bool
var gridOpacityArg float32
var gridColorArg string
var gridScaleArg float32
var formatArg string
var deltaArg []string

var supportedFormats = []string{
	"a2paper",
	"a3paper",
	"a4paper",
}

func main() int {
	flag.BoolVarP(&watchArg, "watch", "w", false, "Start in watch mode with directory")
	flag.StringVarP(&fileArg, "file", "f", "", "Explicit filename given")
	flag.StringVarP(&directoryArg, "directory", "d", "", "Explicit filename given")
	flag.BoolVar(&landscapeArg, "landscape", true, "Wrap the pdf in a landscape format")
	flag.Float32Var(&scaleArg, "scale", 0.9, "Paper scale on the new page")
	flag.StringVar(&pagesArg, "pages", "-", "pdfpages-compliant 'pages' parameter. '-' means all pages of the given pdf")
	flag.BoolVar(&gridArg, "grid", true, "Add a dotted grid with 5mm spacing to the background")
	flag.Float32Var(&gridOpacityArg, "grid-opacity", 0.1, "Specify the opacity of the grid")
	flag.StringVar(&gridColorArg, "grid-color", "black!30", "Specify the color of the grid dots")
	flag.Float32Var(&gridScaleArg, "grid-scale", 1.0, "Specify the scale of the grid")
	flag.StringVar(&formatArg, "format", "a3paper", "Paper format to embed the pdf in")
	flag.StringSliceVar(&deltaArg, "delta", []string{"10cm", "10cm"}, "Specify the page margins in (x,y) order")

	flag.Parse()

	filehandle := flag.Arg(0)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if !fs.ValidPath(filehandle) {
		logger.Fatal("Not a valid path", zap.String("path", filehandle))
		return 1
	}

	f, err := os.Stat(filehandle)
	if err != nil {
		logger.Fatal("Failed to stat file", zap.String("file", filehandle))
		return 1
	}

	config, err := makeWrapperOptions(filehandle)

	if err != nil {
		logger.Fatal(err.Error())
		return 1
	}

	if watchArg {
		if !f.IsDir() {
			logger.Fatal("Given path is not a directory", zap.String("path", filehandle))
			return 1
		}
		_, err := watchDirectory(filehandle)
		if err != nil {
			logger.Fatal(err.Error())
			return 1
		}
	} else {
		if !f.Mode().IsRegular() {
			logger.Fatal("Given file is not regular", zap.String("file", filehandle))
			return 1
		}
		output, err := oneOffBuild(filehandle, config)
		if err != nil {
			logger.Fatal(err.Error())
		}
		logger.Info(output)
	}
	return 0
}

// makeWrapperOptions creates a WrapperOptions struct from the given arguments
func makeWrapperOptions(arg string) (*types.WrapperOptions, error) {
	var options types.WrapperOptions

	if watchArg && len(fileArg) > 0 {
		return nil, errors.New("cannot specify --watch and --file")
	}

	file := fileArg
	if file != "" && len(arg) > 0 {
		return nil, errors.New("ambiguous files, specify either --file or the argument")
	}
	if file == "" {
		file = arg
	}

	options.Filename = file

	delta := [2]string{}
	if len(deltaArg) == 1 {
		delta = [2]string{delta[0], delta[0]}
	}
	if len(delta) != 2 {
		return nil, errors.New("illegal dimensions of delta argument")
	}

	options.Delta.X = delta[0]
	options.Delta.Y = delta[1]

	options.Format = validateFormat(formatArg)

	options.Grid.Enabled = gridArg
	options.Grid.Color = gridColorArg
	options.Grid.Scale = gridScaleArg
	options.Grid.Opacity = gridOpacityArg

	options.Landscape = landscapeArg
	options.Pages = pagesArg
	options.Scale = scaleArg

	return &options, nil
}

// watchDirectory monitors a directory for occuring PDFs and spawns a runner
// TODO: Implement me!
func watchDirectory(dir string) (int, error) {
	return 0, nil
}

// oneOffBuild runs the runner process once and returns the filename of the output PDF
func oneOffBuild(file string, config *types.WrapperOptions) (string, error) {
	return runner.Run(config)
}

// validateFormat checks if format is a supported one, otherwise defaults to a3paper
func validateFormat(format string) string {
	for _, element := range supportedFormats {
		if element == format {
			return format
		}
	}
	return "a3paper"
}
