package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ClientInterface interface {
	Create(path string, body, v interface{}) error
	Fetch(path string, data interface{}) error
	Delete(path string, v interface{}) error
}

type Config struct {
	host string
}

type Client struct {
	config     *Config
	httpClient *http.Client
}

func New(host string) ClientInterface {
	return &Client{
		config:     &Config{host},
		httpClient: &http.Client{},
	}
}

func (c *Client) Create(path string, body, v interface{}) error {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(&body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, c.resolveUrl(path), buffer)
	if err != nil {
		return err
	}

	err = c.doRequest(request, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Fetch(path string, v interface{}) error {
	request, err := http.NewRequest(http.MethodGet, c.resolveUrl(path), nil)
	if err != nil {
		return err
	}

	err = c.doRequest(request, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Delete(path string, v interface{}) error {
	request, err := http.NewRequest(http.MethodDelete, c.resolveUrl(path), nil)
	if err != nil {
		return err
	}

	err = c.doRequest(request, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) doRequest(request *http.Request, v interface{}) error {
	request.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var restError error
		switch res.StatusCode {
		case 400:
			restError = &ErrBadRequest{}
		case 404:
			restError = &ErrNotFound{}
		case 409:
			restError = &ErrConflict{}
		case 500:
			restError = &ErrInternalServer{}
		}

		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		json.Unmarshal(bytes, &restError)
		return restError
	}

	if v != nil {
		err = json.NewDecoder(res.Body).Decode(&v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) resolveUrl(path string) string {
	return c.config.host + path
}

func getErr(statusCode int) error {
	switch statusCode {
	case 400:
		return &ErrBadRequest{}
	case 404:
		return &ErrNotFound{}
	case 409:
		return &ErrConflict{}
	case 500:
		return &ErrInternalServer{}
	}
	return errors.New("Unknown status code")
}
