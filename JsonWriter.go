package radiowatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

/*
JsonWriter implements the Writer interface.
It takes the TrackInfo, converts it to a JSON string and writes it to a file.
 */
type JsonWriter struct {
	fileWriter
}

/*
Returns a new instance of JsonWriter.
Takes a path name at which the results are saved.
 */
func NewJsonWriter(path string) JsonWriter {
	w := *new(JsonWriter)
	w.setPath(path)
	return w
}

/*
Concrete implementation of Writer.Write()
 */
func (j JsonWriter) Write(trackInfo TrackInfo) {
	buffer, err := json.Marshal(trackInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while converting result from station %s to json: %s\n", trackInfo.Station, err.Error())
	}
	b := bytes.NewBuffer(buffer)
	b.WriteString("\n")

	j.writeFile(j.Path + trackInfo.NormalizedStationName() + ".rwjson", b, trackInfo.Station)
}
