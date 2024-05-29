package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/a-x-a/go-loyalty/internal/service/accrualservice/customerrors"
	"github.com/a-x-a/go-loyalty/internal/service/accrualservice/model"
)

type MockHTTPCLient struct {
	expectedResponse *http.Response
	expectedErr      error
}

func (mock *MockHTTPCLient) Get(_ string) (*http.Response, error) {
	return mock.expectedResponse, mock.expectedErr
}

func (mock *MockHTTPCLient) Expect(exResp *http.Response, exErr error) {
	mock.expectedResponse = exResp
	mock.expectedErr = exErr
}

func TestServiceAccrual_ProcessedAccrualData(t *testing.T) {
	type (
		clientMock struct {
			resp *http.Response
			err  error
		}

		expected struct {
			accrualOrder model.AccrualOrder
			err          error
		}
	)

	emptyAccrualOrder := model.AccrualOrder{}
	testOrderNum := "371449635398431"
	someErr := errors.New("some error")
	badResp := fmt.Sprintf("{\"order\":\"%v\",\"status\":\"SomeStatus\",\"accrual\":500}", testOrderNum)
	goodResp := fmt.Sprintf("{\"order\":\"%v\",\"status\":\"PROCESSED\",\"accrual\":500}", testOrderNum)
	testAddr := "testhost:8090"

	tooManyRequestsHeader := make(map[string][]string)
	tooManyRequestsHeader["Retry-After"] = []string{"1"}

	tests := []struct {
		name string
		clientMock
		expected
		clientIsAvailable bool
	}{
		{
			"successful response",
			clientMock{&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(goodResp))),
			}, nil},
			expected{
				model.AccrualOrder{
					Order:   "371449635398431",
					Status:  model.PROCESSED.String(),
					Accrual: 500}, nil},
			true,
		},
		{
			"error response",
			clientMock{nil, someErr},
			expected{emptyAccrualOrder, someErr},
			true,
		},
		{
			"client is no available",
			clientMock{nil, someErr},
			expected{emptyAccrualOrder, customerrors.ErrClientIsNoAvailable},
			false,
		},
		{
			"bad response",
			clientMock{&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(badResp))),
			}, nil},
			expected{emptyAccrualOrder, customerrors.ErrInvalidAccrualOrder},
			true,
		},
		{
			"unknown status",
			clientMock{&http.Response{StatusCode: http.StatusNotImplemented,
				Body: io.NopCloser(bytes.NewReader([]byte(badResp)))}, nil},
			expected{emptyAccrualOrder, nil},
			true,
		},
		// {
		// 	"status ok without body",
		// 	clientMock{&http.Response{StatusCode: http.StatusOK,
		// 		Body: io.NopCloser(bytes.NewReader([]byte(badResp)))}, nil},
		// 	expected{emptyAccrualOrder, someErr},
		// 	true,
		// },
		{
			"status no content",
			clientMock{&http.Response{StatusCode: http.StatusNoContent,
				Body: io.NopCloser(bytes.NewReader([]byte(badResp)))}, nil},
			expected{emptyAccrualOrder, customerrors.ErrNoContent},
			true,
		},
		{
			"status too many requests without header",
			clientMock{&http.Response{StatusCode: http.StatusTooManyRequests,
				Body: io.NopCloser(bytes.NewReader([]byte(badResp)))}, nil},
			expected{emptyAccrualOrder, customerrors.ErrTooManyRequests},
			true,
		},
		{
			"status too many requests with header",
			clientMock{&http.Response{StatusCode: http.StatusTooManyRequests,
				Header: tooManyRequestsHeader,
				Body:   io.NopCloser(bytes.NewReader([]byte(badResp)))}, nil},
			expected{emptyAccrualOrder, customerrors.ErrTooManyRequests},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accrualClient := New(testAddr, zap.L())
			mockClient := MockHTTPCLient{}
			accrualClient.client = &mockClient
			accrualClient.isAvailable.Store(tt.clientIsAvailable)

			mockClient.Expect(tt.clientMock.resp, tt.clientMock.err)

			ctx := context.Background()
			order, err := accrualClient.Get(ctx, testOrderNum)
			if err != nil {
				assert.Equal(t, true, errors.Is(err, tt.expected.err))
			} else {
				assert.Equal(t, tt.expected.accrualOrder, order)
			}
		})
	}
}
