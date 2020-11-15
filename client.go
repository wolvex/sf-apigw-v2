package apigw

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Client struct {
	BaseURL   string
	SecretKey string
	KeyID     string
	Session   *http.Client
}

func New(baseUrl, keyId, secretKey string, timeout int) *Client {
	api := &Client{
		BaseURL:   baseUrl,
		SecretKey: secretKey,
		KeyID:     keyId,
	}

	api.Session = &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	api.SetTransport("")

	return api
}

func (ws *Client) SetProxy(uri string) error {
	ws.SetTransport(uri)
	return nil
}

func (ws *Client) SetTransport(proxyUri string) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
			//KeepAlive: 10 * time.Second,
			//DualStack: true,
		}).DialContext,
		//ForceAttemptHTTP2:     true,
		MaxIdleConns: 10,
		//MaxIdleConnsPerHost:   10,
		IdleConnTimeout:     10 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
		//ExpectContinueTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if proxyUri != "" {
		if proxy, err := url.Parse(proxyUri); err == nil {
			transport.Proxy = http.ProxyURL(proxy)
		}
	}

	ws.Session.Transport = transport
}

func (ws *Client) Post(uri, version string, request []byte) (response []byte, err error) {
	if uri == "" {
		err = fmt.Errorf("Unable to resolve uri")
		return
	}

	path := ws.BaseURL + uri

	req, err := http.NewRequest("POST", path, bytes.NewBufferString(string(request[:])))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	response, err = ws.SubmitRequest(req, version)

	return
}

func (ws *Client) Get(uri, version string, request map[string]string) (response []byte, err error) {
	if uri == "" {
		err = fmt.Errorf("Unable to resolve uri")
		return
	}

	var base *url.URL
	if base, err = url.Parse(ws.BaseURL + uri); err != nil {
		return
	}

	// Query params
	params := url.Values{}
	for k, v := range request {
		params.Add(k, v)
	}
	base.RawQuery = params.Encode()

	//initialize GET request
	req, err := http.NewRequest("GET", base.String(), nil)
	if err != nil {
		return
	}

	response, err = ws.SubmitRequest(req, version)

	return
}

func (ws *Client) SubmitRequest(req *http.Request, version string) (response []byte, err error) {
	//add necessary headers
	req.Header.Add("Date", generateDateHeader())
	signature := generateHMACSignature(ws.KeyID, ws.SecretKey, req.Header.Get("Date"))
	if signature != "" {
		req.Header.Add("Authorization", signature)
	}
	if version != "" {
		req.Header.Add("X-Version", version)
	}

	//dump it for logging
	requestDump, _ := httputil.DumpRequest(req, true)
	fmt.Printf("HTTP Request: %q\n", requestDump)

	//fire !!!
	var res *http.Response
	res, err = ws.Session.Do(req)
	if err != nil {
		fmt.Printf("Exception caught %#v\n", err)
		return
	}
	defer res.Body.Close()

	//dump response for logging
	if responseDump, e := httputil.DumpResponse(res, true); e == nil {
		fmt.Printf("HTTP Response: %q\n", responseDump)
	}

	//get body
	if response, err = ioutil.ReadAll(res.Body); err != nil {
		fmt.Printf("Exception caught %#v\n", err)
	}

	return
}

func generateDateHeader() string {
	loc, _ := time.LoadLocation("MST")
	return time.Now().In(loc).Format("Mon, 02 Jan 2006 15:04:05 MST")
}

func generateHMACSignature(keyId, secretKey, date string) string {
	if keyId == "" || secretKey == "" || date == "" {
		return ""
	}

	// Prepare the signature to include those headers:
	//data := "(request-target): post " + uri + "\n"
	data := "date: " + date
	//fmt.Printf("SECRET: [%s] , DATA: [%s]\n", ws.SecretKey, data)

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))

	// Base64 and URL Encode the string
	sigString := base64.StdEncoding.EncodeToString(h.Sum(nil))
	encodedString := url.QueryEscape(sigString)

	//fmt.Printf("BASE64: [%s] , ESCAPE: [%s]\n", sigString, encodedString)

	//req.Header.Add("Authorization", "Signature keyid="+ws.KeyID+",algorithm=hmac-sha256,headers=date,signature="+encodedString)
	return "Signature keyid=\"" + keyId + "\",algorithm=\"hmac-sha256\",signature=\"" + encodedString + "\""
}
