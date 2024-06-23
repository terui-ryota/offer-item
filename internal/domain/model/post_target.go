package model

// スケジュールタイプ
type PostTarget int

const (
	PostTargetTypeUnknown PostTarget = iota // 不明
	PostTargetAmeba                         // Ameba
	PostTargetX                             // X
	PostTargetInstagram                     // Instagram
)
