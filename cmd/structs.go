package cmd

import uuid "github.com/satori/go.uuid"

type Subscription struct {
	ID          uuid.UUID            `json:"id"`                    // Read-only, can't be set by the client when creating a subscription
	Description string               `json:"description,omitempty"` // Optional description of the subscription
	Name        string               `json:"name,omitempty"`        // Optional when creating a subscription
	Filters     []SubscriptionFilter `json:"filters"`
}

type SubscriptionFilter struct {
	Channel  string `json:"channel,omitempty"`
	GameID   int    `json:"game_id,omitempty"`
	SeriesID int    `json:"series_id,omitempty"`
	MatchID  int    `json:"match_id,omitempty"`
}
