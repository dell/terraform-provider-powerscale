/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// getCopyrightYear get copyright year from git log and file modify time.
func getCopyrightYear(filePath string) (string, error) {
	currYear := fmt.Sprintf("%d", time.Now().Year())
	cmd := exec.Command("bash", "-c", "git log --follow --format=%cd --date=format:%Y "+filePath+" | sort -u")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	fmt.Println("git-log: (", string(output), ") ")
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	// if newly created
	if lines[0] == "" {
		return currYear, nil
	}

	startYear := lines[0]

	// check git status of file
	cmd = exec.Command("bash", "-c", "git status --porcelain "+filePath)
	output, err = cmd.Output()
	if err != nil {
		return "", err
	}

	// if modified
	if string(output) != "" {
		// if createYear equals to now
		if startYear == currYear {
			return currYear, nil
		}
		// otherwise
		return fmt.Sprintf("%s-%s", startYear, currYear), nil
	}

	// if not modified and created in this year
	if len(lines) == 1 {
		return startYear, nil
	}

	// if not modified and created in some other year
	endYear := lines[len(lines)-1]
	return startYear + "-" + endYear, nil
}

// main function to traveser docs folder and update copyright year.
func main() {

	err := filepath.Walk("docs/", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		year, err := getCopyrightYear(path)
		if err != nil {
			return err
		}

		replacedFile := strings.ReplaceAll(string(file), "<copyright-year>", year)
		err = os.WriteFile(path, []byte(replacedFile), 0644)
		if err != nil {
			return err
		}

		println("Copyright Years: " + year + " " + path)
		return nil
	})
	if err != nil {
		return
	}
}
