package tasker

import (
	"testing"
)

func TestCheckDb(t *testing.T) {
	ok, err := checkDb()
	println("ok:", ok)
	println("err: ", err.Error())
	if err != nil {
		t.Error("some error")
	}

}
