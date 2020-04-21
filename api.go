package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	gg "github.com/gookit/color"
	"github.com/gorilla/mux"
)

// model
type Name2 struct {
	Filename   string
	Resolution string
	Took       string
}

func downloadAnimu(w http.ResponseWriter, r *http.Request) {

	// check that it is a post request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid request method\n")
		gg.Red.Println("Invalid request")
		return
	}

	// get name etc. default resolution is 480p
	r.ParseForm()
	ep := r.Form.Get("ep")
	name := r.Form.Get("name")
	reso1 := r.Form.Get("reso")
	// searchForm := r.Form.Get("type")

	// check if custom resolution is given, otherwise default to 480p
	if reso1 == "" {
		reso1 = "480"
		gg.Blue.Println("No resolution given, defaulting to 480p")
	}

	gg.Blue.Println("Downloading: " + name + " episode: " + ep)

	// download and time it
	time1 := time.Now()
	cmdOutput, err := exec.Command("anime-cli", "-q "+name, "-e "+ep, "-r "+reso1).Output()
	realout := string(cmdOutput[:])
	realout = strings.TrimSuffix(realout, "\n")
	took1 := time.Since(time1)

	if err != nil {
		gg.Red.Println("No anime found")
		fmt.Fprintf(w, "No anime found\n")
		return
	}

	// print and send result
	gg.Green.Println("Downloading: " + name + " episode: " + ep + " finished. Took: " + took1.String())

	// convert to json
	name2 := Name2{}
	name2.Filename = realout
	name2.Took = took1.String()
	name2.Resolution = reso1
	realoutJSON, err := json.Marshal(name2)

	if err != nil {
		gg.Red.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(realoutJSON)
	gg.Blue.Println("Waiting for new request...")
}

// search anime
func searchAnimu(w http.ResponseWriter, r *http.Request) {

	// check that it is a post request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid request method\n")
		gg.Red.Println("Invalid request")
		return
	}

	// get name etc. default resolution is 480p
	r.ParseForm()
	ep := r.Form.Get("ep")
	name := r.Form.Get("name")
	reso1 := r.Form.Get("reso")

	// check if custom resolution is given, otherwise default to 480p
	if reso1 == "" {
		reso1 = "480p"
		gg.Blue.Println("No resolution given, defaulting to 480p")
	}

	// vars for request
	api := "https://api.nibl.co.uk/nibl/"
	gg.Blue.Println("Encoding url...")
	anime := "/search?query=" + url.QueryEscape(name) + "%20" + reso1 + "&episodeNumber=" + ep

	// print and send request
	gg.Blue.Println("Search request received. Searching for:", name, ep, reso1)
	gg.Blue.Println("Sending request...")
	rep, e := http.Get(api + anime)

	if e != nil {
		gg.Red.Println("Error:", e.Error())
		fmt.Fprintf(w, "Error:", e.Error())
		return
	}

	defer rep.Body.Close()

	// unmarshal status & convert rep.Body
	var result map[string]interface{}
	json.NewDecoder(rep.Body).Decode(&result)
	gg.Blue.Println("Unmarshalling...")
	resp, e := ioutil.ReadAll(rep.Body)

	if e != nil {
		gg.Red.Println("Error:", e.Error())
		fmt.Fprintf(w, "Error:", e.Error())
		return
	}

	// print and send response
	gg.Green.Println("Status:", result["status"], "\nSearch result sent.")
	fmt.Fprintf(w, string(resp))
	return
}

// cleaning function
func cleaning() {

	// path and other vars
	home, err := os.UserHomeDir()
	dir := home + "/AnimeDownloads/"
	dirRead, _ := os.Open(dir)
	dirFiles, _ := dirRead.Readdir(0)

	// error handling
	if err != nil {
		gg.Red.Println("deleting downloaded files failed, error: ", err)
	}

	// loop it
	gg.Blue.Println("\ndeleting downloaded files...")
	for index := range dirFiles {
		fileHere := dirFiles[index]

		// get names of files the path.
		nameHere := fileHere.Name()
		fullPath := dir + nameHere

		// delete files
		os.Remove(fullPath)
	}

	gg.Green.Println("files deleted, shutting down...")
}

func main() {

	// this makes me look cool
	port := "1337"
	print("\033[H\033[2J")
	gg.Blue.Println("listening on port:", port)

	// clean files after ctrl + c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleaning()
		os.Exit(1)
	}()

	// start listening
	r := mux.NewRouter().StrictSlash(true)

	// download path
	r.HandleFunc("/download", downloadAnimu).Methods("POST")

	// search path
	r.HandleFunc("/search", searchAnimu).Methods("POST")
	err := http.ListenAndServe(":"+port, r)

	if err != nil {
		gg.Red.Println(err)
		return
	}
}
