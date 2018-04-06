package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&subInputFileNameFlag, "file", "f", "", "Input file containing subscription specification (required)")

	updateCmd.MarkFlagRequired("file")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing subscription spec",
	Long:  `Update the subscription spec already registered with the push service`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sub, err := readSubscriptionFile(subInputFileNameFlag)
		if err != nil {
			return err
		}

		t, err := requestAccessToken(apiBaseURLFlag, clientIDFlag, clientSecretFlag)
		if err != nil {
			return err
		}

		subscriptionID := args[0]
		subID, exists, err := updateSubscription(t, subscriptionID, sub)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("ERROR: Could not update subscription: No subscription with id '%s' was registered on the server.\n", subscriptionID)
		}

		if subID.String() != subscriptionID {
			fmt.Printf("Updated subscription with id '%s'\n", subID.String())
		} else {
			fmt.Println("Updated subscription")
		}

		return nil
	},
}

func updateSubscription(accessToken string, subscriptionIDOrName string, sub Subscription) (uuid.UUID, bool, error) {
	URL := wsBaseURLFlag
	URL = URL + "/subscription/" + subscriptionIDOrName
	URL = URL + "?access_token=" + accessToken

	j, _ := json.Marshal(sub)

	req, err := http.NewRequest("PUT", URL, bytes.NewBuffer(j))
	if err != nil {
		return uuid.Nil, false, err
	}
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return uuid.Nil, false, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return uuid.Nil, false, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return uuid.Nil, false, nil
	} else if resp.StatusCode != http.StatusOK {
		return uuid.Nil, false, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var s struct {
		ID uuid.UUID `json:"id"`
	}
	err = json.Unmarshal(respBody, &s)

	return s.ID, true, err

}
