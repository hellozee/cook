package manager

import (
	"encoding/json"
	"hash/crc32"
	"io/ioutil"
	"os"

	lg "../logger"
	ps "../parser"
)

//Entity  Data Structure for holding the file name and hash of an entity
type Entity struct {
	File string `json:"file"`
	Hash uint32 `json:"hash"`
}

//Parent  Data Structure to hold multiple Entity elements
type Parent struct {
	Body struct {
		Entity []Entity `json:"entity"`
	} `json:"body"`
}

//Manager  Data Structure to hold and operate on details.json
type Manager struct {
	FileData       string
	NewFileTimings map[string]uint32
	OldFileTimings map[string]uint32
	FileList       map[string]string
	HashJSONold    Parent
	HashJSONnew    Parent
	Logger         *lg.Logger
}

//ReadDetails  Reading from the details.json file
func (man *Manager) ReadDetails() error {
	jsonFile, err := os.Open("Cooking/details.json")
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(bytes, &man.HashJSONold)

	if err != nil {
		//Someone has tampered with the JSON file
		os.Remove("Cooking/details.json")
		man.Logger.ReportError(err.Error())
		return err
	}

	for _, item := range man.HashJSONold.Body.Entity {
		man.OldFileTimings[item.File] = item.Hash
	}
	man.Logger.ReportSuccess("Successfully read details.json")
	return nil
}

//WriteDetails  Writing the new data onto details.json
func (man *Manager) WriteDetails() error {
	jsonData, err := json.MarshalIndent(man.HashJSONnew, "", " ")

	if err != nil {
		man.Logger.ReportError(err.Error())
		return err
	}

	file, err := os.OpenFile("Cooking/details.json",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)

	defer file.Close()

	if err != nil {
		man.Logger.ReportError(err.Error())
		return err
	}

	_, err = file.Write(jsonData)

	if err != nil {
		man.Logger.ReportError(err.Error())
		return err
	}
	man.Logger.ReportSuccess("Successfully written to details.json")
	return nil
}

//GenerateFileList  Generating the files list which we are going to compile
func (man *Manager) GenerateFileList(par ps.Parser, tag string) error {
	details := par.FileDetails[tag]

	_, err := os.Stat(details.File)

	if err != nil {
		man.Logger.ReportError(err.Error())
		return err
	}

	man.FileList[tag] = details.File

	if details.Deps == nil {
		return nil
	}

	for _, name := range details.Deps {
		err = man.GenerateFileList(par, name)

		if err != nil {
			man.Logger.ReportError(err.Error())
			return err
		}
	}

	man.Logger.ReportSuccess("Successfully generated file list")

	return nil
}

//GenerateList  Generate a brand new details.json
func (man *Manager) GenerateList() error {
	for _, value := range man.FileList {
		file, err := ioutil.ReadFile(value)
		if err != nil {
			man.Logger.ReportError(err.Error())
			return err
		}
		hash := HashFile(file)
		man.NewFileTimings[value] = hash
		man.HashJSONnew.Body.Entity = append(man.HashJSONnew.Body.Entity,
			Entity{File: value, Hash: hash})
	}
	man.Logger.ReportSuccess("Successfully generated details.json")
	return nil
}

//NewManager  Helper function to create a new manager
func NewManager(log *lg.Logger) (Manager, error) {
	temp, err := ioutil.ReadFile("Recipe")

	if err != nil {
		//Missing Recipe File
		log.ReportError(err.Error())
		return Manager{}, err
	}

	recipe := string(temp)

	man := Manager{
		FileData:       recipe,
		NewFileTimings: make(map[string]uint32),
		OldFileTimings: make(map[string]uint32),
		FileList:       make(map[string]string),
		Logger:         log,
	}
	man.Logger.ReportSuccess("Successfully created a Manager Object")
	return man, nil
}

//HashFile  Obtaining the has of the passed file
func HashFile(file []byte) uint32 {
	hash := crc32.ChecksumIEEE(file)
	return hash
}

//CheckHash  Comparing hashes of the passed file with the previous hash
func CheckHash(file []byte, hash uint32) bool {
	generatedHash := crc32.ChecksumIEEE(file)
	return generatedHash == hash
}
