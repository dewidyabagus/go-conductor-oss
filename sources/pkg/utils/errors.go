package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	conductor "github.com/conductor-sdk/conductor-go/sdk/model"
)

type ErrType interface {
	Type() string
}

type errWrap struct {
	typ ErrType
	err error
}

var (
	ErrRequest           = &errRequest{}
	ErrInvalidReqContent = &errInvalidReqContent{}
	ErrDatabase          = &errDatabase{}
	ErrInternal          = &errInternal{}
	ErrRetryableClient   = &errRetryableClient{}
	ErrDataNotFound      = &errDataNotFound{}
)

func (e *errWrap) Error() string {
	return fmt.Sprintf("%v | %v", e.typ, e.err)
}

type errRequest struct{}

func (errRequest) Type() string {
	return "INVALID_REQUEST"
}

type errInvalidReqContent struct{}

func (errInvalidReqContent) Type() string {
	return "INVALID_REQUEST_CONTENT"
}

type errDatabase struct{}

func (errDatabase) Type() string {
	return "ERROR_DATABASE"
}

type errInternal struct{}

func (errInternal) Type() string {
	return "ERROR_INTERNAL_SERVICE"
}

type errGeneral struct{}

func (errGeneral) Type() string {
	return "ERROR_GENERAL"
}

type errRetryableClient struct{}

func (errRetryableClient) Type() string {
	return "ERROR_RETRYABLE_CLIENT"
}

type errDataNotFound struct{}

func (errDataNotFound) Type() string {
	return "ERROR_DATA_NOT_FOUND"
}

func WrapError(typ ErrType, err error) error {
	return &errWrap{typ, err}
}

func UnwrapHttpError(err error) (int, error) {
	if err == nil {
		return 0, nil
	}

	e, ok := err.(*errWrap)
	if !ok {
		return http.StatusInternalServerError, err
	}

	switch e.typ.(type) {
	default:
		return http.StatusInternalServerError, e.err

	case *errRequest, *errInvalidReqContent:
		return http.StatusBadRequest, e.err

	case *errDataNotFound:
		return http.StatusUnprocessableEntity, e.err

	case *errDatabase, *errInternal:
		log.Println("Database Error:", e.err)
		return http.StatusInternalServerError, errors.New("there was a problem with our internal servers")
	}
}

func UnwrapTaskErrorResponse(task *conductor.Task, err error) (*conductor.TaskResult, error) {
	if err == nil {
		return UnwrapTaskErrorResponse(task, WrapError(&errGeneral{}, errors.New("warning: error is nil")))
	}

	e, ok := err.(*errWrap)
	if !ok {
		return UnwrapTaskErrorResponse(task, WrapError(&errGeneral{}, err))
	}

	switch e.typ.(type) {
	default:
		err = conductor.NewNonRetryableError(e.err)

	case *errDatabase, *errRetryableClient:
		err = e.err
	}
	return conductor.NewTaskResultFromTaskWithError(task, err), nil
}

func TaskNonRetryableErrorResponse(task *conductor.Task, err error) (*conductor.TaskResult, error) {
	return conductor.NewTaskResultFromTaskWithError(
		task, conductor.NewNonRetryableError(err),
	), err
}
