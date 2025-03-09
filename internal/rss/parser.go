package internal

import (
	"encoding/xml"
	"fmt"
)

func parseXML(data []byte) (*RSSFeed, error) {
	var feed RSSFeed

	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("Error unmarshaling xml data: %v", err)
	}

	return &feed, nil
}
