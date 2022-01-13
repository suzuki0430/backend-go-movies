package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	// メモリを割り当てる
	wrapper := make(map[string]interface{})

	// wrapキーの構造体にdataを代入する
	wrapper[wrap] = data

	// 構造体をJSONに変換する
	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	
	// ヘッダーにContent-Typeを付加する
	w.Header().Set("Content-Type", "application/json")
	// ステータスコードとともにHTTP応答ヘッダーを返す
	w.WriteHeader(status)
	// ボディを返す
	w.Write(js)

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError {
		Message: err.Error(),
	}

	app.writeJSON(w, http.StatusBadRequest, theError, "error")
}