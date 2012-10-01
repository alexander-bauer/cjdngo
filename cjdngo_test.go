package cjdngo

import "testing"

func TestJSON(t *testing.T) {
	conf := ReadConf("./example.conf")
	_ := WriteConf("./output.conf", conf)
}
