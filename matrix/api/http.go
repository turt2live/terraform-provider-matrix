package api

import (
	"net/http"
	"time"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"fmt"
	"io"
)

var matrixHttpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// Based in part on https://github.com/matrix-org/gomatrix/blob/072b39f7fa6b40257b4eead8c958d71985c28bdd/client.go#L180-L243
func DoRequest(method string, urlStr string, body interface{}, result interface{}, accessToken string) (error) {
	var bodyBytes []byte
	if body != nil {
		jsonStr, err := json.Marshal(body)
		if err != nil {
			return err
		}

		bodyBytes = jsonStr
	}

	return doRawRequest(method, urlStr, bodyBytes, "application/json", result, accessToken)
}

func UploadFile(csApiUrl string, content []byte, name string, mime string, accessToken string) (*ContentUploadResponse, error) {
	qs := make(map[string]string)
	if name != "" {
		qs["filename"] = name
	}
	urlStr := MakeUrlQueryString(qs, csApiUrl, "/_matrix/media/r0/upload")
	result := &ContentUploadResponse{}
	err := doRawRequest("POST", urlStr, content, mime, result, accessToken)
	return result, err
}

func DownloadFile(csApiUrl string, origin string, mediaId string) (*io.ReadCloser, http.Header, error) {
	urlStr := MakeUrl(csApiUrl, "/_matrix/media/r0/download", origin, mediaId)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := matrixHttpClient.Do(req)
	if err != nil {
		return &res.Body, res.Header, err
	}

	if res.StatusCode != http.StatusOK {
		return &res.Body, res.Header, fmt.Errorf("request failed: status code %d", res.StatusCode)
	}

	return &res.Body, res.Header, nil
}

func doRawRequest(method string, urlStr string, bodyBytes []byte, contentType string, result interface{}, accessToken string) (error) {
	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	res, err := matrixHttpClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return err
	}

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		mtxErr := &ErrorResponse{}
		mtxErr.RawError = string(contents)
		mtxErr.StatusCode = res.StatusCode
		err = json.Unmarshal(contents, mtxErr)
		if err != nil {
			return fmt.Errorf("request failed: %s", string(contents))
		}
		return mtxErr
	}

	if result != nil {
		err = json.Unmarshal(contents, &result)
		if err != nil {
			return err
		}
	}

	return nil
}

func MakeUrl(parts ... string) string {
	res := ""
	for i, p := range parts {
		if p[len(p)-1:] == "/" {
			res += p[:len(p)-1]
		} else if p[0] != '/' && i > 0 {
			res += "/" + p
		} else {
			res += p
		}
	}
	return res
}

func MakeUrlQueryString(query map[string]string, parts ... string) string {
	urlStr := MakeUrl(parts...)

	u, _ := url.Parse(urlStr)
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
