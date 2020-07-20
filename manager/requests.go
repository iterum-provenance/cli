package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func _request(method string, url *url.URL, headers map[string]string, body io.Reader) (resp *http.Response, err error) {
	client := http.Client{}
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return resp, err
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	resp, err = client.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return resp, fmt.Errorf("Request returned status code %v", resp.StatusCode)
	}
	return resp, nil
}

// delete performs a simple delete request
func delete(url *url.URL) (err error) {
	_, err = _request("DELETE", url, nil, nil)
	return err
}

// getJSON performs a simple get request, expecting a JSON response and decodes it into target which should be a pointer
func getJSON(url *url.URL, target interface{}) (err error) {
	resp, err := _request("GET", url, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// postJSON performs a simple post request, posting JSON and also expecting a JSON response and decodes it into target which should be a pointer
func postJSON(url *url.URL, body io.Reader, target interface{}) (err error) {
	resp, err := _request("POST", url, map[string]string{"content-type": "application/json"}, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
