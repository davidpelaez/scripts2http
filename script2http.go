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
var port = flag.String("port", "8080", "TCP port to bind to")
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

	if len(urlParts) < 1 {
		// didn't pass the desirde script, error
	}

	scriptName := urlParts[0]
	callArguments := urlParts[1:]

	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	fmt.Println("Error getting request body. Error was:", err)
	//}

	fmt.Println("Executing", scriptName, callArguments)

	scriptPath := filepath.Join(absDir, scriptName)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		fmt.Printf("Cannot run %s. Script not found.\n", scriptName)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if output, exit := ExecuteScript(scriptPath, callArguments); exit != 0 {
		fmt.Println(scriptName, "failed with error:", output)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "%s", output)
	}

}

func main() {
	flag.Parse()
	fmt.Printf("Exposing scripts in '%v' on port %v\n", *scriptsDir, *port)
	absDir, _ = filepath.Abs(*scriptsDir)
	http.HandleFunc("/", ScriptHandler)
	panic(http.ListenAndServe(":"+*port, nil))
}
