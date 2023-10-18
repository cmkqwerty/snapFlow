package views

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cmkqwerty/snapFlow/context"
	"github.com/cmkqwerty/snapFlow/models"
	"github.com/gorilla/csrf"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
)

type public interface {
	Public() string
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	htmlTpl := template.New(pattern[0])
	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField is not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser is not implemented")
			},
			"errors": func() []string {
				return nil
			},
		},
	)

	htmlTpl, err := htmlTpl.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{
		htmlTpl: htmlTpl,
	}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	htmlTpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}

	errMssgs := errMessages(errs...)
	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMssgs
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var buf bytes.Buffer
	err = htmlTpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}

	return t
}

func errMessages(errs ...error) []string {
	var errMessages []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			errMessages = append(errMessages, pubErr.Public())
		} else {
			fmt.Println(err)
			errMessages = append(errMessages, "Something went wrong. Please try again later.")
		}
	}
	return errMessages
}
