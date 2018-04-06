package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete subscriptions",
	Long:  `Deletes the subscription(s) with the given identifier(s) (name or UUID) from the push service`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		accessToken, err := requestAccessToken(apiBaseURLFlag, clientIDFlag, clientSecretFlag)
		if err != nil {
			return err
		}

		err = deleteSubscription(accessToken, args...)

		return err
	},
}

func deleteSubscription(accessToken string, subscriptionIDOrNameSlice ...string) error {
	for _, subscriptionIDOrName := range subscriptionIDOrNameSlice {
		URL := wsBaseURLFlag
		URL = URL + "/subscription/" + subscriptionIDOrName
		URL = URL + "?access_token=" + accessToken

		req, err := http.NewRequest("DELETE", URL, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	}
	return nil
}
