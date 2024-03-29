package api

import (
	"errors"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/docgen"
)

func GenerateDocs(r chi.Router, docs string) {
	if docs == "markdown" {
		// Remove any existing `routes.md` file in the project directory.
		if err := os.Remove("routes.md"); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}

		// Create a new `routes.md` file.
		f, err := os.Create("routes.md")

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		text := docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{ // Here, `r` is the main router.
			ProjectPath: "github.com/lcrownover/hpcadmin-server",
			URLMap: map[string]string{
				"github.com/lcrownover/hpcadmin-server/vendor/github.com/go-chi/chi/v5/": "https://github.com/go-chi/chi/blob/master/",
			},
			ForceRelativeLinks: true, // Without this, you will not see any links to the source code in the Markdown file on your local machine.
			Intro:              "Welcome to the documentation for the HPCAdmin REST API.",
		})

		// Write the Markdown generated by `docgen` to the `routes.md` file.
		if _, err = f.Write([]byte(text)); err != nil {
			log.Fatal(err)
		}

		return
	}
}
