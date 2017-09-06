package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/kyokomi/cloudinary"
	"github.com/urfave/negroni"
)

func main() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/upload", uploadHandler).Methods("POST")
	router.HandleFunc("/favicon.ico", faviconHandler)

	mux := http.NewServeMux()
	mux.Handle("/", router)

	static := http.StripPrefix("/public/", http.FileServer(http.Dir("public")))
	router.PathPrefix("/public").Handler(static)

	n := negroni.Classic()
	n.UseHandler(mux)

	port := getEnv("PORT", ":3000")
	http.ListenAndServe(port, n)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/index.html", "templates/layout.html")
	err = tpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Printf("Error parsing templates %s\n", err.Error())
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("filename")
	fmt.Printf("Filename %s\n", fileName)
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error reading form image %s\n", err.Error())
		return
	}
	defer file.Close()

	path := "./public/temp" + fileName

	out, err := os.Create(path)
	if err != nil {
		log.Printf("Error creating file in public/tmp %s", err.Error())
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Printf("Error writing file to public/tmp %s", err.Error())
		return
	}

	url, err := cloudinaryUpload(path, fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}

	openBrowser(url)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/favicon.ico")
}

func cloudinaryUpload(src string, fileName string) (string, error) {
	ctx := context.Background()

	key := getEnv("CLOUDINARY_API_KEY", "")
	secret := getEnv("CLOUDINARY_API_SECRET", "")
	cloud := getEnv("CLOUDINARY_CLOUD_NAME", "")

	con := fmt.Sprintf("cloudinary://%s:%s@%s", key, secret, cloud)
	ctx = cloudinary.NewContext(ctx, con)

	data, _ := ioutil.ReadFile(src)

	if err := cloudinary.UploadStaticImage(ctx, fileName, bytes.NewBuffer(data)); err != nil {
		log.Println("Error uploading image to cloudinary")
		return "", err
	}

	url := cloudinary.ResourceURL(ctx, fileName)

	return url, nil
}

func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
