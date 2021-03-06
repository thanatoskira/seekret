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

	if hidden, ok := o["hidden"].(bool); ok == true {
		opt.Hidden = hidden
	}
	if recursive, ok := o["recursive"].(bool); ok == true {
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
			if strings.HasPrefix(filepath.Base(path), ".") == true && opt.Hidden == false {
				return filepath.SkipDir
			}

			if firstPath == false && opt.Recursive == false {
				return filepath.SkipDir
			}
			firstPath = false
		} else {
			if strings.HasPrefix(filepath.Base(path), ".") == false || (strings.HasPrefix(filepath.Base(path), ".") == true && opt.Hidden == true) {
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				content, err := ioutil.ReadAll(f)
				if err != nil {
					return err
				}

				o := Object{
					Name:    path,
					Content: content,
				}

				objectList = append(objectList, o)
			}
		}

		return nil
	})

	return objectList, nil
}
