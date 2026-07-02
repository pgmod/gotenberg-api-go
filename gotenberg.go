// Package gotenbergapigo — клиент для Gotenberg (конвертация в PDF).
package gotenbergapigo

import (
	"bytes"
	"errors"
	"fmt"

	"resty.dev/v3"
)

var ErrHostNotSet = errors.New("gotenberg: адрес сервера не задан")

// Client — HTTP-клиент к Gotenberg с общим base URL.
type Client struct {
	baseURL string
	http    *resty.Client
}

// New создаёт клиент. baseURL — например http://localhost:3000 (без завершающего слэша).
func New(baseURL string) (*Client, error) {
	if baseURL == "" {
		return nil, ErrHostNotSet
	}
	for len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	return &Client{
		baseURL: baseURL,
		http:    resty.New(),
	}, nil
}

// Close освобождает ресурсы HTTP-клиента.
func (c *Client) Close() {
	if c != nil && c.http != nil {
		c.http.Close()
	}
}

func httpError(resp *resty.Response) error {
	return fmt.Errorf("HTTP %d: %s", resp.StatusCode(), resp.String())
}

// DocxToPdf конвертирует DOCX (LibreOffice) в PDF. filename — имя части multipart; пустая строка → document.docx.
func (c *Client) DocxToPdf(docx []byte, filename string) ([]byte, error) {
	if len(docx) == 0 {
		return nil, fmt.Errorf("пустой вход DOCX")
	}
	resp, err := c.http.R().
		SetFileReader("file", "document.docx", bytes.NewReader(docx)).
		Post(c.baseURL + "/forms/libreoffice/convert")
	if err != nil {
		return nil, err
	}
	if resp.IsStatusFailure() {
		return nil, httpError(resp)
	}
	return resp.Bytes(), nil
}

// HTMLToPdf конвертирует HTML (Chromium) в PDF. filename — имя части multipart; пустая строка → index.html.
func (c *Client) HTMLToPdf(html []byte) ([]byte, error) {
	if len(html) == 0 {
		return nil, fmt.Errorf("пустой вход HTML")
	}
	resp, err := c.http.R().
		SetFileReader("files", "index.html", bytes.NewReader(html)).
		Post(c.baseURL + "/forms/chromium/convert/html")
	if err != nil {
		return nil, err
	}
	if resp.IsStatusFailure() {
		return nil, httpError(resp)
	}
	return resp.Bytes(), nil
}
