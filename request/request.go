package request

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Send(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("fetching user sales: %w", err)
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making %s request to %s. ERROR: %w", method, url, err)
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code: %d", r.StatusCode)
	}

	return r, nil
}
