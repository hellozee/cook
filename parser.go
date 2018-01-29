package main

import (
	"strings"
)

func getName(file string) string {
	//This fetches the tag of the detail block
	retVal := ""
	for i := 0; i < len(string(file[:])); i++ {
		if string(file[i]) == "{" {
			break
		}
		retVal += string(file[i])
	}
	retVal = strings.TrimSpace(retVal)
	return retVal
}

func getIdentifier(file string) string {
	//Self explanatory, this fetches the Identifier
	retVal := ""
	for i := 0; i < len(string(file[:])); i++ {
		if string(file[i]) == "=" {
			break
		}
		retVal += string(file[i])
	}
	retVal = strings.TrimSpace(retVal)
	return retVal
}

func getParams(file string) string {
	//Fetching the Parameters after fetching the Identifier
	retVal := ""
	count := 2
	for i := 0; i < len(string(file[:])); i++ {
		if count == 0 {
			break
		}
		if string(file[i]) == "'" {
			count--
		}
		if count == 1 {
			retVal += string(file[i])
		}
	}
	retVal = strings.TrimSpace(retVal)
	return retVal[1:]
}

func checkEnd(file string) bool {
	//Checks whether the detail block has ended or not
	file = strings.Trim(file[:], "\n")
	file = strings.TrimSpace(file)
	if string(file[:]) == "}" {
		return true
	}
	return false
}

func parse(file string) {
	compiler := false
	//Removing the blank lines
	file = strings.Replace(file, "\n\n", "\n", -1)
	divided := strings.Split(file, "\n")
	for i := 0; i < len(divided); i++ {
		name := getName(divided[i])
		if name == "$" {
			compiler = true
		} else {
			compiler = false
		}
		i++
		for !checkEnd(divided[i]) {
			if compiler {
				identifier := getIdentifier(divided[i])
				param := getParams(divided[i])
				fillCompilerDetails(identifier, param)
			} else {
				identifier := getIdentifier(divided[i])
				param := getParams(divided[i])
				fillFileDetails(name, identifier, param)
			}
			i++
		}
	}
}
