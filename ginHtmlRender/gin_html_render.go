package ginHtmlRender

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"github.com/gin-gonic/gin/render"
)

const (
	TemplatesDir = "templates/"
	Layout = "layout"
	Ext = ".html"
	Debug = false
)

// Render implements gin's HTMLRender and provides some sugar on top of it
type Render struct {
	Templates    map[string]*template.Template
	Files        map[string][]string
	TemplatesDir string
	Layout       string
	Ext          string
	Debug        bool
}

// Add assigns the name to the template
func (r *Render) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	r.Templates[name] = tmpl
}

// AddFromFiles parses the files and returns the result
func (r *Render) AddFromFiles(name string, files ...string) *template.Template {
	tmpl := template.Must(template.ParseFiles(files...))
	if r.Debug {
		r.Files[name] = files
	}
	r.Add(name, tmpl)
	return tmpl
}

// Instance implements gin's HTML render interface
func (r *Render) Instance(name string, data interface{}) render.Render {
	var tpl *template.Template

	// Check if gin is running in debug mode and load the templates accordingly
	if r.Debug {
		tpl = r.loadTemplate(name)
	} else {
		tpl = r.Templates[name]
	}

	return render.HTML{
		Template: tpl,
		Data:     data,
	}
}

// loadTemplate parses the specified template and returns it
func (r *Render) loadTemplate(name string) *template.Template {
	tpl, err := template.ParseFiles(r.Files[name]...)
	if err != nil {
		panic(name + " template name mismatch")
	}
	return template.Must(tpl, err)
}

// New returns a fresh instance of Render
func New() Render {
	return Render{
		Templates:    make(map[string]*template.Template),
		Files:        make(map[string][]string),
		TemplatesDir: TemplatesDir,
		Layout:       Layout,
		Ext:          Ext,
		Debug:        Debug,
	}
}

// Create goes through the `TemplatesDir` creating the template structure
// for rendering. Returns the Render instance.
func (r *Render) Create() *Render {
	r.Validate()

	layout := r.TemplatesDir + r.Layout + r.Ext

	// root dir
	tplRoot, err := filepath.Glob(r.TemplatesDir + "*" + r.Ext)
	if err != nil {
		panic(err.Error())
	}

	// sub dirs
	tplSub, err := filepath.Glob(r.TemplatesDir + "**/*" + r.Ext)
	if err != nil {
		panic(err.Error())
	}

	for _, tpl := range append(tplRoot, tplSub...) {

		// This check is to prevent `panic: template: redefinition of template "layout"`
		name := r.getTemplateName(tpl)
		if name == r.Layout {
			continue
		}

		r.AddFromFiles(name, layout, tpl)
	}

	return r
}

// Validate checks if the directory and the layout files exist as expected
// and configured
func (r *Render) Validate() {
	// add trailing slash if the user has forgotten..
	if !strings.HasSuffix(r.TemplatesDir, "/") {
		r.TemplatesDir = r.TemplatesDir + "/"
	}

	// check for templates dir
	if ok, _ := exists(r.TemplatesDir); !ok {
		panic(r.TemplatesDir + " directory for rendering templates does not exist.\n Configure this by setting htmlRender.TemplatesDir = \"your-tpl-dir/\"")
	}

	// check for layout file
	layoutFile := r.TemplatesDir + r.Layout + r.Ext
	if ok, _ := exists(layoutFile); !ok {
		panic(layoutFile + " layout file does not exist")
	}
}

// getTemplateName returns the name of the template
// For example, if the template path is `templates/articles/list.html`
// getTemplateName would return `articles/list`
func (r *Render) getTemplateName(tpl string) string {
	dir, file := filepath.Split(tpl)
	dir = strings.Replace(dir, r.TemplatesDir, "", 1)
	file = strings.TrimSuffix(file, r.Ext)
	return dir + file
}

// exists returns whether the given file or directory exists or not
// http://stackoverflow.com/a/10510783/232619
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}