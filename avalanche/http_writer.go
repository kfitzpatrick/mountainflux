package avalanche

import (
	"fmt"
	"net/url"

	"github.com/valyala/fasthttp"
)

type HTTPWriterConfig struct {
	Host     string
	Database string
}

type HTTPWriter struct {
	client fasthttp.Client

	c   HTTPWriterConfig
	url []byte
}

// NewHTTPWriter returns a new HTTPWriter from the supplied HTTPWriterConfig.
func NewHTTPWriter(c HTTPWriterConfig) LineProtocolWriter {
	return &HTTPWriter{
		client: fasthttp.Client{
			Name: "avalanche",
		},

		c:   c,
		url: []byte(c.Host + "/write?db=" + url.QueryEscape(c.Database)),
	}
}

var post = []byte("POST")

func (w *HTTPWriter) WriteLineProtocol(body []byte) error {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethodBytes(post)
	req.Header.SetRequestURIBytes(w.url)
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	err := w.client.Do(req, resp)
	if err == nil {
		sc := resp.StatusCode()
		if sc != fasthttp.StatusNoContent {
			err = fmt.Errorf("Invalid write response (status %d): %s", sc, resp.Body())
		}
	}

	fasthttp.ReleaseResponse(resp)
	fasthttp.ReleaseRequest(req)

	return err
}
