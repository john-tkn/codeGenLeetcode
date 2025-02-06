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

	/*
		Get the function

	*/

	//get code from user
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
	//name the file according to the name of the function
	fileName := fmt.Sprintf("./%s.c", codeVars[1])
	//create the file
	f, err := os.Create(fileName)
	check(err)
	defer f.Close()

	
	_, err = f.WriteString(CFileHeaders)
	check(err)

	//write the code given earlier (the function)
	for i:= range codeLines {
		_, err = f.WriteString(codeLines[i] + "\n")
		check(err)
	}
	_, err = f.WriteString("\n\n\n")
	
	//write main
    _, err = f.WriteString("int main(int argc, char *argv[]) {\n")

	defer f.Close()
	
	/*
		Get the Values of the test case
		and write them
	*/

	//Get the testcode variables
	i := 2
	for true {
		for i < len(codeVars)-1 {
			//print header for test case
			_, err = f.WriteString("//test case" + fmt.Sprintf("%d", testCaseAccumulator) + "\n");
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
			_, err = f.WriteString("\t" + codeVarNoPoint + " " + codeVars[i+1] + fmt.Sprintf("%d", testCaseAccumulator)+ "[]" + " = " + codeLineUserIn + ";\n")
			}
			check(err)
			i+=2
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

			j+=2
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
				_, err = f.WriteString("result"+ fmt.Sprintf("%d", testCaseAccumulator) + "[" + fmt.Sprintf("%d", k) + "] != " + proccessedAnswers[k] + ") {\n\t\treturn -1;\n\t}\n")
			}
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
			break
		}
	}

	//write the final return for if all the functions ran correctly
	_, err = f.WriteString("\treturn 0;\n}")

}