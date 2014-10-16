butler
======

Easily serve file downloads and file uploads in Go

Installation
============

To get up and running, you need to first install Google Go (golang.org)
Please ensure that GOPATH is set to C:\Go\ and that GOROOT is deleted.
First run "go get github.com/gorilla/mux"
Then run "go run butler.go"
You should see Server listening on :9191
Go to a browser and navigate to localhost:9191

You can now do file uploads and serve downloads. 
(Note: You could use http.FileServer(http.Dir("/usr/share/doc")) to serve a directory via the common library, but this will open files in the browser instead of forcing a download.