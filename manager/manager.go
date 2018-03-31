package manager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	ps "github.com/hellozee/cook/parser"
)

type entity struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

type parent struct {
	Body struct {
		Entity []entity `json:"entity"`
	} `json:"body"`
}

type Manager struct {
	FileData       string
	NewFileTimings map[string]string
	OldFileTimings map[string]string
	FileList       map[string]string
	HashJSONold    parent
	HashJSONnew    parent
}

func (man *Manager) ReadDetails() error {
	jsonFile, err := os.Open("Cooking/details.json")
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(bytes, &man.HashJSONold)

	if err != nil {
		//Someone has tampered with the JSON file
		os.Remove("Cooking/details.json")
		return errors.New("Error parsing Cooking/details.json\n" +
			"Please run the program again\n")
	}

	for _, item := range man.HashJSONold.Body.Entity {
		man.OldFileTimings[item.File] = item.Hash
	}

	return nil
}

func (man *Manager) WriteDetails() error {
	jsonData, err := json.MarshalIndent(man.HashJSONnew, "", " ")

	if err != nil {
		return err
	}

	file, err := os.OpenFile("Cooking/details.json",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)

	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)

	if err != nil {
		return err
	}

	return nil
}

func NewManager() (Manager, error) {
	temp, err := ioutil.ReadFile("Recipe")

	if err != nil {
		//Missing Recipe File
		return Manager{}, errors.New("No sane Recipe File found.\n" +
			"Make sure you have a Recipe file with proper syntax\n")
	}

	recipe := string(temp)

	man := Manager{
		FileData:       recipe,
		NewFileTimings: make(map[string]string),
		OldFileTimings: make(map[string]string),
		FileList:       make(map[string]string),
	}

	return man, nil
}

func (man *Manager) GenerateFileList(par ps.Parser, tag string) error {
	details := par.FileDetails[tag]

	_, err := os.Stat(details.File)

	if err != nil {
		return err
	}

	man.FileList[tag] = details.File

	if details.Deps == nil {
		return nil
	}

	for _, name := range details.Deps {
		err = man.GenerateFileList(par, name)

		if err != nil {
			return err
		}
	}

	return nil
}
