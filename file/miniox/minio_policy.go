package miniox

import (
	"encoding/json"
)

type Policy struct {
	Version   string        `json:"Version"`
	Statement StatementList `json:"Statement"`
}

type StatementList []*Statement
type Statement struct {
	Action    []string  `json:"Action"`
	Effect    string    `json:"Effect"`
	Principal Principal `json:"Principal"`
	Resource  []string  `json:"Resource"`
	Sid       string    `json:"Sid"`
}

type Principal struct {
	AWS []string `json:"AWS"`
}

// 存储桶默认配置信息
func defaultBucketPolicy(bucketName string) string {
	if bytes, err := json.Marshal(Policy{
		Version: "2022-10-17",
		Statement: []*Statement{{
			Action:    []string{"s3:ObjectExist"},
			Effect:    "Allow",
			Principal: Principal{AWS: []string{"*"}},
			Resource:  []string{"arn:aws:s3:::" + bucketName + "/*"},
			Sid:       "",
		}},
	}); err != nil {
		return ""
	} else {
		return string(bytes)
	}
}
