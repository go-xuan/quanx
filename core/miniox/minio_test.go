package miniox

import (
	"testing"
	
	"github.com/go-xuan/quanx/core/configx"
)

func TestMinio(t *testing.T) {
	if err := configx.Execute(&Config{
		Host:         "",
		Port:         0,
		AccessId:     "",
		AccessSecret: "",
		SessionToken: "",
		Secure:       false,
		BucketName:   "",
		PrefixPath:   "",
		Expire:       0,
	}); err != nil {
		t.Error(err)
	}
}
