package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
)

func requestAccessToken(baseURL string, clientID string, clientSecret string) (string, error) {
	URL := baseURL + "/oauth/access_token"
	form := url.Values{}
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", URL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("")
	}

	var g struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	err = json.Unmarshal(respBody, &g)

	return g.AccessToken, err
}

func readSubscriptionFile(fileName string) (Subscription, error) {
	var sub Subscription

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return sub, err
	}

	err = json.Unmarshal(b, &sub)

	return sub, err
}

func printJson(msg []byte) error {
	var s []byte
	var v interface{}
	var o map[string]interface{}
	var a []map[string]interface{}

	if bytes.HasPrefix(msg, []byte("[")) {
		err := json.Unmarshal(msg, &a)
		if err != nil {
			return err
		}

		v = a
	} else {
		err := json.Unmarshal(msg, &o)
		if err != nil {
			return err
		}

		v = o
	}

	s, err := coloredPrettyPrint(v)
	if err != nil {
		return err
	}

	fmt.Println(string(s))

	return nil
}

func coloredPrettyPrint(v interface{}) ([]byte, error) {
	s, err := prettyjson.Marshal(v)
	if err != nil {
		return nil, err
	}

	return s, nil
}
