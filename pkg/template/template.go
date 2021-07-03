package template

import (
	"bytes"
	"text/template"

	"github.com/zoomoid/papeer/pkg/types"
)

var wrapperTemplate string = `\documentclass{article}
\usepackage[landscape,a3paper,margin=0pt]{geometry}
\usepackage{pdfpages}
\usepackage{background}
\usepackage{graphicx}
{{ if .Grid.Enabled -}}
\backgroundsetup{
scale={{.Grid.Scale}},
angle=0,
color={{.Grid.Color}},
opacity={{.Grid.Opacity}},
contents={\includegraphics{grid}}
}
{{ end -}}
\title{}
\author{}
\date{}
\begin{document}
\includepdf[pages={{.Pages}},delta={{.Delta.X}} {{.Delta.Y}},scale={{.Scale}},noautoscale=true]{{{.Filename}}}
\end{document}
`

// Template assembles the main tex file that's compiled by latexmk
func Template(options *types.WrapperOptions) (string, error) {
	tpl, err := template.New("wrapper-document").Parse(wrapperTemplate)
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	err = tpl.Execute(&output, *options)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
