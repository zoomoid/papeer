package papeer

import (
	"bytes"
	"text/template"
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

func Template(options *WrapperOptions) (string, error) {
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
