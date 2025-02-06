package main

import (
    "bufio"
    "fmt"
    "os"
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

func getFunctionName(codeLineIn string) string{
	codeLineIn = strings.Replace(codeLineIn, "(", " ", -1)
	fmt.Println(codeLineIn + "\n")
	parts := strings.Split(codeLineIn, " ")
	fmt.Println(parts)
	return parts[1]
}

func getVariableNames(codeLineIn string) []string{
	codeLineIn = strings.Replace(codeLineIn, "(", " ", -1)
	codeLineIn = strings.Replace(codeLineIn, ",", "", -1)
	codeLineIn = strings.Replace(codeLineIn, ")", "", -1)
	codeLineIn = strings.Replace(codeLineIn, "{", "", -1)
	parts := strings.Split(codeLineIn, " ")	
	return parts
}

func main() {
	//get the code

	var testCaseAccumulator int = 0

	var codeLines []string
	var codeLineUserIn string
	fmt.Printf("Enter the code. Press . when done\n")
	scanner := bufio.NewScanner(os.Stdin)
	for true{
		scanner.Scan()
		codeLineUserIn = scanner.Text()
		if codeLineUserIn == "." {
			break
		} else {
			codeLines = append(codeLines, codeLineUserIn)
		}
	}

	codeVars := getVariableNames(codeLines[0])

	fileName := fmt.Sprintf("./%s.c", codeVars[1])
	
	f, err := os.Create(fileName)
	check(err)
	defer f.Close()

	
	_, err = f.WriteString(CFileHeaders)
	check(err)

	for i:= range len(codeLines) {
		_, err = f.WriteString(codeLines[i] + "\n")
		check(err)
	}
	_, err = f.WriteString("\n\n\n")
	

    _, err = f.WriteString("int main(int argc, char *argv[]) {\n")

	defer f.Close()
	
	//Get the testcode variables
	i := 2
	for true {
		for i < len(codeVars)-1 {
			_, err = f.WriteString("//test case" + fmt.Sprintf("%d", testCaseAccumulator + "\n"));
			fmt.Printf("Please enter a value for %s %s\n", codeVars[i], codeVars[i+1])
			scanner.Scan()
			codeLineUserIn = scanner.Text()
		
			if !strings.Contains(codeLineUserIn, "[") || !strings.Contains(codeLineUserIn, "{") && !strings.Contains(codeVars[i], "*") {
				codeVarNoPoint := strings.Replace(codeVars[i], "*", "", -1)
				_, err = f.WriteString("\t" + codeVarNoPoint + " " + codeVars[i+1] + fmt.Sprintf("%d", testCaseAccumulator) + " = " + codeLineUserIn + ";\n")
			check(err)
			} else {
			codeVarNoPoint := strings.Replace(codeVars[i], "*", "", -1)
			codeLineUserIn = strings.Replace(codeLineUserIn, "[", "{", -1)
			codeLineUserIn = strings.Replace(codeLineUserIn, "]", "}", -1)
			_, err = f.WriteString("\t" + codeVarNoPoint + " " + codeVars[i+1] + fmt.Sprintf("%d", testCaseAccumulator)+ "[]" + " = " + codeLineUserIn + ";\n")
			}
			check(err)
			i+=2
		}
		//get the answer
		fmt.Printf("Enter the expected answer\n")
		scanner.Scan()
		codeLineUserIn = scanner.Text()
		answers := codeLineUserIn
		_, err = f.WriteString("\t" + codeVars[0] + " result" + fmt.Sprintf("%d", testCaseAccumulator) + " = " + codeVars[1] + "(")
		j := 2
		for j < len(codeVars)-1 {
			//codeVarNoPoint := strings.Replace(codeVars[j], "*", "", -1)
			_, err = f.WriteString(codeVars[j+1] + fmt.Sprintf("%d", testCaseAccumulator))
			if j+2 < len(codeVars)-1 {
				_, err = f.WriteString(", ")
			}

			j+=2
		}
		_, err = f.WriteString(");\n")


		//see if expected answer is an array
		if strings.Contains(answers, "[") || strings.Contains(answers, "{") {
			answers = strings.Replace(answers, "[", "", -1)
			answers = strings.Replace(answers, "]", "", -1)
			answers = strings.Replace(answers, " ", "", -1)
			answers = strings.Replace(answers, ",", " ", -1)
			var proccessedAnswers []string
			proccessedAnswers = strings.Split(answers, " ")
			
			for k := range len(proccessedAnswers) {
				//fmt.Println("result"+ fmt.Sprintf("%d", testCaseAccumulator) + "[" + fmt.Sprintf("%d", k) + "] != " + proccessedAnswers[k] + ") {\n\t\treturn -1;\n\t}\n")
				_, err = f.WriteString("\tif(")
				_, err = f.WriteString("result"+ fmt.Sprintf("%d", testCaseAccumulator) + "[" + fmt.Sprintf("%d", k) + "] != " + proccessedAnswers[k] + ") {\n\t\treturn -1;\n\t}\n")
			}
		}
		
		
		testCaseAccumulator++
		fmt.Printf("Whould you like to continue? y/n\n")
		scanner.Scan()
		codeLineUserIn = scanner.Text()
		if strings.Contains(codeLineUserIn, "y") || strings.Contains(codeLineUserIn, "Y") {
			_, err = f.WriteString("\n\n")
			i = 2
		} else if strings.Contains(codeLineUserIn, "n") || strings.Contains(codeLineUserIn, "N") {
			break
		}
	}
	_, err = f.WriteString("\treturn 0;\n}")

}