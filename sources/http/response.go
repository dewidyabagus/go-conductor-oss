package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dewidyabagus/go-payout-workflow/sources/model"
	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/utils"
	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/validator"
)

func SetJSONResponse(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}

func JSONErrorResponse(w http.ResponseWriter, e error) {
	code, err := utils.UnwrapHttpError(e)
	SetJSONResponse(w, code)

	response := model.Response{
		Status: "error",
		Error:  &model.ErrorResponse{},
	}
	switch t := err.(type) {
	case *validator.Errors:
		response.Error.Details = t
		response.Error.Message = "invalid request content"

	case error:
		response.Error.Details = err.Error()
		response.Error.Message = "others"
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("JSON Decode: %v", err), http.StatusInternalServerError)
	}
}

func JSONSuccessResponse(w http.ResponseWriter, code int, data any) {
	SetJSONResponse(w, code)

	response := model.Response{
		Status: "success",
		Data:   data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("JSON Decode: %v", err), http.StatusInternalServerError)
	}
}
