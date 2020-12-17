package file

import (
	"errors"
	"testing"
	"time"

	"github.com/jbsmith7741/trial"
)

const filePath = "../../../test/"

type SimpleStruct struct {
	Name   string
	Value  int
	Enable bool
	Dura   time.Duration
	//Time   time.Time `format:"2006-01-02"`
}

func TestLoad(t *testing.T) {
	fn := func(args ...interface{}) (interface{}, error) {
		c := &SimpleStruct{}
		f := args[0].(string)
		err := Load(f, c)
		return c, err
	}
	cases := trial.Cases{
		"toml": {
			Input: filePath + "test.toml",
			Expected: &SimpleStruct{
				Name:   "toml",
				Value:  10,
				Enable: true,
				Dura:   10 * time.Second},
		},
		"json": {
			Input: filePath + "test.json",
			Expected: &SimpleStruct{
				Name:   "json",
				Value:  10,
				Enable: true,
				//Dura: 10 * time.Second, //TODO add support
			},
		},
		"yaml": {
			Input:    filePath + "test.yaml",
			Expected: &SimpleStruct{Name: "yaml", Value: 10, Enable: true, Dura: 10 * time.Second},
		},
		"unknown": {
			Input:       "test.unknown",
			ExpectedErr: errors.New("unknown file type"),
		},
		"missing file": {
			Input:       filePath + "missing.toml",
			ExpectedErr: errors.New("no such file or directory"),
		},
	}
	trial.New(fn, cases).Test(t)
}
