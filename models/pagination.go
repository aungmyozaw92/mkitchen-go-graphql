package models

import (
	"encoding/base64"
)

type PageInfo struct {
	StartCursor string `json:"startCursor"`
	EndCursor   string `json:"endCursor"`
	HasNextPage *bool  `json:"hasNextPage,omitempty"`
}

func DecodeCursor(cursor *string) (string, error) {
	decodedCursor := ""
	if cursor != nil {
		b, err := base64.StdEncoding.DecodeString(*cursor)
		if err != nil {
			return decodedCursor, err
		}
		decodedCursor = string(b)
	}
	return decodedCursor, nil
}

func EncodeCursor(cursor string) string {
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}
