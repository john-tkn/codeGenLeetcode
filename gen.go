package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type codeVars struct {
	varType string
	varName string
}

var CFileHeaders string = "#include <stdio.h>\n#include <string.h>\n#include <stdbool.h>\n#include <stdlib.h>\n\n\n\n"

func getFunctionName(codeLineIn string) string {
	codeLineIn = strings.Replace(codeLineIn, "(", " ", -1)
	fmt.Println(codeLineIn + "\n")
	parts := strings.Split(codeLineIn, " ")
	fmt.Println(parts)
	return parts[1]
}

func getVariableNames(codeLineIn string) []string {
	//remove extra symbols
	codeLineIn = strings.Replace(codeLineIn, "(", " ", -1)
	codeLineIn = strings.Replace(codeLineIn, ",", "", -1)
	codeLineIn = strings.Replace(codeLineIn, ")", "", -1)
	codeLineIn = strings.Replace(codeLineIn, "{", "", -1)
	//split it by spaces
	parts := strings.Split(codeLineIn, " ")
	//return the vars and var names
	return parts
}

func main() {
	//get the code

	//keeps track of what test case we are on for naming vars
	var testCaseAccumulator int = 0

	var codeLines []string
	var codeLineUserIn string

	var ltDescin string
	var ltDesc []string

	var ltProbNum string
	scanner := bufio.NewScanner(os.Stdin)

	/*
		Get the problem Description

	*/
	fmt.Printf("LTgen v2.0\n\n")
	fmt.Printf("Enter the problem Number\n")
	scanner.Scan()
	ltProbNum = scanner.Text()

	fmt.Printf("Enter the problem Description. Press . when done\n")
	for true {
		scanner.Scan()
		ltDescin = scanner.Text()
		if ltDescin == "." {
			break
		} else {
			ltDesc = append(ltDesc, ltDescin)
		}
	}

	//get code from user
	fmt.Printf("Enter the code. Press . when done\n")

	for true {
		scanner.Scan()
		codeLineUserIn = scanner.Text()
		if codeLineUserIn == "." {
			break
		} else {
			codeLines = append(codeLines, codeLineUserIn)
		}
	}
	//process first line
	codeVars := getVariableNames(codeLines[0])

	//create the folder for all the code
	//fileFolderName := fmt.Sprint(codeVars[1], ltProbNum)
	fileFolderName := filepath.Join(".", fmt.Sprint(codeVars[1], ltProbNum))

	fmt.Println(fileFolderName)
	os.Mkdir(fileFolderName, os.ModePerm)

	//create the txt file with problem desc

	//txtfileName := fmt.Sprintf("./%s/%s%s.txt", fileFolderName, codeVars[1], ltProbNum)

	txtfileName := filepath.Join(".", fileFolderName, fmt.Sprintf("%s%s.txt", codeVars[1], ltProbNum))

	//fmt.Println("codeVars[1]:", codeVars[1])
	//fmt.Println("ltProbNum:", ltProbNum)
	//fmt.Println(txtfileName)
	//create the file
	txtFileDesc, err := os.Create(txtfileName)
	for i := range ltDesc {
		_, err = txtFileDesc.WriteString(ltDesc[i] + "\n")
		check(err)
	}

	defer txtFileDesc.Close()

	/**

	Write PROBLEM.h
	*/
	//name the file according to the name of the function

	FolderSrcPath := filepath.Join(".", fileFolderName, "src")

	os.Mkdir(FolderSrcPath, os.ModePerm)
	//fileNameHeaderFile := fmt.Sprintf("%s/src/%s.h", fileFolderName, codeVars[1])
	fileNameHeaderFile := filepath.Join(".", fileFolderName, "src", fmt.Sprintf("%s.h", codeVars[1]))

	//create the file
	fHeaderFile, err := os.Create(fileNameHeaderFile)
	check(err)

	headerFileForTestText := fmt.Sprintf("#ifndef %s_H\n#define %s_H\n", strings.ToUpper(codeVars[1]), strings.ToUpper(codeVars[1]))
	fHeaderFile.WriteString(headerFileForTestText)

	firstLine := strings.Replace(codeLines[0], " {", ";", -1)
	fHeaderFile.WriteString(firstLine + "\n")

	headerFileForTestText = fmt.Sprintf("#endif %s_H", strings.ToUpper(codeVars[1]))
	fHeaderFile.WriteString(headerFileForTestText)
	defer fHeaderFile.Close()

	/**

	Write PROBLEM.c
	*/

	//name the file according to the name of the function
	//fileName := fmt.Sprintf("%s/src/%s.c", fileFolderName, codeVars[1])

	fileName := filepath.Join(".", fileFolderName, "src", fmt.Sprintf("%s.c", codeVars[1]))

	//create the file
	f, err := os.Create(fileName)
	check(err)
	defer f.Close()

	testFunctionHeaders := fmt.Sprintf("#include \"../src/%s.h\"\n", codeVars[1])

	_, err = f.WriteString(testFunctionHeaders)

	_, err = f.WriteString(CFileHeaders)
	check(err)

	//write the code given earlier (the function)
	for i := range codeLines {
		_, err = f.WriteString(codeLines[i] + "\n")
		check(err)
	}
	_, err = f.WriteString("\n\n\n")

	//write main
	//_, err = f.WriteString("int main(int argc, char *argv[]) {\n")

	defer f.Close()

	/*

		Write tests/test_problemname.c

	*/

	//name the file according to the name of the function

	testsDir := filepath.Join(".", fileFolderName, "tests")

	os.Mkdir(testsDir, os.ModePerm)

	fileNameTests := filepath.Join(".", fileFolderName, "tests", fmt.Sprintf("%s_test.c", codeVars[1]))

	fmt.Println(fileNameTests)

	//fileName = fmt.Sprintf("%s/tests/%s_test.c", fileFolderName, codeVars[1])

	//create the file
	f, err = os.Create(fileNameTests)
	check(err)
	//write headers
	_, err = f.WriteString(CFileHeaders)
	testFunctionHeaders = fmt.Sprintf("#include \"../src/%s.h\"\n", codeVars[1])
	_, err = f.WriteString(testFunctionHeaders)

	/*
		Get the Values of the test case
		and write them
	*/

	//Get the testcode variables
	i := 2
	for true {
		_, err = f.WriteString("//test case" + fmt.Sprintf("%d", testCaseAccumulator) + "\n")
		_, err = f.WriteString("void test_case_" + fmt.Sprintf("%d", testCaseAccumulator) + "() {")
		for i < len(codeVars)-1 {
			//print header for test case

			fmt.Printf("Please enter a value for %s %s\n", codeVars[i], codeVars[i+1])
			scanner.Scan()
			codeLineUserIn = scanner.Text()
			//check if the test case value is an array
			if !strings.Contains(codeLineUserIn, "[") || !strings.Contains(codeLineUserIn, "{") && !strings.Contains(codeVars[i], "*") {
				//not an array
				codeVarNoPoint := strings.Replace(codeVars[i], "*", "", -1)
				_, err = f.WriteString("\t" + codeVarNoPoint + " " + codeVars[i+1] + fmt.Sprintf("%d", testCaseAccumulator) + " = " + codeLineUserIn + ";\n")
				check(err)
			} else {
				//is an array
				codeVarNoPoint := strings.Replace(codeVars[i], "*", "", -1)
				codeLineUserIn = strings.Replace(codeLineUserIn, "[", "{", -1)
				codeLineUserIn = strings.Replace(codeLineUserIn, "]", "}", -1)
				_, err = f.WriteString("\t" + codeVarNoPoint + " " + codeVars[i+1] + fmt.Sprintf("%d", testCaseAccumulator) + "[]" + " = " + codeLineUserIn + ";\n")
			}
			check(err)
			i += 2
		}
		//get the answer to the test case
		fmt.Printf("Enter the expected answer\n")
		scanner.Scan()
		codeLineUserIn = scanner.Text()
		answers := codeLineUserIn
		//create a var and hold the answer to the function
		_, err = f.WriteString("\t" + codeVars[0] + " result" + fmt.Sprintf("%d", testCaseAccumulator) + " = " + codeVars[1] + "(")
		//write the vars we created earlier. Start at 2 to avoid function var type and name
		j := 2
		for j < len(codeVars)-1 {
			_, err = f.WriteString(codeVars[j+1] + fmt.Sprintf("%d", testCaseAccumulator))
			if j+2 < len(codeVars)-1 {
				_, err = f.WriteString(", ")
			}

			j += 2
		}
		_, err = f.WriteString(");\n")

		/*
			Start writing if statements to check if the result
			is the expected answer
		*/

		//see if expected answer is an array
		if strings.Contains(answers, "[") || strings.Contains(answers, "{") {
			answers = strings.Replace(answers, "[", "", -1)
			answers = strings.Replace(answers, "]", "", -1)
			answers = strings.Replace(answers, " ", "", -1)
			answers = strings.Replace(answers, ",", " ", -1)
			var proccessedAnswers []string
			proccessedAnswers = strings.Split(answers, " ")

			for k := range proccessedAnswers {
				//
				_, err = f.WriteString("\tif(")
				_, err = f.WriteString("result" + fmt.Sprintf("%d", testCaseAccumulator) + "[" + fmt.Sprintf("%d", k) + "] != " + proccessedAnswers[k] + ") {\n\t\treturn -1;\n\t}\n")
			}
			//check if the anwser is a string
		} else if strings.Contains(answers, "\"") {
			_, err = f.WriteString("\tif(")
			//_, err = f.WriteString("result"+ fmt.Sprintf("%d", testCaseAccumulator) + " != " + answers + ") {\n\t\treturn -1;\n\t}\n")
			answers = strings.Replace(answers, "\"", "", -1)
			_, err = f.WriteString("strcmp(result" + fmt.Sprintf("%d", testCaseAccumulator) + ", \"" + answers + "\") == 1 ) {\n\t\tprintf(\"Test" + fmt.Sprintf("%d", testCaseAccumulator) + " FAILED\\n\");\n\t} else {\n\tprintf(\"Test" + fmt.Sprintf("%d", testCaseAccumulator) + " PASSED\\n\");\n\t}\n}\n")
		} else {
			_, err = f.WriteString("\tif(")
			_, err = f.WriteString("result" + fmt.Sprintf("%d", testCaseAccumulator) + " != " + answers + ") {\n\t\tprintf(\"Test" + fmt.Sprintf("%d", testCaseAccumulator) + " FAILED\\n\");\n\t} else {\n\tprintf(\"Test" + fmt.Sprintf("%d", testCaseAccumulator) + " PASSED\\n\");\n\t}\n}\n")
		}

		testCaseAccumulator++
		//ask for another test case
		fmt.Printf("Whould you like to continue? y/n\n")
		scanner.Scan()
		codeLineUserIn = scanner.Text()
		if strings.Contains(codeLineUserIn, "y") || strings.Contains(codeLineUserIn, "Y") {
			_, err = f.WriteString("\n\n")
			i = 2
		} else if strings.Contains(codeLineUserIn, "n") || strings.Contains(codeLineUserIn, "N") {
			//write main
			_, err = f.WriteString("int main(int argc, char *argv[]) {\n")
			for i = 0; i < testCaseAccumulator; i++ {
				f.WriteString("\ttest_case_" + fmt.Sprintf("%d", i) + "();")
			}
			f.WriteString("\n}")
			break
		}
	}

	//write the makefile

	//name the file according to the name of the function
	//fileNameMake := fmt.Sprintf("%s/Makefile", fileFolderName)
	fileNameMake := filepath.Join(fileFolderName, "Makefile")

	//create the file
	fMakeFile, err := os.Create(fileNameMake)

	fMakeFile.WriteString(codeVars[1] + ":\n\tcc ./tests/" + codeVars[1] + "_test.c" + " ./src/" + codeVars[1] + ".c")
	defer fMakeFile.Close()

}
