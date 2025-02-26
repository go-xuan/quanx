package miniox

import (
	"testing"
)

func TestMinio(t *testing.T) {
	if err := (&Config{
		Host:         "",
		Port:         0,
		AccessId:     "",
		AccessSecret: "",
		SessionToken: "",
		Secure:       false,
		BucketName:   "",
		PrefixPath:   "",
		Expire:       0,
	}).Execute(); err != nil {
		t.Error(err)
	}
}
