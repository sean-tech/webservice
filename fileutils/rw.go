package fileutils

import (
	"bufio"
	"io"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	inputFile, inputError := os.Open(path)
	if inputError != nil {
		return nil, inputError
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	buf := make([]byte, 1024)
	for {

		n, err := inputReader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if (n == 0) { break}
	}
	return buf, nil
}