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
