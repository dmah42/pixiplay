package pixiplay

import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"

	"appengine"
	"appengine/datastore"
)

type Script struct {
	Name	string
	Author	string
	Content	string
}

var templateFuncs = map[string]interface{} {
	"eq": func(a, b interface{}) bool {
		switch a := a.(type) {
		case string, int, int64, byte, float32, float64:
			return a == b
		}

		return reflect.DeepEqual(a, b)
	},
}

var rootTemplate = template.Must(template.ParseFiles("templates/index.html"))
var scriptTemplate = template.Must(template.ParseFiles("templates/script.html")).Funcs(templateFuncs)
var authorTemplate = template.Must(template.ParseFiles("templates/author.html"))
var submitTemplate = template.Must(template.ParseFiles("templates/submit.html"))

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/author/", author)
	http.HandleFunc("/script/", script)
	http.HandleFunc("/submit", submit)
}

func allScripts(c appengine.Context) (scripts []Script, err error) {
	q := datastore.NewQuery("Script").Order("Name")
	scripts = make([]Script, 0)
	_, err = q.GetAll(c, &scripts)
	return
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	scripts, err := allScripts(c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := rootTemplate.Execute(w, scripts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func script(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get script name from request path
	name := r.URL.Path[len("/script/"):]

	c.Infof(fmt.Sprintf("name: %q", name))

	key := datastore.NewKey(c, "Script", name, 0, nil)
	var script Script
	err := datastore.Get(c, key, &script)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := scriptTemplate.Execute(w, script); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func author(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get author name from request path
	name := r.URL.Path[len("/author/"):]
	c.Infof(fmt.Sprintf("name: %q", name))

	q := datastore.NewQuery("Script").Filter("Author =", author).Order("Name")
	scripts := make([]Script, 0)
	_, err := q.GetAll(c, &scripts)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := authorTemplate.Execute(w, scripts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func submit(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if (r.Method == "GET") {
		scripts, err := allScripts(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := submitTemplate.Execute(w, scripts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// TODO: 'test' link
	} else if (r.Method == "POST") {
		// TODO: authenticate and set 'author' name
		script := &Script {
			Name: r.FormValue("name"),
			Author: r.FormValue("author"),
			Content: r.FormValue("content"),
		}

		key := datastore.NewKey(c, "Script", script.Name, 0, nil)
		_, err := datastore.Put(c, key, script)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/game/" + script.Name, http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
