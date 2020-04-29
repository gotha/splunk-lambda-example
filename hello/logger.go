package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type SplunkLogger struct {
	token      string
	URL        string
	httpClient httpClient
	items      map[string]string
	err        error
}

func NewSplunkLogger(token, URL string, httpClient httpClient) *SplunkLogger {

	l := &SplunkLogger{
		token:      token,
		URL:        URL,
		httpClient: httpClient,
	}
	l.items = make(map[string]string, 0)
	return l
}

func (l *SplunkLogger) WithProperty(key, value string) *SplunkLogger {
	l.items[key] = value
	return l
}

func (l *SplunkLogger) WithError(err error) *SplunkLogger {
	l.err = err
	return l
}

func (l *SplunkLogger) WithServiceName(value string) *SplunkLogger {
	l.items["service_name"] = value
	return l
}
func (l *SplunkLogger) WithLogLevel(value string) *SplunkLogger {
	l.items["logLevel"] = value
	return l
}

func (l *SplunkLogger) reset() {
	l.items = make(map[string]string, 0)
	l.err = nil
}

func (l *SplunkLogger) Log() error {

	if l.err != nil {
		l.items["error"] = l.err.Error()
	}
	l.items["logtime"] = fmt.Sprintf("%d", time.Now().Unix())
	data, err := json.Marshal(l.items)
	if err != nil {
		l.reset()
		return err
	}

	event := make(map[string]interface{}, 0)

	event["sourcetype"] = "_json"
	event["event"] = string(data)

	log, err := json.Marshal(event)
	if err != nil {
		l.reset()
		return err
	}

	fmt.Printf("%s\n", "--------------")
	fmt.Printf("%s\n", string(log))
	fmt.Printf("%s\n", "--------------")

	req, err := http.NewRequest("POST", l.URL, bytes.NewReader(log))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Splunk %s", l.token))

	resp, err := l.httpClient.Do(req)
	if err != nil {
		l.reset()
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}
	return nil

}
