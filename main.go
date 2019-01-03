package main

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var cliOutput bool
var checksumOnly bool
var filename string
var executable string

func init() {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	flag.BoolVar(&cliOutput, "i", false, "Outputs checksum to stdout instead of file.")
	flag.BoolVar(&checksumOnly, "c", false, "Disables filename print.")
	flag.StringVar(&filename, "o", "checksum.txt", "Sets output filename.")
	filename = filepath.Join(filepath.Dir(executable), filename)
	flag.Parse()
}

func main() {
	files := make([]string, 0)
	for _, arg := range os.Args[1:] {
		if isExecutable(arg) {
			continue
		}
		if stat, err := os.Stat(arg); err == nil && !stat.IsDir() {
			files = append(files, arg)
		}
	}

	if len(files) <= 0 {
		fmt.Println("No valid files specified")
		return
	}

	var output *os.File

	if cliOutput {
		output = os.Stdout
	} else {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			output, err = os.Create(filename)
			if err != nil {
				panic(err)
			}
		} else {
			output, err = os.OpenFile(filename, os.O_APPEND|os.O_RDWR, 0660)
			if err != nil {
				panic(err)
			}
		}

		defer output.Close()
	}

	checksumsChan := make(chan string)
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, filepath := range files {
		go func(filepath string) {
			checksum, err := calculateChecksum(filepath)
			var msg bytes.Buffer
			if err != nil {
				msg.WriteString(err.Error())
			} else {
				msg.WriteString(checksum)
				if !checksumOnly {
					msg.WriteString(fmt.Sprintf(" - %s", filepath))
				}
			}
			checksumsChan <- msg.String()
		}(filepath)
		wg.Add(1)
	}

	go func() {
		for checksum := range checksumsChan {
			fmt.Fprintln(output, checksum)
			wg.Done()
		}
	}()
}

func isExecutable(path string) bool {
	return path == executable || filepath.Clean(path) == filepath.Clean(filepath.Base(executable))
}

func calculateChecksum(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", errors.New("Unable to open file " + filepath)
	}
	hash := sha512.New384()
	defer f.Close()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
