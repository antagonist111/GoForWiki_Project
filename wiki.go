package main

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte // Body element is 'byte slice' because io libaries expect that type!
}

// GLOBAL VARS:
var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// Saves page to persistent storage
// Takes a pointer to Page (p) as a receiver. Has no paramters, and returns error type
// Returns error type as it is return type of WriteFile
// 0600 means: File has read-write permissions for current user only
func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Load pages
// Reads files contents (constructed from title) into new variable body
// Returns a pointer to the Page literal constructed.
// the '_' is there because io.ReadFile returns []byte and error
// We aren't handling error yet, so it's 'blank' or '_' (FIXED)
func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

// Handler: Takes http.ResponseWriter and http.Request as arguments
// ResponseWriter assembles the HTTP response. (Write to = sends to client)
// Request is data structure = client HTTP request. r.URL.Path is
// the path component of request URL. [1:] drops leading "/"

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// editHandler loads a page, or creates an empty one if it doesn't exist.
// It then displays an HTML form

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p) // Passes: ResponseWriter, templatename, template page fill
	//t.Execute(w, p)  Executes template, writes HTML to http.ResponseWriter (w)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("body") // type String

	p := &Page{Title: title, Body: []byte(body)} // Convert string to byte
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/"+"FrontPage", http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {

	// Replace instances of [] as page links
	buf := new(bytes.Buffer)
	err := templates.ExecuteTemplate(buf, tmpl+".html", p)
	re := regexp.MustCompile(`\[(.+?)\]`)
	output := re.ReplaceAllString(buf.String(), `<a href="/view/$1">$1</a>.`)
	w.Write([]byte(output))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
