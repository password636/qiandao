
package main

import (
	"encoding/json"
	"fmt"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	_ "statik"
	"strings"
	"syscall"
	"regexp"
)

var signed_map = make(map[string]string, 0)
var unsigned_array = make([]string, 0)
var my_array = make([]string, 0)

func writeResponse(rw http.ResponseWriter, status int, v1 interface{}) (err error) {
	var b []byte
	b, err = json.Marshal(v1)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(status)
		rw.Write(b)
	}
	return
}

func getIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func names(w http.ResponseWriter, r *http.Request) {

	v1 := make(map[string][]string)

	content, _ := ioutil.ReadFile("signed.txt")
	my_array = strings.Split(string(content), "\r\n")
	for _, v := range my_array {
		if match, _ :=regexp.MatchString(`^\s*$`, v); ! match {
			v1["signed"] = append(v1["signed"], v)
		}
	}

	content1, _ := ioutil.ReadFile("unsigned.txt")
	my_array = strings.Split(string(content1), "\r\n")
	for _, v := range my_array {
		if match, _ :=regexp.MatchString(`^\s*$`, v); ! match {
			v1["unsigned"] = append(v1["unsigned"], v)
		}
	}

	writeResponse(w, http.StatusOK, v1)
}

func sign(w http.ResponseWriter, req *http.Request) {

	formData := map[string]interface{}{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&formData)
	if err != nil {
		log.Fatal(err)
	}

	name := formData["name"].(string)
	date := formData["date"].(string)

	if _, ok := signed_map["name"]; ok {

		w.WriteHeader(http.StatusConflict)

	} else {
		i := 0
		for in, v := range unsigned_array {
			if v == name {
				i = in
				break
			}
		}
		unsigned_array = append(unsigned_array[:i], unsigned_array[i+1:]...)
		signed_map[name] = date
		writeFile()
		names(w, req)
	}

}

func printLines(filePath string, values interface{}) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	switch v := values.(type) {
	case []string:
		for _, value := range v {
			fmt.Fprint(f, value+"\r\n")
		}
	case map[string]string:
		for k, value := range v {
			fmt.Fprint(f, k+"\t"+value+"\r\n")
		}
	}
	return nil
}

func writeFile() {
	printLines("unsigned.txt", unsigned_array)
	printLines("signed.txt", signed_map)
}

func main() {
	content, _ := ioutil.ReadFile("attendees.txt")
	my_array := strings.Split(string(content), "\r\n")
	for _, v := range my_array {
		if match, _ :=regexp.MatchString(`^\s*$`, v); ! match {
			unsigned_array = append(unsigned_array, v)
		}
	}
	

	writeFile()

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/names", names)
	http.HandleFunc("/api/sign", sign)

	http.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	var local, network string

	if runtime.GOOS == "windows" {
		local = "http://localhost:3000/"
		network = fmt.Sprintf("http://%v:3000/", getIP())
	} else {
		local = "\033[36mhttp://localhost:3000/\033[0m"
		network = fmt.Sprintf("\033[36mhttp://%v:3000/\033[0m", getIP())
	}

	log.Printf(`Started successfully!

You can now enter the following URL in your browser.

  Local:            ` + local + `

  On Your Network:  ` + network + `
`)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)

	<-sigchan
	log.Println("exit!")
}
