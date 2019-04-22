package main

import (
	"fmt"
	"strings"
	"CamelCaseWords/oxfordDict"
)

func processSentence(sen string) string {
	words := [10]string{"i", "like", "sam", "sung", "samsung", "dell", "love", "me", "mine", "apple"}
	var dp []bool
	var prevWord []int
	wordMap := make(map[string]float64) 

	for i:=0; i < len(words); i++ {
		wordMap[words[i]] = 1
	}

	for i := 0; i < len(sen); i++ {
		dp = append(dp, false)
		prevWord = append(prevWord, -1)
		for j := i-1; j > -2; j-- {
			if j == -1 || dp[j] == true {
				exists := oxfordDict.Connect(sen[j+1:i+1])
				fmt.Println(sen[j+1:i+1], exists)
				// _, exists := wordMap[sen[j+1:i+1]]
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
			res = append(res, prevWord[i]+1)
			i = prevWord[i]
		}

		cSen = makeCamelCase(sen, res)

	}
	return cSen
}

func makeCamelCase(sentence string, res []int) string {
	for i:=len(res)-1; i>0; i--  {
		letter := string(sentence[res[i]])
		sentence = sentence[:res[i]] + strings.ToUpper(letter) + sentence[res[i]+1:]
	}
	return sentence
}

func main() {
	input := "onetwothree"
    fmt.Println("Checking for ", input)
	fmt.Println("Result: ", processSentence(input))

	//fmt.Println(oxfordDict.Connect("i"))

}
