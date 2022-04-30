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

var PreviousPath int64 = 1

var PATHS = make(map[string]string)

func main() {
	Init()
	fmt.Println("Initializing /create handler")
	fmt.Println("Listening on 127.0.0.1:8082")
	http.ListenAndServe("127.0.0.1:8082", nil)
}

// Initialize the basic handlers
func Init() {
	http.HandleFunc(CreatePath, createHandler)
}

// createHandler which vould create new handlers for shortned url
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

// The function which redirects the request to the url
func Redirector(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Fprintf(resp,
			"Method %s is not supported. Only GET method is supported",
			req.Method)
	}
	fmt.Printf("Redirecting %s to %s\n", req.URL.Path, PATHS[req.URL.Path])
	http.Redirect(resp, req, PATHS[req.URL.Path], 301)
}

// Genereate numeric path based on now unix time and previous path
func genPath() string {
	seed := time.Now().Unix() * PreviousPath
	rand.Seed(seed)
	ipath := int(rand.Uint32())
	PreviousPath = int64(ipath)
	path := strconv.Itoa(ipath)
	return string(path)
}
