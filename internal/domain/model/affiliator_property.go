//go:generate go run github.com/terui-ryota/gen-getter -type=AffiliatorProperty

package model

// アフィリエイター情報
// データの管理元は affiliator コンテキスト
// TODO: 実装が追いつき次第 nolint:unusedを削除
// nolint:unused
type AffiliatorProperty struct {
	// アフィリエイターID
	affiliatorID AffiliatorID
	// メディアユーザーID（ASIDを想定）
	mediaUserID string
	// メディアURL（アメーバブログURLを想定）
	mediaURL string
	// アカウントの有効性
	validity AffiliatorValidity
}

// アフィリエイターID
type AffiliatorID string

func (ai AffiliatorID) String() string {
	return string(ai)
}

// アフィリエイターIDリスト
type AffiliatorIDList []AffiliatorID

// アカウントの有効性
type AffiliatorValidity int

const (
	AffiliatorValidityUnknown   AffiliatorValidity = iota // 不明
	AffiliatorValidityValid                               // 有効
	AffiliatorValiditySuspended                           // アカウント停止中
	AffiliatorValidityBanned                              // 利用禁止
	AffiliatorValidityDeleted                             // アカウント削除済み
)

func (av AffiliatorValidity) Int() int {
	return int(av)
}
