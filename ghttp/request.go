package ghttp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"
	"github.com/turing-era/turingera-shared/cutils"
	"github.com/turing-era/turingera-shared/log"
)

var client *http.Client

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
}

// Get 发送Get请求
func Get(path string, rsp interface{}, header map[string]string) error {
	httpReq, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			httpReq.Header.Set(k, v)
		}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debugf("http rsp: %s", body)
	if err = jsoniter.Unmarshal(body, rsp); err != nil {
		return err
	}
	return nil
}

// Post 发送post请求
func Post(path string, req interface{}, rsp interface{}, header map[string]string) error {
	mJson, err := jsoniter.Marshal(req)
	if err != nil {
		return err
	}
	log.Debugf("http req: %s", mJson)
	httpReq, err := http.NewRequest("POST", path, bytes.NewReader(mJson))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			httpReq.Header.Set(k, v)
		}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debugf("http rsp: %s", body)
	if err = jsoniter.Unmarshal(body, rsp); err != nil {
		return err
	}
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
