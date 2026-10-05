package main

import (
    "html/template"
    "os"
    "net/http"
    "log"
    "regexp"
)

type Page struct {
    Title string
    Body  []byte
}
const baseTemplatesPath = "templates"

func buildTemplatePath(path string) string {
    return baseTemplatesPath + "/" + path
}

const baseDataPath = "data"

func buildDataPath(path string) string {
    return baseDataPath + "/" + path
}

var templates = template.Must(template.ParseFiles(buildTemplatePath("edit.html"), buildTemplatePath("view.html")))

var validPath = regexp.MustCompile("^/(edit|save|view|delete)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
    filename := buildDataPath(p.Title + ".txt")
    return os.WriteFile(filename, p.Body, 0600)
}

func (p *Page) delete() error {
    filename := buildDataPath(p.Title + ".txt")
    return os.Remove(filename)
}

func loadPage(title string) (*Page, error) {
    filename := buildDataPath(title + ".txt")
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    fullFileName := tmpl + ".html"
    err := templates.ExecuteTemplate(w, fullFileName, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    pToSave := &Page{Body: []byte(body), Title: title}
    err := pToSave.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, title string) {
    pToDelete := &Page{Title: title}
    err := pToDelete.delete()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write([]byte("File deleted"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/view/FrontPage/", http.StatusFound)
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    http.HandleFunc("/delete/", makeHandler(deleteHandler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}