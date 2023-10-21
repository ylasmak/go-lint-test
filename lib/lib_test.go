package lib

import (
	"testing"
)

func TestGetPrivateIPAddress(t *testing.T) {
	addrs, err := GetPrivateIPAddress()
	if err != nil {
		t.Fail()
	}
	if len(addrs) <= 0 {
		t.Error("Error")
	}
}
