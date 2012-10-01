package cjdngo

import (
	"testing"
	"encoding/json"
)

func TestJSON(t *testing.T) {
	conf, err := ReadConf("./example.conf")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	println(string(b))
}
