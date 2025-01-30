// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

package fleetautomationextension

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"

	"go.uber.org/zap"
)

const (
	internalIntakeURLTemplate = "%s/v2/track/%s"
)

const (
	defaultAsyncIntakeTimeout = time.Second * 30
	defaultAsyncIntakeRetries = 4
	retryMaxWait              = 5 * time.Second
	retryMultiplier           = 1.5
	retryMaxJitter            = 500 * time.Millisecond
	retryMinInitialWait       = 500 * time.Millisecond
	retryMaxInitialWait       = time.Second
	retryMinJitter            = 0
	defaultTrack              = "genresources"
	agentMetadataTrack        = "agentmetadata"
)

// Clock provides an interface to get current time
type Clock interface {
	Now() time.Time
}

// SystemClock implements Clock using time.Now
type SystemClock struct{}

// Now returns the current time
func (SystemClock) Now() time.Time { return time.Now() }

func getSupportedTracks() []string {
	return []string{agentMetadataTrack}
}

// AsyncIntakeClientConfig holds configurations for the resource grpc client
type AsyncIntakeClientConfig struct {
	Clock    Clock
	Source   string
	Track    string
	Endpoint string

	// the max total number of calls
	MaxRetries          int
	Timeout             time.Duration
	MaxWait             time.Duration
	RetryMinInitialWait time.Duration
	RetryMaxInitialWait time.Duration
	RetryMultiplier     float64
	RetryMinJitter      time.Duration
	RetryMaxJitter      time.Duration

	supportedTracks StringSet
}

// apply runs the given function(s) on the calling AsyncIntakeClientConfig object.
func (cc *AsyncIntakeClientConfig) apply(options ...func(*AsyncIntakeClientConfig)) *AsyncIntakeClientConfig {
	for _, option := range options {
		option(cc)
	}
	return cc
}

func (cc *AsyncIntakeClientConfig) validate() error {
	if cc.Source == "" {
		return errors.New("source is required and should be set to the name of your service")
	}
	if cc.Endpoint == "" {
		return errors.New("the Event Platform endpoint is required")
	}
	if !cc.isTrackSupported(cc.Track) {
		return fmt.Errorf("track '%s' is not supported", cc.Track)
	}
	return nil
}

func (cc *AsyncIntakeClientConfig) isTrackSupported(track string) bool {
	return cc.supportedTracks.Contains(track)
}

// NewAsyncIntakeClientConfig returns the default AsyncIntakeClientConfig
func NewAsyncIntakeClientConfig() *AsyncIntakeClientConfig {
	result := &AsyncIntakeClientConfig{
		Clock:               SystemClock{},
		MaxRetries:          defaultAsyncIntakeRetries,
		Timeout:             defaultAsyncIntakeTimeout,
		MaxWait:             retryMaxWait,
		RetryMinInitialWait: retryMinInitialWait,
		RetryMaxInitialWait: retryMaxInitialWait,
		RetryMinJitter:      retryMinJitter,
		RetryMaxJitter:      retryMaxJitter,
		RetryMultiplier:     retryMultiplier,
		Track:               defaultTrack,
	}
	result.supportedTracks = NewStringSetFromSlice(getSupportedTracks())
	return result
}

// RedaplAsyncIntakeClient is a client that ingests REDAPL resources via the Event Platform intake
// The resources are not validated against the schema by this client. If a resource is invalid (e.g. it has a field
// not specified in the type's schema) it will be successfully ingested by this client but discarded later in the pipeline.
type RedaplAsyncIntakeClient interface {
	// Send upserts a single resource to REDAPL via the Event Platform Intake
	// Returned errors can be safely retried
	Send(ctx context.Context, source string, encodedMessage RawResourceV3) error
	// SendBatch upserts a batch of resources to REDAPL via the Event Platform Intake
	// The request is atomic, i.e. either the whole batch succeeds or fails
	// Returned errors can be safely retried
	SendBatch(ctx context.Context, source string, resources []RawResourceV3) error
}

// HTTPRedaplAsyncIntakeClient implements RedaplAsyncIntakeClient
type HTTPRedaplAsyncIntakeClient struct {
	clock  Clock
	client HTTPClient
	cfg    *AsyncIntakeClientConfig
	log    *zap.Logger
}

// SendError is returned when there is an error with sending to event platform
type SendError struct {
	Retriable  bool
	Err        error
	Reason     string
	ShouldExit bool
}

// HTTPClient is an interface that http.Client implements
// This is needed in order to be able to mock http.Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (e *SendError) Error() string {
	return e.Err.Error()
}

func newSendError(err error, retriable bool, reason string, shouldExit bool) error {
	return &SendError{
		Retriable:  retriable,
		Err:        err,
		Reason:     reason,
		ShouldExit: shouldExit,
	}
}

var _ RedaplAsyncIntakeClient = (*HTTPRedaplAsyncIntakeClient)(nil) // Forces implementation of the interface

// NewHTTPRedaplAsyncIntakeClient ..
func NewHTTPRedaplAsyncIntakeClient(log *zap.Logger, options ...func(*AsyncIntakeClientConfig)) (*HTTPRedaplAsyncIntakeClient, error) {
	cfg := NewAsyncIntakeClientConfig()
	cfg.apply(options...)
	err := cfg.validate()
	if err != nil {
		return nil, err
	}

	return &HTTPRedaplAsyncIntakeClient{
		clock: cfg.Clock,
		cfg:   cfg,
		log:   log,
	}, nil
}

// Send sends a message to event platform intake
// If a context with a deadline (timeout) is passed it will apply to all retry attempts combined
// If you choose to set the deadline it should be large enough to allow all the retry attempts according to your
// client config. The suggested timeout for the default config is 9 seconds or more.
func (c *HTTPRedaplAsyncIntakeClient) Send(ctx context.Context, source string, r RawResourceV3) error {
	err := c.SendBatch(ctx, source, []RawResourceV3{r})

	return err

}

// SendBatch sends messages to event platform intake
func (c *HTTPRedaplAsyncIntakeClient) SendBatch(ctx context.Context, source string, resources []RawResourceV3) error {
	var err error

	now := c.clock.Now()

	req, err := createRequest(now, resources, source)
	if err != nil {
		return err
	}
	body, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	err = c.send(ctx, body)

	if err != nil {
		return err
	}

	return nil
}

func (c *HTTPRedaplAsyncIntakeClient) send(ctx context.Context, body []byte) error {
	doRequest := func() error {
		url := fmt.Sprintf(internalIntakeURLTemplate, c.cfg.Endpoint, c.cfg.Track)
		reqCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return newSendError(err, false, "invalid_request", true)
		}
		req.Header.Add("Content-Type", "application/x-protobuf")

		res, err := c.client.Do(req)

		if err != nil {
			return newSendError(err, true, "client_error", true)
		}
		defer res.Body.Close()

		switch {
		case 200 <= res.StatusCode && res.StatusCode < 300:
			c.log.Debug("Request to EvP sent successfully", zap.Int("status code", res.StatusCode), zap.String("status", res.Status))
			return nil
		case res.StatusCode == 400:
			buf := new(strings.Builder)
			_, _ = io.Copy(buf, res.Body)
			respBody := buf.String()

			var reqProto internalIntakeRequest
			_ = proto.Unmarshal(body, &reqProto)
			reqString := reqProto.String()

			return newSendError(fmt.Errorf("request returned 400 status code : %s ;\n request was : %s", respBody, reqString), false, respBody, false)
		case res.StatusCode == 403:
			// ignore 403 errors because in this case it means that some
			// of the orgs in the batch have no or expired Logs subscription
			// intake will process the rest of the orgs' payloads correctly
			c.log.Info("EvP Intake returned 403")
			return nil
		case res.StatusCode == 413:
			// 413 errors mean that the payload is too big for the event platform.
			// We definitely don't want to retry in this case
			return newSendError(fmt.Errorf("request returned 413 status code"), false, fmt.Sprintf("status_code:%d", res.StatusCode), false)
		default:
			return newSendError(fmt.Errorf("request returned non-200 status code (code=%d): %s",
				res.StatusCode, res.Body), true, fmt.Sprintf("%d", res.StatusCode), true)
		}
	}

	doAndLogRequest := func() error {
		if err := doRequest(); err != nil {
			c.log.Error("Error sending to event platform", zap.Any("error", err))
			return err
		}
		return nil
	}

	backoffRetrier := c.getbackoffRetrier()
	if err := backoffRetrier.DoContext(ctx, doAndLogRequest); err != nil {
		return err
	}
	return nil
}

func createRequest(now time.Time, resources []RawResourceV3, source string) (*internalIntakeRequest, error) {
	var events []*internalIntakeRequestEvent
	for _, msg := range resources {

		rProto := &RawResourceV3{
			OrgID:        msg.OrgID,
			Type:         msg.Type,
			Name:         msg.Name,
			FieldsByName: msg.FieldsByName,
			SeenAt:       msg.SeenAt,
			ExpireAt:     msg.ExpireAt,
		}

		ser, err := proto.Marshal(rProto)
		if err != nil {
			return nil, err
		}

		redaplEvent, err := proto.Marshal(&redaplEvent{
			Source:  source,
			Message: ser,
		})
		if err != nil {
			return nil, err
		}

		event := &internalIntakeRequestEvent{
			// UUID is set to discard duplicate events in case of retries
			UUID:      &wrappers.StringValue{Value: uuid.NewString()},
			Payload:   redaplEvent,
			OrgID:     &wrappers.Int64Value{Value: msg.OrgID},
			Timestamp: &wrappers.Int64Value{Value: now.UnixMilli()},
		}

		events = append(events, event)
	}
	return &internalIntakeRequest{Events: events}, nil
}

func (c *HTTPRedaplAsyncIntakeClient) getbackoffRetrier() backoffRetrier {
	return backoffRetrier{
		Times: c.cfg.MaxRetries,
		Strategy: BackoffStrategy{
			InitialWait: []time.Duration{c.cfg.RetryMinInitialWait, c.cfg.RetryMaxInitialWait},
			Multiplier:  c.cfg.RetryMultiplier,
			Jitter:      []time.Duration{c.cfg.RetryMinJitter, c.cfg.RetryMaxJitter},
			MaxWait:     c.cfg.MaxWait,
		},
		Retriable: func(err error) RetriableCode {
			var e *SendError
			if !errors.As(err, &e) || !e.Retriable {
				return DoNotRetry
			}
			return ShouldRetry
		},
	}
}
