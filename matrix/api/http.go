package api

import (
	"net/http"
	"time"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"fmt"
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

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
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
