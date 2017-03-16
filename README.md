# disksinfo

[![travis Status](https://travis-ci.org/mh-cbon/disksinfo.svg?branch=master)](https://travis-ci.org/mh-cbon/disksinfo)[![appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/disksinfo?branch=master&svg=true)](https://ci.appveyor.com/project/mh-cbon/disksinfo)
[![Go Report Card](https://goreportcard.com/badge/github.com/mh-cbon/disksinfo)](https://goreportcard.com/report/github.com/mh-cbon/disksinfo)

[![GoDoc](https://godoc.org/github.com/mh-cbon/disksinfo?status.svg)](http://godoc.org/github.com/mh-cbon/disksinfo)


Package disksinfo provides the list of partitions


Compatible __windows__ / __linux__

# Install

#### Go

```sh
go get github.com/mh-cbon/disksinfo
```


# Usage


###### $ disksinfo 
```sh
[
    {
        "Label": "",
        "IsRemovable": false,
        "Size": "1,9G",
        "SpaceLeft": "1,9G",
        "Path": "devtmpfs",
        "MountPath": "/dev"
    },
    {
        "Label": "",
        "IsRemovable": false,
        "Size": "1,9G",
        "SpaceLeft": "1,9G",
        "Path": "tmpfs",
        "MountPath": "/dev/shm"
    },
    {
        "Label": "",
        "IsRemovable": false,
        "Size": "32G",
        "SpaceLeft": "3,4G",
        "Path": "/dev/mapper/fedora-root",
        "MountPath": "/"
    },
    {
        "Label": "",
        "IsRemovable": false,
        "Size": "126G",
        "SpaceLeft": "15G",
        "Path": "/dev/sda5",
        "MountPath": "/home"
    },
    {
        "Label": "whatever",
        "IsRemovable": true,
        "Size": "932G",
        "SpaceLeft": "515G",
        "Path": "/dev/sdb1",
        "MountPath": "/run/media/mh-cbon/whatever"
    },
    {
        "Label": "Recovery",
        "IsRemovable": false,
        "Size": "",
        "SpaceLeft": "",
        "Path": "/dev/sda1",
        "MountPath": ""
    },
    {
        "Label": "stockage",
        "IsRemovable": false,
        "Size": "",
        "SpaceLeft": "",
        "Path": "/dev/sda6",
        "MountPath": ""
    },
    {
        "Label": "System Reserved",
        "IsRemovable": false,
        "Size": "",
        "SpaceLeft": "",
        "Path": "/dev/sda2",
        "MountPath": ""
    }
]
```

# API example


###### > main.go
```go
// Package disksinfo provides the list of partitions
package main

import (
	"encoding/json"
	"os"

	"github.com/mh-cbon/disksinfo/diskinfo"
)

func main() {
	loader := diskinfo.NewMultiOsLoader()
	p, err := loader.Load()
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(p); err != nil {
		panic(err)
	}
}
```

# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```
