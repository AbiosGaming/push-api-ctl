package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show account config",
	Long:  `Show account config for the push service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := requestAccessToken(apiBaseURLFlag, clientIDFlag, clientSecretFlag)
		if err != nil {
			return err
		}

		configJson, err := fetchPushServiceConfig(t)
		if err != nil {
			return err
		}

		printJson(configJson)

		return nil
	},
}

func fetchPushServiceConfig(accessToken string) ([]byte, error) {
	URL := wsBaseURLFlag
	URL = URL + "/config"
	URL = URL + "?access_token=" + accessToken

	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	return respBody, err
}
