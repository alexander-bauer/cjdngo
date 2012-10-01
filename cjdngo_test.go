package cjdngo

import (
	"testing"
)

//TestJSON reads a Conf object from ./example.conf, then logs it.
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
