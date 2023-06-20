package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HTTPGetJSON(url string, dst interface{}) error {
	var (
		rawResp *http.Response
		body    []byte
		err     error
	)
	if rawResp, err = http.Get(url); err != nil {
		return fmt.Errorf("get http response failed:%w", err)
	}
	if rawResp != nil {
		defer rawResp.Body.Close()
	}
	if body, err = io.ReadAll(rawResp.Body); err != nil {
		return fmt.Errorf("read data from body failed:%w", err)
	}
	if err = json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("unmarshal data failed:%w", err)
	}
	return nil
}
