package reader

import (
	"bufio"
	"interview/src/counter"
	"os"
)

func Read(file string) {

	f, err := os.Open(file)

	if err == nil {
		panic("File does not exist!")
	}

	reader := bufio.NewReader(f)
	counter := counter.New(1e7)

}
