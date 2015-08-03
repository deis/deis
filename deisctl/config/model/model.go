package model

import "time"

// ConfigNode represents running Deis services
type ConfigNode struct {
	Key        string     `json:"key"`
	Value      string     `json:"value,omitempty"`
	Expiration *time.Time `json:"expiration,omitempty"`
	TTL        int64      `json:"ttl,omitempty"`
}
