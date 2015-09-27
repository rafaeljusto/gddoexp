// gddoexp is a command line tool crated to list eligible packages for
// archiving in GoDoc.org
package main

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/golang/gddo/database"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/rafaeljusto/gddoexp"
)

func main() {
	// add cache to avoid to repeated requests to Github
	gddoexp.HTTPClient = &http.Client{
		Transport: httpcache.NewTransport(
			diskcache.New(path.Join(os.Getenv("HOME"), ".gddoexp")),
		),
	}

	db, err := database.New()
	if err != nil {
		fmt.Println("error connecting to database:", err)
		return
	}

	pkgs, err := db.AllPackages()
	if err != nil {
		fmt.Println("error retrieving all packages:", err)
		return
	}

	for _, pkg := range pkgs {
		if archive, err := gddoexp.ShouldArchivePackage(pkg.Path, db); err != nil {
			fmt.Println(err)
		} else if archive {
			fmt.Printf("package “%s” should be archived\n", pkg.Path)
		}
	}
}
