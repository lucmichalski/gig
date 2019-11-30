package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/shihanng/gi/internal/file"
	"github.com/shihanng/gi/internal/order"
	"gopkg.in/src-d/go-git.v4"
)

const sourceRepo = `https://github.com/toptal/gitignore.git`

func main() {
	path := filepath.Join(xdg.CacheHome(), `gi`)

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      sourceRepo,
		Depth:    1,
		Progress: os.Stdout,
	})
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		log.Fatal(err)
	}

	args := os.Args[1:]
	languages := make(map[string]bool, len(args))

	for _, arg := range args {
		languages[file.Canon(arg)] = true
	}

	files, err := ioutil.ReadDir(filepath.Join(path, `templates`))
	if err != nil {
		log.Fatal(err)
	}

	giFiles := []file.File{}

	for _, f := range files {
		filename := f.Name()
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)

		if languages[file.Canon(base)] {
			giFiles = append(giFiles, file.File{Name: base, Typ: ext})
		}
	}

	orders, err := order.ReadOrder(filepath.Join(path, `templates`, `order`))
	if err != nil {
		log.Fatal(err)
	}

	giFiles = file.Sort(giFiles, orders)

	if err := file.Compose(os.Stdout, filepath.Join(path, `templates`), giFiles...); err != nil {
		log.Fatal(err)
	}
}
