package model

// RecordID defines a record id. Together with RecordType identifies unique records across all types
type RecordID string

// RecordType defines a record type. Together with RecordID identifies unique records across all types.
type RecordType string

// RecordTypeMovie existing record types.
const (
	RecordTypeMovie = RecordType("movie")
)

// UserID defines a user id
type UserID string

type RatingValue int

// Rating defines an individual rating created by a user some record.
type Rating struct {
	RecordID   string      `json:"recordID"`
	RecordType string      `json:"recordType"`
	UserID     UserID      `json:"userID"`
	Value      RatingValue `json:"value"`
}
