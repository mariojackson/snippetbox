package main

import (
	"html/template"
	"jackson.software/snippetbox/pkg/models"
	"path/filepath"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// NewTemplateCache creates a cache of templates indexed by their page name,
// coming from all *.page.tmpl files found in the given directory, by reading
// all template files into a map.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize map to act like a cache
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := getTemplateForPage(page, dir)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

// GetTemplateForPage returns a template for the given page eg. home.page.tmpl,
// by parsing all layout and partial tmpl files which are found in the given directory.
//
// GetTemplateForPage returns an error when a file could not be parsed.
func getTemplateForPage(page, dir string) (*template.Template, error) {
	ts, err := template.ParseFiles(page)
	if err != nil {
		return nil, err
	}

	ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
	if err != nil {
		return nil, err
	}

	ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
	if err != nil {
		return nil, err
	}

	return ts, nil
}
