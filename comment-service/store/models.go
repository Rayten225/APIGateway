package store

import (
    "time"
)

type Comment struct {
    ID       int       `json:"id"`
    NewsID   int       `json:"news_id"`
    ParentID *int      `json:"parent_id,omitempty"`
    Text     string    `json:"text"`
    Created  time.Time `json:"created"`
}
