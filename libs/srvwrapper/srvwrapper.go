package srvwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"route256/libs/logger"
	"net/http"
)

type Wrapper[Req Validator, Res any] struct {
	fn func(ctx context.Context, req Req) (Res, error)
}

type Validator interface {
	Validate() error
}

func New[Req Validator, Res any](fn func(ctx context.Context, req Req) (Res, error)) *Wrapper[Req, Res] {
	return &Wrapper[Req, Res]{
		fn: fn,
	}
}

func (w *Wrapper[Req, Res]) ServeHTTP(resWriter http.ResponseWriter, httpReq *http.Request) {

	var req Req

	err := json.NewDecoder(httpReq.Body).Decode(&req)
	if err != nil {
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "parse request", err)
		return
	}

	//reqValidation, ok := any(req).(Validator)
	//if ok {
	//	errValidation := reqValidation.Validate()
	//	if errValidation != nil {
	//		resWriter.WriteHeader(http.StatusBadRequest)
	//		writeErrorText(resWriter, "bad request", errValidation)
	//		return
	//	}
	//}

	errValidation := req.Validate()
	if errValidation != nil {
		resWriter.WriteHeader(http.StatusBadRequest)
		writeErrorText(resWriter, "bad request", errValidation)
		return
	}

	resp, err := w.fn(httpReq.Context(), req)
	if err != nil {
		logger.Info("executor fail: %s", err)
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "exec handler", err)
		return
	}

	rawData, err := json.Marshal(&resp)
	if err != nil {
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "decode response", err)
		return
	}

	_, _ = resWriter.Write(rawData)
}

func writeErrorText(w http.ResponseWriter, text string, err error) {
	buf := bytes.NewBufferString(text)

	buf.WriteString(": ")
	buf.WriteString(err.Error())
	buf.WriteByte('\n')

	w.Write(buf.Bytes())
}
