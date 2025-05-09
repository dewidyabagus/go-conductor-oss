package workflow

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/http"

	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/settings"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
	"github.com/conductor-sdk/conductor-go/sdk/workflow/executor"
)

type Config struct {
	BaseURL       string
	Authorization Authorization
}

type TaskRunner = worker.TaskRunner
type WorkflowExecutor = executor.WorkflowExecutor

type Authorization interface {
	Encode() string
}

type BasicAuth struct {
	Username string
	Password string
}

func (b *BasicAuth) Encode() string {
	return base64.RawURLEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", b.Username, b.Password))
}

type workflow struct {
	apiClient *client.APIClient

	httpClient  http.HttpClient
	httpHeaders map[string]string
	mx          *sync.Mutex
	runner      *worker.TaskRunner
	executor    *executor.WorkflowExecutor
}

func New(cfg Config) *workflow {
	httpSettings := &settings.HttpSettings{
		BaseUrl: cfg.BaseURL,
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"Accept":          "application/json",
			"Accept-Encoding": "gzip",
		},
	}
	if cfg.Authorization != nil {
		switch auth := cfg.Authorization.(type) {
		default:
			// Default handling

		case *BasicAuth:
			httpSettings.Headers["Authorization"] = "Basic " + auth.Encode()
		}
	}

	client := client.NewAPIClient(nil, httpSettings)

	url, _ := url.Parse(cfg.BaseURL)

	workflow := &workflow{
		apiClient:   client,
		mx:          new(sync.Mutex),
		httpClient:  http.NewHttpClient(fmt.Sprintf("%s://%s", url.Scheme, url.Host)),
		httpHeaders: httpSettings.Headers,
	}
	return workflow
}

func (w *workflow) TaskRunner() *TaskRunner {
	w.mx.Lock()
	defer w.mx.Unlock()

	if w.runner == nil {
		w.runner = worker.NewTaskRunnerWithApiClient(w.apiClient)
	}
	return w.runner
}

func (w *workflow) WorkflowExecutor() *WorkflowExecutor {
	w.mx.Lock()
	defer w.mx.Unlock()

	if w.executor == nil {
		w.executor = executor.NewWorkflowExecutor(w.apiClient)
	}
	return w.executor
}

func (w *workflow) Workers() *workers {
	return &workers{taskRunner: w.TaskRunner()}
}

func (w *workflow) HealthCheck(ctx context.Context) error {

	if resp, err := w.httpClient.Get(ctx, "/health", w.httpHeaders); err != nil {
		fmt.Println("HTTP Client Error:", err.Error())
		return err

	} else if resp.StatusCode() != 200 {
		fmt.Printf("Unexpected Response:\nCode: %d\nResponse: %s", resp.StatusCode(), string(resp.BodyBytes()))
		return errors.New("unhealthy workflow engine")
	}
	return nil
}
