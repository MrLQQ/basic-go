package demo

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 问题：我怎么在序列化 JSON 的时候，加密一些字段？

type User struct {
	Password MaskedString
	Email    string
	Phone    MaskedPhone
}

type MaskedPhone string

func (m MaskedPhone) MarshalJSON() ([]byte, error) {
	str := string(m)
	var res = []byte(str[0:3])
	res = append(res, "****"...)
	res = append(res, str[7:]...)
	return res, nil
}

type MaskedString string

func (m MaskedString) MarshalJSON() ([]byte, error) {
	return []byte(`"****"`), nil
}

func TestJson(t *testing.T) {
	u := User{
		Password: "123456",
		Email:    "123@qq.com",
	}
	val, err := json.Marshal(u)
	assert.NoError(t, err)
	t.Log(string(val))
}
