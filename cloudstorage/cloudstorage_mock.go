// Code generated by MockGen. DO NOT EDIT.
// Source: cloudstorage.go

// Package cloudstorage is a generated GoMock package.
package cloudstorage

import (
	context "context"
	io "io"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// GoMockClient is a mock of Client interface.
type GoMockClient struct {
	ctrl     *gomock.Controller
	recorder *GoMockClientMockRecorder
}

// GoMockClientMockRecorder is the mock recorder for GoMockClient.
type GoMockClientMockRecorder struct {
	mock *GoMockClient
}

// NewGoMockClient creates a new mock instance.
func NewGoMockClient(ctrl *gomock.Controller) *GoMockClient {
	mock := &GoMockClient{ctrl: ctrl}
	mock.recorder = &GoMockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *GoMockClient) EXPECT() *GoMockClientMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *GoMockClient) Delete(ctx context.Context, object string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, object)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *GoMockClientMockRecorder) Delete(ctx, object interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*GoMockClient)(nil).Delete), ctx, object)
}

// Upload mocks base method.
func (m *GoMockClient) Upload(ctx context.Context, file io.Reader, object string, expires time.Time) (*CloudFile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", ctx, file, object, expires)
	ret0, _ := ret[0].(*CloudFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload.
func (mr *GoMockClientMockRecorder) Upload(ctx, file, object, expires interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*GoMockClient)(nil).Upload), ctx, file, object, expires)
}