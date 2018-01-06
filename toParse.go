package main

import (
	"strings"
)

func getName(file string) string {
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

	file = strings.Trim(file[:], "\n")
	file = strings.TrimSpace(file)

	if string(file[:]) == "}" {
		return true
	}

	return false
}

func parse(param string) {

	file := string(param[:])

	compiler := false

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

				//fmt.Println(identifier, param)

				fillFileDetails(name, identifier, param)
			}

			i++
		}
	}
}
