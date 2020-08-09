package data

import (
	"fmt"
	"github.com/go-playground/validator"
	"strings"
	"time"
)

const contentTimestampPrecision int64 = 1e9

// Subscription ID
type ID int64

type Subscription struct {
	ID       ID     `json:"id" bson:"_id"`
	URL      string `json:"url" bson:"url" validate:"required,url"`
	Interval int64  `json:"interval" bson:"interval" validate:"required,gt=0"` // Frequency of downloads for worker - in SECONDS
}

func (s *Subscription) Validate() error {
	v := validator.New()
	if err := v.Struct(s); err != nil {
		return err
	}
	return nil
}

type History []Content

// Represents single fetch operation. That is, result of single worker's request.
type Content struct {
	SubID     ID      `json:"-" bson:"subID"`           // id of corresponding subscription
	Response  *string `json:"response" bson:"response"` // We want to distinguish between empty response and "null-one"
	Duration  float64 `json:"duration" bson:"duration"`
	CreatedAt string  `json:"created_at" bson:"created_at"`
}

func (c *Content) DeepCopy() *Content {
	contCopy := *c
	if contCopy.Response != nil {
		response := *contCopy.Response
		contCopy.Response = &response
	}
	return &contCopy
}

func NewContent(subID ID, response *string, duration float64) *Content {
	return &Content{
		SubID:     subID,
		Response:  response,
		Duration:  duration,
		CreatedAt: getTimeFormatted(),
	}
}

// Returns current Unix timestamp formatted to our use case.
func getTimeFormatted() string {
	t := time.Now().UnixNano()
	ft := fmt.Sprintf("%d.%d", t/contentTimestampPrecision, t%contentTimestampPrecision)
	return strings.TrimRight(ft, "0")
}
