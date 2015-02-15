package gopywatch

import (
	"bufio"
	"fmt"
	"os"
)

func GetNextLine() (line string, err error) {
	fmt.Print(">>> ")
	in := bufio.NewReader(os.Stdin)
	line, err = in.ReadString('\n')
	return
}
