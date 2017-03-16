# {{.Name}}

{{template "badge/travis" .}}{{template "badge/appveyor" .}}{{template "badge/goreport" .}}{{template "badge/godoc" .}}

{{pkgdoc}}

Compatible __windows__ / __linux__

# Install

#### Go
{{template "go/install" .}}

# Usage

{{cli "disksinfo"}}

# API example

{{file "main.go"}}

# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```
