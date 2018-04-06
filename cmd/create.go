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

var subInputFileNameFlag string
var subNameFlag string

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&subInputFileNameFlag, "file", "f", "", "Input file containing subscription specification (required)")
	createCmd.Flags().StringVarP(&subNameFlag, "name", "n", "", "An optional name to register the subscription with")

	createCmd.MarkFlagRequired("file")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new subscription",
	Long:  `Create and register a subscription it with the push service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sub, err := readSubscriptionFile(subInputFileNameFlag)
		if err != nil {
			return err
		}

		// If the name flag was used it will override the name set
		// in the file (if any)
		if subNameFlag != "" {
			sub.Name = subNameFlag
		}

		t, err := requestAccessToken(apiBaseURLFlag, clientIDFlag, clientSecretFlag)
		if err != nil {
			return err
		}

		id, alreadyExisted, err := registerSubscription(t, sub)
		if err != nil {
			return err
		}

		if alreadyExisted {
			return fmt.Errorf("ERROR: Subscription with name '%s' already exists. id=%s\n", sub.Name, id.String())
		}

		// Everything ok, subscription was registered with the service
		fmt.Printf("Subscription with name '%s' created. id=%s\n", sub.Name, id.String())

		return nil
	},
}

func registerSubscription(accessToken string, sub Subscription) (uuid.UUID, bool, error) {
	URL := wsBaseURLFlag
	URL = URL + "/subscription"
	URL = URL + "?access_token=" + accessToken

	j, _ := json.Marshal(sub)

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(j))
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

	// The subscription POST endpoint response have 2 normal status codes:
	//  * Unprocessable Entity (422)
	//    This is returned by the server if client tries to register a subscription
	//    with a name that has already been registered on the server.
	//  * OK (200)
	//    If the registration was successful
	if resp.StatusCode == http.StatusUnprocessableEntity {
		var existingID uuid.UUID

		// If we get HTTP response code 422 the server has also set
		// the 'Location' header with the ID of the existing subscription
		if resp.Header.Get("Location") != "" {
			existingID, err = uuid.FromString(resp.Header.Get("Location"))
			if err != nil {
				// Location header didn't contain a valid UUID
				return uuid.Nil, true, err
			}

			return existingID, true, nil
		}

		// Server didn't set a valid ID in the 'Location' header, this should never happen
		return uuid.Nil, true, fmt.Errorf("Subscription with name already exists, but failed to retrieve ID")
	} else if resp.StatusCode != http.StatusOK {
		return uuid.Nil, false, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var s struct {
		ID uuid.UUID `json:"id"`
	}
	err = json.Unmarshal(respBody, &s)

	return s.ID, false, err

}
