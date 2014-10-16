package main

import (
	"fmt"
	"io"
	"net/http"
	"io/ioutil"
	"log"
	"bytes"
	"github.com/gorilla/mux"
	"html/template"
	"os"
	"strconv"
)

var uploadTemplate = template.Must(template.ParseFiles("tmpl/upload.html"))
var filesTemplate = template.Must(template.ParseFiles("tmpl/files.html"))
 
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	uploadTemplate.ExecuteTemplate(w, tmpl+".html", data)
}

/* 
	Serve a file download
*/
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		filename := vars["filename"]

		path := "./docs/"
		fullPath := path + filename

		fmt.Println("Serving " + fullPath)

		content, err := ioutil.ReadFile(fullPath)

		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(content)
		nopCloser := ioutil.NopCloser(reader)

		w.Header().Set("Content-Disposition", "attachment; filename=" + filename)
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))

		io.Copy(w, nopCloser)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

/* 
	Accept a file upload
*/
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	hasUpload := false

	switch r.Method {
	case "GET":
		display(w, "upload", nil)
 
	case "POST":
		reader, err := r.MultipartReader()
 
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
 
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
 
			if part.FileName() == "" {
				continue
			}

			filePath := "./docs/" + part.FileName()
			_, err = os.Stat(filePath)
			if err == nil {
			    display(w, "upload", "Upload failed, file already exists.")
			    return
			}

			dst, err := os.Create(filePath)
			hasUpload = true
			defer dst.Close()
 
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		
		if hasUpload {
			display(w, "upload", "Upload successful.")
		} else {
			display(w, "upload", "No file sent.")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

/* 
	Show an index of all the files
*/
func filesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		files, err := ioutil.ReadDir("./docs/")

		if err != nil {
			log.Fatal(err)
		}

		content := ""
		for _, file := range files {
			content += fmt.Sprintf("<a href='/file/%s'>%s</a><br />", file.Name(), file.Name())
		}
		content += "<br />"

		filesTemplate.ExecuteTemplate(w, "files.html", template.HTML(content))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	port := ":9191"

	r := mux.NewRouter()

	r.HandleFunc("/", filesHandler)
	r.HandleFunc("/index", filesHandler)
	r.HandleFunc("/index/", filesHandler)
	r.HandleFunc("/files", filesHandler)
	r.HandleFunc("/files/", filesHandler)

	r.HandleFunc("/file/{filename}", downloadHandler)
	r.HandleFunc("/file/{filename}/", downloadHandler)

	r.HandleFunc("/upload", uploadHandler)
	r.HandleFunc("/upload/", uploadHandler)

	http.Handle("/", r)

	fmt.Println("Server listening on " + port)
	err := http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Println(err)
	}
}