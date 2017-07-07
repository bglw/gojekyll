package site

import (
	"time"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/templates"
	"github.com/osteele/liquid/evaluator"
)

// ToLiquid returns the site variable for template evaluation.
func (s *Site) ToLiquid() interface{} {
	// double-checked lock is okay here, since it's okay if this gets
	// written twice
	if len(s.drop) == 0 {
		s.Lock()
		defer s.Unlock()
		if len(s.drop) == 0 {
			s.initializeDrop()
		}
	}
	return s.drop
}

// MarshalYAML is part of the yaml.Marshaler interface
// The variables subcommand uses this.
func (s *Site) MarshalYAML() (interface{}, error) {
	return s.ToLiquid(), nil
}

func (s *Site) initializeDrop() {
	vars := templates.MergeVariableMaps(s.config.Variables, map[string]interface{}{
		"data":         s.data,
		"documents":    s.docs,
		"html_pages":   s.htmlPages(),
		"pages":        s.Pages(),
		"static_files": s.StaticFiles(),
		// TODO read time from _config, if it's available
		"time": time.Now(),
		// TODO static_files, html_files, tags.TAG
	})
	collections := []interface{}{}
	for _, c := range s.Collections {
		vars[c.Name] = c.Pages()
		collections = append(collections, c.ToLiquid())
	}
	evaluator.SortByProperty(collections, "label", true)
	vars["collections"] = collections
	s.drop = vars
	s.setPostVariables()
}

func (s *Site) setPageContent() error {
	for _, c := range s.Collections {
		if err := c.SetPageContent(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *Site) htmlPages() (out []pages.Page) {
	for _, p := range s.Pages() {
		if p.OutputExt() == ".html" {
			out = append(out, p)
		}
	}
	return
}
