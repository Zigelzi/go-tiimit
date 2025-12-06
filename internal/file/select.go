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

	var fileNames []FileName
	for _, path := range filePaths {
		fileName := FileName{Path: filepath.Base(path)}
		date, err := ParseDate(path)
		if err != nil {
			// Errors from parsing don't need to stop the whole function.
			fmt.Println(err)
			continue
		}
		fileName.Date = date
		fileNames = append(fileNames, fileName)
	}

	if len(fileNames) == 0 {
		return "", fmt.Errorf("no files with yyyy-mm-dd date format in directory %s, [%v]", path, filePaths)
	}

	sortByNewestDate(fileNames)

	result, err := choose(fileNames)
	if err != nil {
		return "", err
	}

	return result, nil
}

func choose(fileNames []FileName) (string, error) {
	templates := &promptui.SelectTemplates{
		Inactive: "  {{ .Path }}",
		Active:   fmt.Sprintf("%s {{ .Path | underline }}", promptui.IconSelect),
		Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Path | faint }}`, promptui.IconGood),
	}

	prompt := promptui.Select{
		Label:     "Select file to import",
		Items:     fileNames,
		Templates: templates,
	}
	i, _, err := prompt.Run()

	if err != nil {
		fmt.Println("Unable to get input for selecting action")
		return "", err
	}

	return fileNames[i].Path, nil
}
