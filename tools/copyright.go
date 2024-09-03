package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// getCopyrightYear get copyright year from git log and file modify time.
func getCopyrightYear(filePath string) (string, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	modYear := fmt.Sprintf("%d", file.ModTime().Year())
	cmd := exec.Command("bash", "-c", "git log --follow --format=%cd --date=format:%Y "+filePath+" | sort -u")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if lines[0] != "" && lines[0] != modYear {
		return lines[0] + "-" + modYear, nil
	}
	return modYear, nil
}

// main function to traveser docs folder and update copyright year.
func main() {
	var wg sync.WaitGroup
	err := filepath.Walk("docs/", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()

			file, err := os.ReadFile(path)
			if err != nil {
				return
			}

			year, err := getCopyrightYear(path)
			if err != nil {
				return
			}

			replacedFile := strings.ReplaceAll(string(file), "<copyright-year>", year)
			err = os.WriteFile(path, []byte(replacedFile), 0644)
			if err != nil {
				return
			}

			println("Copyright Years: " + year + " " + path)
		}()
		return nil
	})
	wg.Wait()
	if err != nil {
		return
	}
}
