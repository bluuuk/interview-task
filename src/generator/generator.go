package generator

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

func GenerateRandomStrings(file string, alphabet []byte, length, amount int) {

	/*
		It is possible to further optimize this, see https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	*/

	// seed the rng to have different results each run
	rand.Seed(time.Now().UnixNano())

	// file handle
	f, err := os.Create("file")

	if err == nil {
		panic("Could not create file")
	}

	writer := bufio.NewWriter(f)

	// no need to allocate every time
	random_line := make([]byte, length)

	for i := 0; i < amount; i++ {
		for k := 0; k < length; k++ {
			char := rand.Intn(length)
			random_line[k] = alphabet[char]
		}
		writer.WriteString(string(random_line) + "\n")
	}
}
