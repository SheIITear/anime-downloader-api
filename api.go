package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	gg "github.com/gookit/color"
	"github.com/gorilla/mux"
)

// model
type Name2 struct {
	Filename string
	Reso     string
	Took     string
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
	name2.Reso = reso1
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

func main() {
	port := "1337"

	// this makes me look cool
	print("\033[H\033[2J")
	gg.Blue.Println("listening on port:", port)

	// start listening
	r := mux.NewRouter()
	r.HandleFunc("/download", downloadAnimu)
	err := http.ListenAndServe(":"+port, r)

	if err != nil {
		gg.Red.Println(err)
		return
	}
}
