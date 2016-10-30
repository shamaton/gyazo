package main

/*
  gyazo.go

  Application for private gyazo service.
*/
import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

var HOST string
var BIND string

const HOST_DEFAULT = "localhost"
const PORT_DEFAULT = 8080

/*
  initialize
*/
func init() {
	// random seed for filename.
	rand.Seed(time.Now().UnixNano())

	var port int
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.StringVar(&HOST, "h", HOST_DEFAULT, "hostname")
	f.StringVar(&HOST, "host", HOST_DEFAULT, "hostname")
	f.IntVar(&port, "p", PORT_DEFAULT, "port")
	f.IntVar(&port, "port", PORT_DEFAULT, "port")
	f.Parse(os.Args[1:])

	BIND = ":" + strconv.Itoa(port)
}

/*
   entry point.
*/
func main() {
	// routing
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/images/", imagesHandler)
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("start listen... ", BIND)

	// listen
	err := http.ListenAndServe(BIND, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

/*
  upload image from gyazo client.

  arguments
    w : response writer
    r : request
*/
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// get project directory
	dir, err := getProjectDir()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// make directory if it don't exist.
	imageDir := path.Join(dir, "images")
	if err := os.Mkdir(imageDir, 0755); err != nil && !os.IsExist(err) {
		fmt.Fprintln(w, err)
		return
	}

	// get image data from request/
	imgData, _, err := r.FormFile("imagedata")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer imgData.Close()

	// get filename from md5
	key := strconv.FormatInt(rand.Int63(), 10)
	timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)

	h := md5.New()
	io.WriteString(h, key+timeStr)
	baseName := fmt.Sprintf("%x", h.Sum(nil))

	// create emplty file.
	imageFile := path.Join(imageDir, baseName)
	file, err := os.Create(imageFile)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()

	// copy to file from request data.
	_, err = io.Copy(file, imgData)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// redirect created address.
	fmt.Fprintf(w, "http://%s/images/%s", HOST, baseName)
}

/*
  show clipped image from redirect.

  arguments
    w : response writer
    r : request
*/
func imagesHandler(w http.ResponseWriter, r *http.Request) {
	// get project directory
	dir, err := getProjectDir()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	// path to image files
	imagefile := path.Join(dir, r.URL.Path)
	http.ServeFile(w, r, imagefile)
}

/*
  check response. (ping test)

  arguments
    w : response writer
    r : request
*/
func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

/*
  get project directory
*/
func getProjectDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(dir, ".."), nil
}
