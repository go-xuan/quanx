package filex

import (
	"fmt"
	"testing"
)

func TestFileSplit(t *testing.T) {
	filePath := "./nohup.log"
	fmt.Println(Analyse(filePath))
	files, err := FileSplit(filePath, 5000)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
}

func TestFileWrite(t *testing.T) {
	WriteFileString("nohup.log", "1111\n", Append)
	WriteFileString("nohup.log", "2222\n", Append)
	WriteFileString("nohup.log", "3333\n", Append)

	WriteFileLine("nohup.log", []string{
		"4444",
		"5555",
		"6666",
	}, Append)

	WriteFileLine("nohup.log", []string{
		"7777",
		"8888",
		"9999",
	}, Append)
}
