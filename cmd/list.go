package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered subscriptions",
	Long:  `List all registered subscriptions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := requestAccessToken(apiBaseURLFlag, clientIDFlag, clientSecretFlag)
		if err != nil {
			return err
		}

		subsJson, err := listSubscriptions(t)
		if err != nil {
			return err
		}

		var subs []Subscription
		err = json.Unmarshal(subsJson, &subs)
		if err != nil {
			return err
		}

		maxNameLength := 1
		for _, s := range subs {
			if len(s.Name) > maxNameLength {
				maxNameLength = len(s.Name)
			}
		}

		header := fmt.Sprintf("%s%s%s%s%s%s",
			strings.Repeat(" ", 18), "ID",
			strings.Repeat(" ", 18+maxNameLength/2), "Name",
			strings.Repeat(" ", maxNameLength/2), "Description")
		fmt.Printf(" %s\n", header)
		fmt.Printf(" %s\n", strings.Repeat("=", len(header)+18))

		for _, s := range subs {
			fmt.Printf(" %s   %s   %s\n", s.ID, s.Name, s.Description)
		}

		return nil
	},
}

func listSubscriptions(accessToken string) ([]byte, error) {
	URL := wsBaseURLFlag
	URL = URL + "/subscription"
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
