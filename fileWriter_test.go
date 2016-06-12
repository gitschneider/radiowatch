package radiowatch

import (
	"testing"
	"bytes"
	"io/ioutil"
	"os"
)

func TestFileWriter_SetString(t *testing.T) {
	table := []struct {
		input    string
		expected string
	}{
		{"/path", "/path/"},
		{"/path/", "/path/"},
		{"path", "path/"},
		{"path/", "path/"},
	}
	fw := fileWriter{}
	for _, tl := range table{
		fw.setPath(tl.input)
		if fw.Path != tl.expected{
			t.Errorf("SetPath(%v): Expected %v, got %v", tl.input, tl.expected, fw.Path)
		}
	}
}

func TestFileWriter_WriteFile(t *testing.T) {
	fileContents := "This is a test for \n file handling."
	fw := fileWriter{""}
	buf := bytes.Buffer{}
	buf.WriteString(fileContents)
	fw.writeFile("test.test", &buf, "TestStation")
	defer os.Remove("test.test")

	con, err := ioutil.ReadFile("test.test")
	if err != nil {
		t.Errorf("WriteFile: Error when accessing the written file: %v", err.Error())
	}
	if string(con) != fileContents{
		t.Errorf("WriteFile(%v): Expected %v, got %v", `buf, "test.test"`, fileContents, string(con))
	}
}
