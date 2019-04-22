package oxfordDict

import (
	"strings"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
)

func Connect(word_id string) bool{
	app_id := "	0043c514"
	app_key := "6039d61fd17213e07123cc1f79b2205b"

	endpoint := "entries"
	language_code := "en-us"

	categories :=  "noun,verb,adjective,adverb,conjunction,numeral,particle,preposition,pronoun"
	url := "https://od-api.oxforddictionaries.com/api/v2/" + endpoint + "/" + language_code + "/" + strings.ToLower(word_id) + "?lexicalCategory=" + categories

	client := &http.Client{}
	//resp, err := http.Get(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		os.Exit(1)
	}
	req.Header.Add("app_id", app_id)
	req.Header.Add("app_key", app_key)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	responseString := string(body)

	var jsonResp map[string]interface{}
	json.Unmarshal([]byte(responseString), &jsonResp)

	_, ok := jsonResp["error"]

	if ok == false {

	}

	return !ok

}
