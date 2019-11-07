package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

//HTTPGet get 请求
func HTTPGet(uri string) ([]byte, error) {
	response, err := http.Get(uri)
	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//PostJSON post json 数据请求
func PostJSON(uri string, obj interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(uri, "application/json;charset=utf-8", body)
	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

// PostJSONWithRespContentType post json数据请求，且返回数据类型
func PostJSONWithRespContentType(uri string, obj interface{}) ([]byte, string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, "", err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(uri, "application/json;charset=utf-8", body)
	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()
	if err != nil {
		return nil, "", err
	}

	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	contentType := response.Header.Get("Content-Type")
	return responseData, contentType, err
}

//PostFile 上传文件
func PostFile(fieldname, filename, uri string) ([]byte, error) {
	fields := []MultipartFormField{
		{
			IsFile:    true,
			Fieldname: fieldname,
			Filename:  filename,
		},
	}
	return PostMultipartForm(fields, uri)
}

func PostFileV2(fieldname, filename string, fileReader io.Reader, uri string) (respBody []byte, err error) {
	field := MultipartFormField{
		IsFile:    true,
		Fieldname: fieldname,
		Filename:  filename,
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
	if e != nil {
		err = fmt.Errorf("error writing to buffer , err=%v", e)
		return
	}
	if _, err = io.Copy(fileWriter, fileReader); err != nil {
		return
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, e := http.Post(uri, contentType, bodyBuf)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if e != nil {
		err = e
		return
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	return
}

//MultipartFormField 保存文件或其他字段信息
type MultipartFormField struct {
	IsFile    bool
	Fieldname string
	Value     []byte
	Filename  string
}

//PostMultipartForm 上传文件或其他多个字段
func PostMultipartForm(fields []MultipartFormField, uri string) (respBody []byte, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for _, field := range fields {
		if field.IsFile {
			fileWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
			if e != nil {
				err = fmt.Errorf("error writing to buffer , err=%v", e)
				return
			}

			fh, e := os.Open(field.Filename)
			if e != nil {
				err = fmt.Errorf("error opening file , err=%v", e)
				return
			}
			defer fh.Close()

			if _, err = io.Copy(fileWriter, fh); err != nil {
				return
			}
		} else {
			partWriter, e := bodyWriter.CreateFormField(field.Fieldname)
			if e != nil {
				err = e
				return
			}
			valueReader := bytes.NewReader(field.Value)
			if _, err = io.Copy(partWriter, valueReader); err != nil {
				return
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, e := http.Post(uri, contentType, bodyBuf)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	return
}

//PostXML perform a HTTP/POST request with XML body
func PostXML(uri string, obj interface{}, xmlHead string, client *http.Client) ([]byte, error) {
	xmlData, err := xml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	_xmlData := string(xmlData)
	_xmlData = strings.Replace(_xmlData, xmlHead, "xml", 2)
	body := bytes.NewBuffer([]byte(_xmlData))
	if client == nil {
		client = http.DefaultClient
	}
	var response *http.Response
	response, err = client.Post(uri, "application/xml;charset=utf-8", body)
	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

// NewHTTPSClient 获取默认https客户端
func NewTLSHttpClientByFile(certFile, keyFile string) (httpClient *http.Client, err error) {
	certPEMBlock, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return NewTLSHttpClient(certPEMBlock, keyPEMBlock)
}

func NewTLSHttpClient(certPEMBlock, keyPEMBlock []byte) (httpClient *http.Client, err error) {
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			TLSClientConfig:     tlsConfig,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: 15 * time.Second,
	}
	return
}
