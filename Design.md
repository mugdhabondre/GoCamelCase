# Design 
In this document, I am describing the _What_, _Why_ and _How_ of the application. I will also be analysing the performance of the application and suggest improvements. 


## Understanding the problem
So lets see what would be our inputs and output. Input to the application would be a string supplied as a field in the URL. This string could be a concatenation of words or just random letters. 
Output would also be a string. This string would contain the same letters as the input, however, the letters that are the first letter of a recognized word would be capitalized. Oh, and except the first letter of the string. 

For example:<br/>
Input: gameofthrones<br/>
Output: gameOfThrones<br/>

Input: helloworld<br/>
Output:helloWorld<br/>

Input: gocamelcase<br/>
Output: goCamelCase<br/>

So the real problem lies in identifying all words in the input string, or atleast a list of words which concatenated together will from the same order of letters (case insensitive) as the input. <br/>
Breaking it down further, <br/>
1. We need to find an application which can tell us if a given sequence of letters is a valid word<br/>
2. Find out the right words from a pool of words that can be formed from input's letters, that when concatenated, form the input itself. 

## Requirements
Now that the problem is a bit clear, lets try to see what are the requirements:

### Go
We will be developing our application in GoLang. No justification, just a choice.

### Oxford Dictionaries API
To recognize if a word is valid or not, we will be using the [Oxford Dictionaries API](https://developer.oxforddictionaries.com/). I will be using the ProtoType plan. 
I chose this because Oxford Dictionaries are one of the most trusted ones when it comes to English. <br/>
But is it one of the best from performance perspective of our application? We'll find out soon ;)

### Docker
I have used docker to generate my image and run the application as a container. Reason is nothing more than familarity. 

### Microsoft Azure
I deployed my application in Azure at the end. It was pretty easy to deploy it and also because I had a few free credits ;)

## Assumptions
1. I am expecting the input to be small, may be 7-60 letters. This is because the oxford dictionaries api allows only 60 request per minute. Also, it has a limit of 1000 requests per month. These limits can be increased by purchasing the enterprise or research version. <br/>
2. Another assumption is, I am considering the input to contain words that are in their root form. If a word is not in the root form, for example, _rocking_ instead of _rock_, it needs to be converted using the _lemmas_ api of oxford dictionaries. This requires an additional http request, and exceeds the limits specified. <br />
3. Oxford Dictionaries API assumes every letter of the alphabet to be a _noun_. So I am not checking for any single letters to be words, except for the letter _i_/_I_. 
4. I am filtering the API results for the following lexical categories: verb,adjective,adverb,conjunction,numeral,particle,preposition,pronoun,noun. Any parts of speech other than this, wont be recognized as a part of the algorithm.


## Core Algorithm
Now, lets take a look at the core algorithm. Lets assume we already have implemented a method - _ConnectAndCheck_ that connects to oxford dictionaries API and tells us if a given word is valid or not.

Basics of the algorithm is as follows: 

```bash
start from the last index
traverse the input in reverse
  When you come across an index(j) such that input[j:] is a valid word using _ConnectAndCheck_
    recursively call the algorithm on input[:j] to check if it can form a valid combination of words
      Break if yes, continue traversing if no
```

But recursive implementation has exponential time complexity. So, I have used Dynamic Programming Bottom Up approach to figure out if the string forms a combination of valid words at a particular index. <br/>

Following is the algorithm that I have implemented:

```bash
isWord = [] #Store information about valid words till index i
prevWord = [] #Store information about the previous valid chain of words till index i
FOR each index i in input
	FOR each index j in input, such that j < i, j ∈ {i-1, 0} #parse in reverse order
		isWord[i] = FALSE, prevWord[i] = -1 #initialization
		IF _ConnectAndCheck(input[j:i]) = TRUE && isWord[j-1] = TRUE THEN
			isWord[i] = TRUE 
			prevWord[i] = j-1
		ENDIF
	ENDFOR
ENDFOR

IF isWord[len(input)-1] != TRUE THEN
	# if cant find a chain of words till the last index, find the latest index for which. isWord = TRUE
	FOR each index j in input
		IF isWord[j] = TRUE THEN
			endIndex = j
		ENDIF
	ENDFOR
ENDIF 

# Now find indices to be capitalized using isWord and prevWord
res = []
i = endIndex
WHILE prevWord[i] != -1
	res.push(i)
	i = prevword[i]
ENDWHILE

#Capitalize letters at indices in res array
FOR index i in res
	input = input[:i] + i.upper() + input[i:]
ENDFOR
```

This algorithm runs in O(n^2). We will discuss improvements over the algorithm in Improvements section. 

A peculiar thing to observe about this algorithm is, it picks a chain of shortest words. 
For example, a string _mango_, could be the noun _mango_ or a combination of two words, _man_ and _go_. The algorithm displays the output as _manGo_. This is because the inner loop of the DP traverses the input in reverse and tries to find the shortest word from index j to i, that can be included in a chain of words till j.

## Implementation

[oxforddict.go](https://github.com/mugdhabondre/goCamelCase/blob/master/oxforddict/oxforddict.go) implements the _ConnectAndCheck_ method which connects to the oxford dictionaries API's "entries" endpoint and queries for the given word. The. response is checked to see if the word exists. I have used fields and filters to reduce the response size.

[gocamelcaseimpl.go](https://github.com/mugdhabondre/goCamelCase/blob/master/gocamelcaseimpl/gocamelcaseimpl.go) implements the actual core algorithm explained above. 

[httphandler.go](https://github.com/mugdhabondre/goCamelCase/blob/master/httphandler.go) is a wrapper function to expose RESTful apis. It invokes gocamelcaseimpl's _ProcessPhrase_ method to return camelCased version of an input sent in along with the URL as a field.

All information about supporting build files and credentials files, is explained in [InstallAndGo.md](https://github.com/mugdhabondre/goCamelCase/blob/master/InstallAndGo.md).

## Latency Analysis

Lets do an approximate latency analysis.
I used Google Chrome's dev tools to record time taken to complete various requests. The reponse time increases quadratically (as expected with O(n^2)) with the length of the input. 
Here's the average latency for different lengths of words, tested over \~10 API calls. 

| len(input) | Latency(s)  |
| -----------| -----------:|
| 5	         |        1.27s|
| 10         |       4.375s|
| 14         |        7.14s|

These latencies are way more than the acceptable limit of a normal service. Lets see reasons and improvements in the next section.

## Improvements

1. Algorithmic improvements<br/><br/>
__Bottom__ __Up__ __Approach__: The algorithm is bottom up currently. If we make it top down, we can save some time, by checking only for those indices which can form a potential word. Worst case would still be the same as bottom up, but we do better on the average case.<br/><br/>
__Reverse__ __Traversal__:  Another improvement could be start the second scanning of indices (j) in forward manner instead of reverse. I tried this, but it increases the latency by atleast 1s for sentences with 13-14 letters.<br/><br/>
__Multiple__ __Queries__ __from__ __a__ __single__ __connection__: The algorithm is quadratic. At each index (i) it can fire upto i queries, n-1 worst case. This is a lot of calls to the oxford dictionary API. One of the solutions could be to open an client connection and fire queries repeatedly to save time on handshakes. 
This can result in a few extra queries, but we would save some time by this. <br/><br/>

2. Oxford Dictionaries API
__Input__ __Array__ __of__ __Words__: The API does not have an endpoint which can take in an array of words and return if they are valid or not. This type of a request can save a lot of time <br/><br/>
__API__ __Response__: The "entries" endpoint returns a lot of data even after limiting the data using fields and filters.<br/><br/>
__Another__ __API?__: Many other APIs are avialable online which provide validation of words. They should be checked if they are any faster and provide the endpoints that we need.<br/><br/>

3. Use of higher programming primitives<br/><br/>
__Parallel__ __requests__: The queries going to the oxford dictionaries api can be multithreaded to achieve the same goal in lesser amount of time. We can parallelize the queries in the inner loop of the algorithm.  

4. Caching results<br/><br/>
We can cache results for words that were queried. This can help reduce latency if we have multiple incoming requests. 

5. Input structure<br/><br/>
We can reduce the number of calls going to the oxford dictionaries API by defining an input strcture which helps the user define the input with separators and then we can capitalize it. Most of the camel case generators online use this trick. This wont be ideal for the problem, just a suggestion.

