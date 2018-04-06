package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch a subscription spec",
	Long:  `Fetch the subscription spec for the given id or name`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := requestAccessToken(apiBaseURLFlag, clientIDFlag, clientSecretFlag)
		if err != nil {
			return err
		}

		subJson, err := getSubscription(t, args[0])
		if err != nil {
			return err
		}

		printJson(subJson)

		return nil
	},
}

func getSubscription(accessToken string, subID string) ([]byte, error) {
	URL := wsBaseURLFlag
	URL = URL + "/subscription/" + subID
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
		return nil, fmt.Errorf("Unexpected status code: %d. Body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, err
}
