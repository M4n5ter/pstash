package zo

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/m4n5ter/pstash/stash/config"
	"io"
	"net/http"
	"strings"
)

type Writer struct {
	client http.Client
	zoUrl  string
	// 最近一次的响应
	LastResp ZOResp
}

type ZOResp struct {
	Code   int `json:"code"`
	Status []struct {
		Name       string `json:"name"`
		Successful int    `json:"successful"`
		Failed     int    `json:"failed"`
	} `json:"status"`
}

func NewWriter(conf config.ZincObserveConf) *Writer {
	zoUrl := fmt.Sprintf("%s://%s:%s@%s/%s/%s/%s",
		conf.Schema, conf.Username, conf.Password, conf.Host, conf.Organization, conf.Stream, conf.IngestionType)

	return &Writer{
		client:   http.Client{},
		zoUrl:    zoUrl,
		LastResp: ZOResp{},
	}
}

// Write writes the val to the ZincObserve endpoint
func (c *Writer) Write(_, val string) error {
	req, err := http.NewRequest("POST", c.zoUrl, strings.NewReader(val))
	if err != nil {
		return fmt.Errorf("fail to create request: %w", err)
	}
	resp, err := c.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("fail to send request to ZincObserve: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code is not 200: %d", resp.StatusCode)
	}

	zpb, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("fail to read response body: %w", err)
	}

	err = json.Unmarshal(zpb, &c.LastResp)
	if err != nil {
		return fmt.Errorf("fail to unmarshal response body: %w", err)
	}
	return nil
}
