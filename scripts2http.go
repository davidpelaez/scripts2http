package main

import (
	"flag"
	"fmt"
	//"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	//"path"
)

var absDir string
var address = flag.String("address", "0.0.0.0:8080", "interface and port to bind to")
var scriptsDir = flag.String("scripts-dir", ".", "Path to folder with the scripts to be exposed")

func ExecuteScript(scriptPath string, callArguments []string) (string, int) {

	cmd := exec.Command(scriptPath, callArguments...)
	cmd.Dir = *scriptsDir

	eventOutput, err := cmd.CombinedOutput()

	if err != nil {
		return err.Error(), 1
	}

	return string(eventOutput), 0
}

func ScriptHandler(w http.ResponseWriter, r *http.Request) {
	// the first element in empty, because if starts with /
	trimmedPath := strings.Trim(r.URL.Path, "/")
	urlParts := strings.Split(trimmedPath, "/")

	scriptName := urlParts[0]
	callArguments := urlParts[1:]

	if scriptName == "" {
		// didn't pass the desirde script
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error, no script requested")
		return
	}

	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	fmt.Println("Error getting request body. Error was:", err)
	//}

	scriptPath := filepath.Join(absDir, scriptName)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		fmt.Printf("Cannot run %s. Script not found\n", scriptName)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Script not found")
		return
	}

	fmt.Println("Executing", scriptName)

	w.Header().Set("Content-Type", "text/plain")
	if output, exit := ExecuteScript(scriptPath, callArguments); exit != 0 {
		fmt.Println(scriptName, "failed with error:", output)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", output)
	} else {
		fmt.Fprintf(w, output)
	}

}

func main() {
	flag.Parse()
	fmt.Printf("Exposing scripts in '%v' on %v\n", *scriptsDir, *address)
	absDir, _ = filepath.Abs(*scriptsDir)
	http.HandleFunc("/", ScriptHandler)
	panic(http.ListenAndServe(*address, nil))
}
