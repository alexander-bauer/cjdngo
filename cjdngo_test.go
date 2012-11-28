package cjdngo

import (
	"testing"
)

//TestJSON reads a Conf object from ./example.conf, then writes it to ./temp.conf.
func TestJSON(t *testing.T) {
	conf, err := ReadConf("./example.conf")
	if err != nil {
		t.Fatal(err)
	}
	err = WriteConf("./temp.conf", *conf)
	if err != nil {
		t.Fatal(err)
	}
}
