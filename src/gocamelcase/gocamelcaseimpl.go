package main

import (
	"fmt"
	"os"
	"strings"
	"gocamelcase/oxfordDict"
)

func processSentence(sen string) string {
	var dp []bool
	var prevWord []int

	for i := 0; i < len(sen); i++ {
		dp = append(dp, false)
		prevWord = append(prevWord, -1)
		if (i ==0 || (i > 0 && dp[i-1] == true)) && (string(sen[i]) == "i" || string(sen[i]) == "I") {
			dp[i] = true
			prevWord[i] = i-1
			continue
		}
		for j := i-2; j > -2; j-- {
			if j == -1 || dp[j] == true {
				exists := oxfordDict.Connect(sen[j+1:i+1])
				// fmt.Println(sen[j+1:i+1], exists)
				if exists == true {
					dp[i] = true
					if i > 0 {
						prevWord[i] = j
					}
					break
				}
			}
		}
	}

	cSen := ""
	var res []int
	if dp[len(sen) - 1] == true {
		for i:=len(sen)-1; i>0; {
			if prevWord[i] == -1 {
				break
			}
			res = append(res, prevWord[i]+1)
			i = prevWord[i]
		}
		// fmt.Println(res)
		cSen = makeCamelCase(sen, res)

	}
	return cSen
}

func makeCamelCase(sentence string, res []int) string {
	for i:=0; i<len(res); i++  {
		letter := string(sentence[res[i]])
		sentence = sentence[:res[i]] + strings.ToUpper(letter) + sentence[res[i]+1:]
	}
	return sentence
}

func main() {
	args := os.Args
	input := args[1]
    fmt.Println("Checking for", input)
	fmt.Println("Result:", processSentence(input))

}
