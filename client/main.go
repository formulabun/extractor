package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.formulabun.club/extractor/env"
	"go.formulabun.club/metadatadb"
)

type Client struct {
	c *http.Client

	buff []metadatadb.File
}

func NewClient() *Client {
	var capacity = 5
	c := Client{
		http.DefaultClient,
		make([]metadatadb.File, 0, capacity),
	}
	return &c
}

func (c *Client) ExtractFile(data metadatadb.File) error {
	var err error
	if len(c.buff) == cap(c.buff) {
		if err = c.ExtractFiles(c.buff); err != nil {
			return nil
		}
		c.buff = c.buff[:0]
	} else {
		c.buff = append(c.buff, data)
	}
	return nil
}

func (c *Client) ExtractFiles(data []metadatadb.File) error {
	if len(data) == 0 {
		return nil
	}

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.Encode(data)

	resp, err := c.c.Post("http://"+env.Host, "application/json", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		body, err := io.ReadAll(resp.Body)
		return errors.Join(fmt.Errorf("Bad request: %s", body), err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.New("Unknown error in the extractor server")
	}

	return nil
}
