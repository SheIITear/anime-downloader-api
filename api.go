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
	"math/rand"
	"path/filepath"

	gg "github.com/gookit/color"
	"github.com/gorilla/mux"
)

// model
type Name2 struct {
	Filename   		string
	Resolution 		string
	TookDownload    string
	TookCompress	string
	CompressError 	bool
	Compressed		bool
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
	compress := r.Form.Get("compress")
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
	took2 := took1.String()

	if err != nil {
		gg.Red.Println("No anime found")
		fmt.Fprintf(w, "No anime found\n")
		return
	}

	// print and send result
	gg.Green.Println("Downloading: " + name + " episode: " + ep + " finished. Took: " + took1.String())

	if compress == "true" {
		var extension = filepath.Ext(realout)
		var name = realout[0:len(realout)-len(extension)]
		compressAnime(name, extension, took2, w, r)
	} else {

	// convert to json
	name2 := Name2{}
	name2.Filename = realout
	name2.TookDownload = took2
	name2.Resolution = reso1
	name2.Compressed = false
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
}

// compress
func compressAnime(name string, extension string, took2 string, w http.ResponseWriter, r *http.Request) {

	home, err := os.UserHomeDir()

	if err != nil{
		fmt.Println("Getting dir failed")
	}

	dir := home + "/AnimeDownloads/"

	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	length := 8
	var a strings.Builder
	var b strings.Builder


	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
		a.WriteRune(chars[rand.Intn(len(chars))])
	}

	newname := b.String() + extension
	name3 := a.String() + extension

	name = strings.Replace(name, " ", "_", -1)
	err = os.Rename(dir + name + extension, dir + name3)

	if err != nil {
		gg.Red.Println("Filename error:", err)
	}

	name2 := Name2{}

	command := "ffmpeg -i " + dir + name3 + " -vcodec libx265 -crf 28 -preset fast -c:a copy " + dir + newname
	cmdString := strings.TrimSuffix(command, "\n")
	cmdString2 := strings.Fields(cmdString)
	gg.Blue.Println("Received file: " + name3 + "\nStarting compression, please stand by...")

	// compress and time it
	time1 := time.Now()
	cmdOutput, err := exec.Command(cmdString2[0], cmdString2[1:]...).Output()
	realout := string(cmdOutput[:])
	realout = strings.TrimSuffix(realout, "\n")
	took1 := time.Since(time1)

	if err != nil {
		name2.CompressError = true
		gg.Red.Println("Compress failed:", err)
	} else {
		name2.CompressError = false
	}

	name2.Filename = newname + extension
	name2.TookDownload = took2
	name2.TookCompress = took1.String()
	name2.Compressed = true
	realoutJSON, err := json.Marshal(name2)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(realoutJSON)

	gg.Green.Println("Anime compressed and response sent.\nWaiting for new request...")
	return
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
	gg.Blue.Println("Search request received. Searching for: ", name, ep, reso1)
	gg.Blue.Println("Sending request...")
	rep, e := http.Get(api + anime)

	if e != nil {
		gg.Red.Println("Error:", e.Error())
		fmt.Fprintf(w, "Error:", e.Error())
		return
	}

	defer rep.Body.Close()

	// response
	resp, e := ioutil.ReadAll(rep.Body)

	if e != nil {
		gg.Red.Println("Error:", e.Error())
		fmt.Fprintf(w, "Error:", e.Error())
		return
	}

	// send response
	fmt.Fprintf(w, string(resp))
	gg.Green.Println("Search result sent.")
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

	// delete files after ctrl + c
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
