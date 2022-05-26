// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"encoding/json"
	"fmt"
)

type decisions map[string]interface{}

func (d *decisions) update(u *decisions) {
	for key, value := range *u {
		(*d)[key] = value
	}
}

func (d *decisions) encode() (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", fmt.Errorf("encode: error marshaling: %w", err)
	}
	return string(b), nil
}

func parseDecisions(s string) (*decisions, error) {
	d := decisions{}
	if s != "" {
		if err := json.Unmarshal([]byte(s), &d); err != nil {
			return nil, fmt.Errorf("parseDecisions: error unmarshaling: %w", err)
		}
	}
	return &d, nil
}
