package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	start(path.Join("./samples"))

}

func start(root string) error {
	filepath.WalkDir(root, walkFunc)

	return nil
}

func walkFunc(p string, d fs.DirEntry, err error) error {
	if d.IsDir() {
		return filepath.SkipDir
	}
	if d != nil && !d.IsDir() {
		rx, e := regexp.Compile(`.*_\d{3}\.txt$`)
		if e != nil {
			fmt.Println(e)
			return e
		}
		if rx.MatchString(d.Name()) {
			usIdx := strings.LastIndex(d.Name(), "_")
			title, number, extension := d.Name()[:usIdx], d.Name()[usIdx+1:usIdx+4], d.Name()[usIdx+4:]

			newName := fmt.Sprintf("%s-%s%s", number, title, extension)
			newPath := path.Join(path.Dir(p), newName)
			e := os.Rename(p, newPath)
			if e != nil {
				println(err)
				return e
			}
		}
	}
	return nil
}
