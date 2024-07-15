package model

// アメーバID
type AmebaID string

func (ai AmebaID) String() string {
	return string(ai)
}
