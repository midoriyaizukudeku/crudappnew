package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"dream.website/internal/model"
	"dream.website/ui"
)

type Dynamicdata struct {
	CurrentYear     int
	SingleSnippet   *model.Snippet
	Snippets        []*model.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func HumanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 January 2006 at 15:04")
}

var functions = template.FuncMap{
	"HumanDate": HumanDate,
}

func NewTemplatecache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		pattern := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, pattern...)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %s: %w", name, err)
		}

		cache[name] = ts
	}
	return cache, nil
}
