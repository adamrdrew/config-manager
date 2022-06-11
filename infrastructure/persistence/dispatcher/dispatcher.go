package dispatcher

// package dispatcher contains code generated from an OpenAPI spec derived from
// the playbook-dispatcher project. First, swagger-cli dereferences and bundles
// the public.openapi.yaml and private.openapi.yaml files. Then the boilerplate
// code is generated by the github.com/deepmap/oapi-codegen project. At the time
// of this writing, it was generated using version 1.5.1. Go does not support
// path@version synytax with 'go run', so oapi-codegen must be installed
// manually before code can be generated:
// 'go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@1.5.1'

//go:generate swagger-cli bundle --outfile /tmp/openapi.yaml https://github.com/RedHatInsights/playbook-dispatcher/raw/6950cf45d9c0ff464ecbfcd1d83d2cc0797fe41f/schema/private.openapi.yaml
//go:generate oapi-codegen -generate types -package dispatcher -o ./dispatcher_types.gen.go /tmp/openapi.yaml
//go:generate oapi-codegen -generate client -package dispatcher -o ./dispatcher.gen.go /tmp/openapi.yaml

import (
	"config-manager/internal/config"
	"context"
	"fmt"
	"net/http"
	"time"
)

// DispatcherClient provides REST client API methods to interact with the
// platform playbook-dispatcher application.
type DispatcherClient interface {
	Dispatch(ctx context.Context, inputs []RunInput) ([]RunCreated, error)
}

// dispatcherClientImpl implements DispatcherClient interface.
type dispatcherClientImpl struct {
	client ClientWithResponsesInterface
}

// NewDispatcherClientWithDoer returns a DispatchClient by constructing a
// dispatcher.Client, configured with request headers and host information.
func NewDispatcherClientWithDoer(doer HttpRequestDoer) DispatcherClient {
	client := &ClientWithResponses{
		ClientInterface: &Client{
			Server: config.DefaultConfig.DispatcherHost.String(),
			Client: doer,
			RequestEditors: []RequestEditorFn{
				func(ctx context.Context, req *http.Request) error {
					req.Header.Set("Authorization", fmt.Sprintf("PSK %s", config.DefaultConfig.DispatcherPSK))
					req.Header.Set("Content-Type", "application/json")
					return nil
				},
			},
		},
	}

	return &dispatcherClientImpl{
		client: client,
	}
}

// NewDispatcherClient creates a new DispatcherClient.
func NewDispatcherClient() DispatcherClient {
	client := &http.Client{
		Timeout: time.Duration(int(time.Second) * config.DefaultConfig.DispatcherTimeout),
	}

	return NewDispatcherClientWithDoer(client)
}

// Dispatch performs the CreateWithResponse API method of the
// playbook-dispatcher service.
func (dc *dispatcherClientImpl) Dispatch(ctx context.Context, inputs []RunInput) ([]RunCreated, error) {
	res, err := dc.client.ApiInternalRunsCreateWithResponse(ctx, inputs)
	if err != nil {
		return nil, err
	}

	if res.HTTPResponse.StatusCode != 207 {
		return nil, fmt.Errorf("Unexpected error code %d received: %v", res.HTTPResponse.StatusCode, string(res.Body))
	}

	return *res.JSON207, nil
}
