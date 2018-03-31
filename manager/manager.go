package manager

import (
	"encoding/json"
	"errors"
	"hash/crc32"
	"io/ioutil"
	"os"

	ps "github.com/hellozee/cook/parser"
)

type Entity struct {
	File string `json:"file"`
	Hash uint32 `json:"hash"`
}

type Parent struct {
	Body struct {
		Entity []Entity `json:"entity"`
	} `json:"body"`
}

type Manager struct {
	FileData       string
	NewFileTimings map[string]uint32
	OldFileTimings map[string]uint32
	FileList       map[string]string
	HashJSONold    Parent
	HashJSONnew    Parent
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

func (man *Manager) GenerateList() error {
	for _, value := range man.FileList {
		file, err := ioutil.ReadFile(value)
		if err != nil {
			return err
		}
		hash := HashFile(file)
		man.NewFileTimings[value] = hash
		man.HashJSONnew.Body.Entity = append(man.HashJSONnew.Body.Entity,
			Entity{File: value, Hash: hash})
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
		NewFileTimings: make(map[string]uint32),
		OldFileTimings: make(map[string]uint32),
		FileList:       make(map[string]string),
	}

	return man, nil
}

func HashFile(file []byte) uint32 {
	hash := crc32.ChecksumIEEE(file)
	return hash
}

//Comparing hashes of the current timestamp with the previous one
func CheckHash(file []byte, hash uint32) bool {
	generatedHash := crc32.ChecksumIEEE(file)
	return generatedHash == hash
}
