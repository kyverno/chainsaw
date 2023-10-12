package client

import (
	"fmt"

	petname "github.com/dustinkirkland/golang-petname"
)

func Pet() string {
	return fmt.Sprintf("chainsaw-%s", petname.Generate(2, "-"))
}
