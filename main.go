package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const CreatePath string = "/create"

var PATHS = make(map[string]string)

func main() {
	Init()
	fmt.Println("Initializing /create handler")
	http.ListenAndServe("127.0.0.1:8082", nil)
	fmt.Println("Listening on 127.0.0.1:8082")
}

// Initialize the basic handlers
func Init() {
	http.HandleFunc(CreatePath, createHandler)
}

func createHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Fprintf(resp,
			"Method %s is not supported. Only POST method with url in body is supported\n",
			req.Method)
	} else {
		path := "/" + genPath()
		burl, _ := io.ReadAll(req.Body)
		PATHS[path] = string(burl)
		fmt.Fprintln(resp, path)
		http.HandleFunc(path, Redirector)
		fmt.Printf("The path %s was created for %s\n", path, PATHS[path])
	}
}

func Redirector(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Fprintf(resp,
			"Method %s is not supported. Only GET method is supported",
			req.Method)
	}
	fmt.Printf("Redirecting %s to %s\n", req.URL.Path, PATHS[req.URL.Path])
	http.Redirect(resp, req, PATHS[req.URL.Path], 301)
}

func genPath() string {
	rand.Seed(time.Now().Unix())
	return string(strconv.Itoa(int(rand.Uint32())))
}
