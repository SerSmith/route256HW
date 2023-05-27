package cliwrapper

import (
	"time"
	"context"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	TIMEOUT = 5 * time.Second
)


func RequestAPI[Req any, Res any](ctx context.Context, handle string, url string, request Req) (Res, error) {
	
	var out Res

	rawData, err := json.Marshal(&request)
	if err != nil {
		return out, fmt.Errorf("encode url: %s handle: %s request: %w", url, handle, err)
	}

	ctx, fnCancel := context.WithTimeout(ctx, TIMEOUT)
	defer fnCancel()

	httpRequest, err := http.NewRequestWithContext(ctx, handle, url, bytes.NewBuffer(rawData))
	if err != nil {
		return out, fmt.Errorf("prepare url: %s handle: %s request: %w", url, handle, err)
	}

	response, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return out, fmt.Errorf("do url: %s : %w", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return out, fmt.Errorf("wrong status code get url: %s handle: %s: %d",  url, handle, response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&out)
	if err != nil {
		return out, fmt.Errorf("decode  url: %s handle: %s request: %w",  url, handle, err)
	}

	return out, nil
}
