package store

import (
    "time"
)

type News struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Published time.Time `json:"published"`
}
