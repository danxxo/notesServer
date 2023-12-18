package dto

import "encoding/json"

type Response struct {
	Data   json.RawMessage
	Error  string
	Result string
}

func (resp *Response) Wrap(result string, data json.RawMessage, err error) {
	if err == nil {
		resp.Error = ""
	} else {
		resp.Error = err.Error()
	}
	resp.Data = data
	resp.Result = result
}
