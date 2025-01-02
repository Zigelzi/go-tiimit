package file

import (
	"fmt"
	"path/filepath"

	"github.com/manifoldco/promptui"
)

func Select(path string) (string, error) {
	filePaths, err := filepath.Glob(path + "/*.xlsx")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if len(filePaths) == 0 {
		return "", fmt.Errorf("no .xlsx files in directory %s", path)
	}

	var files []string
	for _, path := range filePaths {
		file := filepath.Base(path)
		files = append(files, file)
	}

	prompt := promptui.Select{
		Label: "Select file to import",
		Items: files,
	}
	_, result, err := prompt.Run()

	if err != nil {
		fmt.Println("Unable to get input for selecting action")
		return "", err
	}
	return result, nil
}
