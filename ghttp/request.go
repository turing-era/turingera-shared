package ghttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"
	"github.com/turing-era/turingera-shared/cutils"
	"github.com/turing-era/turingera-shared/log"
)

// Post 发送post请求
func Post(path string, req interface{}, rsp interface{}) error {
	mJson, err := jsoniter.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Post(path, "application/json", bytes.NewReader(mJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("rsp: %s\n", body)
	if err = jsoniter.Unmarshal(body, rsp); err != nil {
		return err
	}
	log.Debugf("http request: %v, req: %s, rsp: %s", path, mJson, cutils.Obj2Json(rsp))
	return nil
}

// PostForm 发送post表单请求
func PostForm(path string, req interface{}, rsp interface{}) error {
	reqFrom := url.Values{}
	mJson, err := jsoniter.Marshal(req)
	if err != nil {
		return err
	}
	reqMap := make(map[string]interface{})
	if err = jsoniter.Unmarshal(mJson, &reqMap); err != nil {
		return err
	}
	for k, v := range reqMap {
		reqFrom.Set(k, fmt.Sprintf("%s", v))
	}
	log.Debugf("reqFrom: %v", reqFrom)
	resp, err := http.PostForm(path, reqFrom)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = jsoniter.Unmarshal(body, rsp); err != nil {
		return err
	}
	log.Debugf("http request: %v, req: %s, rsp: %s", path, mJson, cutils.Obj2Json(rsp))
	return nil
}
