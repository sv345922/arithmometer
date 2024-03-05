package tasker

import (
	"fmt"
	"testing"
)

func TestLoadDB(t *testing.T) {
	res, err := baseToSave.LoadDB()
	fmt.Printf("%v", res)
	println("err: ", err.Error())
	if err != nil {
		t.Error("some error")
	}

}
