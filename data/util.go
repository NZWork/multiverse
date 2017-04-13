package data

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func post(url string, params map[string]interface{}) ([]byte, error) {
	resp, err := http.PostForm(url, paramsFormator(params))
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return []byte(body), nil
}

// 构造参数 (URLEncode)
func paramsFormator(params map[string]interface{}) url.Values {
	v := url.Values{}
	for k, p := range params {
		v.Add(k, p.(string))
	}
	return v
}
