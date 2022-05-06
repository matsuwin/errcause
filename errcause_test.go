package errcause

import (
	"github.com/pkg/errors"
	"os"
	"testing"
)

func Test(t *testing.T) {
	defer Recover()

	if err := mkError(); err != nil {
		panic(err)
	}
}

func mkError() (_ error) {
	_, err := os.ReadFile("xxx.txt")
	if err != nil {
		return errors.New(err.Error())
	}
	return
}
