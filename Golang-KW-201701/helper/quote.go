package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func GetQuote() (string, error) {
	resp, err := http.Get("http://quotes.stormconsultancy.co.uk/random.json")
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unexpected status %s", resp.Status)
	}
	defer resp.Body.Close()
	var q struct {
		Quote string
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&q)
	if err != nil {
		return "", err
	}
	return q.Quote, nil
}

func ContainsInappropriateLanguage(text string) (bool, error) {
	resp, err := http.Get("http://www.purgomalum.com/service/containsprofanity?add=sex&text=" +
		url.QueryEscape(text))
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Unexpected status %s", resp.Status)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(string(data))
}
