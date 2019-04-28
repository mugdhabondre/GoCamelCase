package gocamelcaseimpl

import (
	"strings"
	"github.com/mugdhabondre/gocamelcase/oxforddict"
)

// Method to process given phrase to find legitimate words
func ProcessPhrase(phrase string) (string, error) {

	// maintain information about which indices denote words
	// prevWord denotes information about which index before i was a word
	var isWord []bool
	var prevWord []int

	for i := 0; i < len(phrase); i++ {
		isWord = append(isWord, false)
		prevWord = append(prevWord, -1)

		// for single letters, only letter i is a valid pronoun, rest are not considered
		if ((i == 0 || isWord[i-1] == true) && (string(phrase[i]) == "i" || string(phrase[i]) == "I")) {
			isWord[i] = true
			prevWord[i] = i-1
			continue
		}

		// if first letter was not an  "i", go to next letter
		if i == 0 {
			continue
		}

		// find the complete word nearest to the left of index i, ignoring single letters
		for j := i-2; j > -2; j-- {
			if j == -1 || isWord[j] == true {

				// if a word is found till index j, check if phrase[j+1:i+1] is a word
				exists, error := oxforddict.ConnectAndCheck(phrase[j+1:i+1])

				if error != nil {
					return "", error
				}
				if exists == true {
					isWord[i] = true
					if i > 0 {
						prevWord[i] = j
					}
					break
				}
			}
		}
	}

	// find capitalized letter indices
	capitals := findWordIndices(phrase, isWord, prevWord)
	// Capitalize the letters
	resPhrase := makeCamelCase(phrase, capitals)

	return resPhrase, nil
}

//.Method to find indices which have completed words starting the zeroeth index
func findWordIndices(phrase string, isWord []bool, prevWord []int) []int {
	var capitalIndices []int
	lastWordIndex := len(phrase)-1

	// if last index of phrase is not a word, find the last index that is a word
	if isWord[lastWordIndex] != true {
		i := 0
		for i =len(isWord)-1; i>0 && isWord[i] != true; {
			i--
		}
		lastWordIndex = i
	}

	// Find the letters to capitalize, hop from one valid word to another starting from the end of the string
	for i:=lastWordIndex; i>0; {
		// break if we arrive at an index that is one word from start
		if prevWord[i] == -1 {
			break
		}
		capitalIndices = append(capitalIndices, prevWord[i]+1)
		i = prevWord[i]
	}

	return capitalIndices
}

// Method to capitalize letters at given indices
func makeCamelCase(sentence string, res []int) string {
	for i:=0; i<len(res); i++  {
		letter := string(sentence[res[i]])
		sentence = sentence[:res[i]] + strings.ToUpper(letter) + sentence[res[i]+1:]
	}
	return sentence
}

