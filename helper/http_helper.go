package helper

import (
	"net"
	"net/http"
)

/***********************************************************************************************
 * Description:容器管理平台(CMP) -> API http访问客户端
 *             中润四方版权所有
 * Author: BaoQL
 * Date:   2015-09-29
 **********************************************************************************************/

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type (
	TCommonResult struct {
		Result int
		Error  string
		Data   interface{}
	}
)

func GetHttpClient(sUrl string) *http.Client {
	tr := &http.Transport{
		DisableKeepAlives: true,
		Dial: func(netw, addr string) (net.Conn, error) {
			deadline := time.Now().Add(time.Duration(30) * time.Second)
			c, err := net.DialTimeout(netw, addr, time.Duration(30)*time.Second)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		},
	}

	nPos := strings.Index(sUrl, "https://")
	if nPos == 0 {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		tr.DisableCompression = true
	}

	return &http.Client{Transport: tr}
}

func CloseRespose(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

func ApiGet(sUrl string, resp *TCommonResult) error {
	WriteLog(0, "GET请求到："+sUrl)
	httpClient := GetHttpClient(sUrl)
	reqest, err := http.NewRequest("GET", sUrl, nil)
	reqest.Header.Set("Connection", "close")
	res, err := httpClient.Do(reqest)
	defer CloseRespose(res)

	if err != nil {
		WriteLog(2, "http post err:"+err.Error())
		return err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		WriteLog(2, "http read err:"+err.Error())
		return err
	}

	//logger.Debug("返回结果：%s\r\n" + string(result))
	if err = json.Unmarshal(result, resp); err != nil {
		WriteLog(2, "parse result err:"+err.Error())
		return err
	}

	if resp.Result != 0 {
		return errors.New(resp.Error)
	}

	return nil
}

func ApiPost(sUrl string, req interface{}, resp *TCommonResult) error {
	var body *bytes.Buffer
	if req != nil {
		b, err := json.Marshal(req)
		if err != nil {
			WriteLog(2, "api client, http json err:"+err.Error())
			return err
		}
		body = bytes.NewBuffer(b)
	} else {
		body = bytes.NewBuffer([]byte(""))
	}
	WriteLog(0, "POST请求到："+sUrl+"，字节长度："+strconv.Itoa(body.Len()))

	httpClient := GetHttpClient(sUrl)
	reqest, err := http.NewRequest("POST", sUrl, body)
	reqest.Header.Set("Connection", "close")
	res, err := httpClient.Do(reqest)
	defer CloseRespose(res)

	if err != nil {
		WriteLog(2, "POST ERROR："+sUrl+"，字节长度："+strconv.Itoa(body.Len())+", "+err.Error())
		return err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		WriteLog(2, "http read err:"+err.Error())
		return err
	}

	//logger.Debug("返回结果：%s\r\n" + string(result))
	if err = json.Unmarshal(result, resp); err != nil {
		WriteLog(2, "parse result err:"+err.Error())
		return err
	}

	if resp.Result != 0 {
		return errors.New(resp.Error)
	}

	return nil
}

func ApiDelete(sUrl string, req interface{}, resp *TCommonResult) error {
	var body *bytes.Buffer
	if req != nil {
		b, err := json.Marshal(req)
		if err != nil {
			WriteLog(2, "api client, http json err:"+err.Error())
			return err
		}
		body = bytes.NewBuffer(b)
	} else {
		body = bytes.NewBuffer([]byte(""))
	}
	WriteLog(0, "DELETE请求到："+sUrl+"，字节长度："+strconv.Itoa(body.Len()))

	httpClient := GetHttpClient(sUrl)
	reqest, err := http.NewRequest("DELETE", sUrl, body)
	reqest.Header.Set("Connection", "close")
	res, err := httpClient.Do(reqest)
	defer CloseRespose(res)

	if err != nil {
		WriteLog(2, "http delete err:"+err.Error())
		return err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		WriteLog(2, "http read err:"+err.Error())
		return err
	}

	//WriteLog(0, "返回结果：%s\r\n" + string(result))
	if err = json.Unmarshal(result, resp); err != nil {
		WriteLog(2, "parse result err:"+err.Error())
		return err
	}

	if resp.Result != 0 {
		return errors.New(resp.Error)
	}

	return nil
}

func ApiRequest(sUrl, method string) (*http.Response, error) {
	httpClient := GetHttpClient(sUrl)
	reqest, _ := http.NewRequest(method, sUrl, nil)
	reqest.Header.Set("Connection", "close")
	res, err := httpClient.Do(reqest)
	return res, err
}

func TranGet(sUrl string) ([]byte, error) {
	WriteLog(0, "发送请求到："+sUrl)
	httpClient := GetHttpClient(sUrl)
	reqest, err := http.NewRequest("GET", sUrl, nil)
	reqest.Header.Set("Connection", "close")
	res, err := httpClient.Do(reqest)
	defer CloseRespose(res)

	if err != nil {
		WriteLog(2, "http post err:"+err.Error())
		return []byte(""), err
	}

	if res.StatusCode != 200 {
		return []byte(""), errors.New(strconv.Itoa(res.StatusCode))
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		WriteLog(2, "http read err:"+err.Error())
		return []byte(""), err
	}

	//logger.Debug("返回结果：%s\r\n" + string(result))
	return result, nil
}

func TranPost(sUrl string, req []byte) ([]byte, error) {
	var body *bytes.Buffer
	if req != nil {
		WriteLog(0, "发送请求到："+sUrl)
		body = bytes.NewBuffer(req)
	} else {
		WriteLog(0, "发送请求到："+sUrl)
		body = bytes.NewBuffer([]byte(""))
	}

	httpClient := GetHttpClient(sUrl)
	reqest, err := http.NewRequest("POST", sUrl, body)
	reqest.Header.Set("Connection", "close")
	res, err := httpClient.Do(reqest)
	defer CloseRespose(res)

	if err != nil {
		WriteLog(2, "http post err:"+err.Error())
		return []byte(""), err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		WriteLog(2, "http read err:"+err.Error())
		return []byte(""), err
	}

	//WriteLog(0, "返回结果：%s\r\n"+string(result))
	return result, nil
}

func WriteLog(level int, Msg string) {

	if level == 0 {
		log.Debug(Msg)
	} else if level == 1 {
		log.Info(Msg)
	} else {
		log.Error(Msg)
	}

}
