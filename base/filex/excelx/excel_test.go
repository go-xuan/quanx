package excelx

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-xuan/quanx/base/filex"
)

func TestExcelParse(t *testing.T) {
	var inputFile = "./user.xlsx"
	var headerMapping = map[string]string{
		"用户名": "name",
		"邮箱":  "email",
		"手机号": "mobile",
		"备注":  "comment",
	}
	if result, err := ReadXlsxWithMapping(inputFile, "", headerMapping); err != nil {
		fmt.Println(err)
	} else {
		var lines []string
		for _, row := range result {
			if mobile, ok := row["mobile"]; ok && mobile != "" {
				lines = append(lines, fmt.Sprintf(`update "magicuser" set "mobile" = '%s' where "email" = '%s';`, mobile, row["email"]))
			}
		}
		outputFile := fmt.Sprintf("./user_update_%s.sql", time.Now().Format("20060102"))
		if err = filex.WriteFileLine(outputFile, lines); err != nil {
			fmt.Println(err)
		}
	}
}
