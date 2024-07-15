package id

import (
	"math/rand"
	"time"

	"github.com/eknkc/basex"
	"github.com/google/uuid"
)

// IDを半角数字+半角英字(大文字小文字)の62文字で表現
const base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// UUIDをベースに変換
const idlen = 22

var (
	baselen   int
	enc       *basex.Encoding
	baseChars []byte
	random    *rand.Rand
)

func init() {
	baselen = len(base62)
	enc, _ = basex.NewEncoding(base62)
	baseChars = []byte(base62)
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 新しいIDを生成します。
func New() string {
	// gen UUID (ver4)
	uuid := uuid.New()
	bytes, _ := uuid.MarshalBinary()

	id := enc.Encode(bytes)
	idBytes := []byte(id)
	genlen := len(idBytes)

	buf := make([]byte, idlen)

	// バイト列を逆順に詰める
	for i := 0; i < genlen; i++ {
		buf[idlen-i-1] = idBytes[i]
	}

	// 指定文字数に達しない場合はランダムに文字をくっつける
	for i := idlen - genlen; i > 0; i-- {
		buf[i-1] = baseChars[random.Intn(baselen)]
	}

	return string(buf)
}
