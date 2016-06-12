package radiowatch

import "bytes"

type CsvWriter struct {
	fileWriter
}

const _separator string = "\t"

/*
Returns a new instance of CsvWriter.
Takes a path name at which the results are saved.
 */
func NewCsvWriter(path string) CsvWriter {
	w := *new(CsvWriter)
	w.setPath(path)
	return w
}

func (c CsvWriter) Write(ti TrackInfo) {
	var b bytes.Buffer
	b.WriteString(ti.Artist)
	b.WriteString(_separator)
	b.WriteString(ti.Title)
	b.WriteString(_separator)
	b.WriteString(ti.Station)
	b.WriteString(_separator)
	b.WriteString(ti.CrawlTime.String())
	b.WriteString("\n")

	c.writeFile(c.Path + ti.NormalizedStationName() + ".csv", &b, ti.Station)
}