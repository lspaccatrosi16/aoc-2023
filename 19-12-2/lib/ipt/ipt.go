package ipt

import (
	"io"
	"os"
	"strings"
)

func GetInput() string {
	stdin, err := io.ReadAll(os.Stdin)

	if err != nil {
		panic(err)
	}

	str := string(stdin)

	return strings.TrimSuffix(str, "\n")
}
