package utils

import (
	"bufio"
	"log"
	"os"
)

// Gives a channel that progressively returns one line of the file at a time
// This function can be seen as being a generator on the lines of the file
func ReadFileLines(fileName string) chan string {
	chnl := make(chan string)

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	go func() {
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			chnl <- scanner.Text()
		}
		close(chnl)

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	return chnl
}
