package main

import (
	"html/template"
	"net/url"
	"path/filepath"
	"time"

	"jackson.software/snippetbox/pkg/models"
)

type templateData struct {
	CurrentYear int
	FormData    url.Values
	FormErrors  map[string]string
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// Transforms the given date time to a better human readable presentation.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
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

		ts, err := getTemplateForPage(name, page, dir)
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
func getTemplateForPage(name, page, dir string) (*template.Template, error) {
	ts, err := template.New(name).Funcs(functions).ParseFiles(page)
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
