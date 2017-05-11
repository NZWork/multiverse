package data

import (
	"encoding/json"
	"errors"
	"log"
)

// TokenResponse 结构
type TokenResponse struct {
	Code   int    `json:"code"`
	Result Result `json:"result"`
	Msg    string `json:"msg"`
}

// Result 结构
type Result struct {
	Token     string `json:"token"`
	ProjectID string `json:"project_id"`
	UID       int    `json:"uid"`
	Time      int    `json:"time"`
}

// APIDomain API 域名
const APIDomain = "https://tiki.im/api/nz/"

// XAuth 认证
const XAuth = "5DE0CB6960FDD55B9F7C26E6554617B5"

// DeserializeToken 反序列化 Token
func DeserializeToken(token, pubkey string) (r Result, err error) {
	data, err := post(APIDomain+"identity", map[string]interface{}{
		"xauth":  XAuth,
		"xtoken": token,
		"pubkey": pubkey,
	})
	if err != nil {
		log.Println(err)
	}

	t := &TokenResponse{}
	json.Unmarshal(data, t)
	if t.Code == 2404 {
		err = errors.New(t.Msg)
		return
	}
	r = t.Result
	return
}

// PersistenTiki Tiki 持久化
func PersistenTiki(projectID, fileToken string) {
	var data []byte
	var err error
	data, err = post(APIDomain+"save", map[string]interface{}{
		"xauth": XAuth,
		"token": fileToken,
		"pid":   projectID,
	})

	if err != nil {
		log.Printf("persistence request api error: %s", err.Error())
	}

	t := &TokenResponse{}
	json.Unmarshal(data, t)
	if t.Code != 2201 {
		log.Printf("persistence error: [%d]%s", t.Code, t.Msg)
	}
}
