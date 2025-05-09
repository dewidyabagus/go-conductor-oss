package main

import (
	"encoding/json"
	"net/http"

	"github.com/dewidyabagus/go-payout-workflow/sources/model"
	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/utils"
	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/validator"
)

type handler struct {
	validator *validator.Validate
	service   *service
}

func (h *handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	request := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		JSONErrorResponse(w, utils.WrapError(utils.ErrRequest, err))
		return
	}

	if resp, err := h.service.HelloWorld(r.Context(), request["message"]); err != nil {
		JSONErrorResponse(w, err)
	} else {
		JSONSuccessResponse(w, http.StatusOK, resp)
	}
}

func (h *handler) PrepaidPaymentHandler(w http.ResponseWriter, r *http.Request) {
	request := model.PrepaidPaymentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		JSONErrorResponse(w, utils.WrapError(utils.ErrRequest, err))
		return
	}

	response, err := h.service.PrepaidPayment(r.Context(), request)
	if err != nil {
		JSONErrorResponse(w, err)
		return
	}
	JSONSuccessResponse(w, http.StatusOK, response)
}

func (h *handler) HarsyaPaymentNotificationHandler(w http.ResponseWriter, r *http.Request) {
	request := model.HarsyaPaymentNotificationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		JSONErrorResponse(w, utils.WrapError(utils.ErrRequest, err))
		return
	}

	if err := h.validator.Struct(request); err != nil {
		JSONErrorResponse(w, utils.WrapError(utils.ErrInvalidReqContent, err))
		return
	}

	if err := h.service.HarsyaPaymentNotification(r.Context(), request); err != nil {
		JSONErrorResponse(w, err)
		return
	}
	JSONSuccessResponse(w, http.StatusOK, map[string]string{"status": "OK"})
}
