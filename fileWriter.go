package radiowatch

import (
	"strings"
	"bytes"
	"os"
	"fmt"
)

type fileWriter struct {
	Path string
}

/*
Set the Path at which the file will be saved.
 */
func (f *fileWriter) setPath(path string) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	f.Path = path
}

/*
Write Contents to the specified path
 */
func (f *fileWriter) writeFile(path string, content *bytes.Buffer, stationName string){
	file, err := os.OpenFile(path, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file %v: %v\n", path, err.Error())
	}
	defer file.Close()

	_, err = file.Write(content.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving results from %v to file %v: %v\n", stationName, path, err.Error())
	}
}