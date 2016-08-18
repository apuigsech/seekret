package lib

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	SourceTypeDir = &SourceDir{}
)

type SourceDir struct{}

type SourceDirLoadOptions struct {
	Hidden    bool
	Recursive bool
}

func prepareDirLoadOptions(o LoadOptions) SourceDirLoadOptions {
	opt := SourceDirLoadOptions{
		Hidden:    false,
		Recursive: true,
	}

	if hidden, ok := o["hidden"].(bool); ok {
		opt.Hidden = hidden
	}
	if recursive, ok := o["recursive"].(bool); ok {
		opt.Hidden = recursive
	}

	return opt
}

func (s *SourceDir) LoadObjects(source string, o LoadOptions) ([]Object, error) {
	var objectList []Object

	opt := prepareDirLoadOptions(o)

	firstPath := true

	filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {

		if fi.IsDir() {
			if strings.HasPrefix(filepath.Base(path), ".") && !opt.Hidden {
				return filepath.SkipDir
			}

			if !firstPath && !opt.Recursive {
				return filepath.SkipDir
			}
			firstPath = false
		} else {
			if !strings.HasPrefix(filepath.Base(path), ".") || (strings.HasPrefix(filepath.Base(path), ".")  && opt.Hidden) {
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				
				content, err := ioutil.ReadAll(f)
				if err != nil {
					return err
				}

				o := NewObject(path, content)
		
				objectList = append(objectList, *o)
			}
		}

		return nil
	})

	return objectList, nil
}
