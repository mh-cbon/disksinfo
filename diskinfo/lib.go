package diskinfo

import "runtime"

// PropertiesLoader can load the list of partitions and their properties.
type PropertiesLoader interface {
	Load() ([]*Properties, error)
}

// NewMultiOsLoader prepares a multi os loader for the current runtime operating system.
func NewMultiOsLoader() PropertiesLoader {
	var loader PropertiesLoader
	loader = &LinuxLoader{}
	if runtime.GOOS == "windows" {
		loader = &WindowsLoader{}
	}
	return loader
}

// Properties provides information about a partition on the system.
type Properties struct {
	Label       string
	IsRemovable bool
	Size        string
	SpaceLeft   string
	Path        string
	MountPath   string
}

// NewProperties is a constructor.
func NewProperties() *Properties {
	return &Properties{}
}

// PropertiesList ...
type PropertiesList []*Properties

// Merge ...
func (l PropertiesList) Merge(some PropertiesList, what ...string) []*Properties {
	for i, d := range l {
		s := some.FindByPath(d.Path)
		if s != nil {
			for _, w := range what {
				switch w {
				case "Label":
					if s.Label != "" {
						d.Label = s.Label
					}
				case "Size":
					if s.Size != "" {
						d.Size = s.Size
					}
				case "SpaceLeft":
					if s.SpaceLeft != "" {
						d.SpaceLeft = s.SpaceLeft
					}
				case "MountPath":
					if s.MountPath != "" {
						d.MountPath = s.MountPath
					}
				case "Path":
					if s.Path != "" {
						d.Path = s.Path
					}
				case "IsRemovable":
					d.IsRemovable = s.IsRemovable
				}
			}
			l[i] = d
		}
	}
	return l
}

// Append ...
func (l PropertiesList) Append(some PropertiesList) []*Properties {
	for _, d := range some {
		s := l.FindByPath(d.Path)
		if s == nil {
			l = append(l, d)
		}
	}
	return l
}

// FindByPath ...
func (l PropertiesList) FindByPath(path string) *Properties {
	for _, d := range l {
		if d.Path == path {
			return d
		}
	}
	return nil
}
