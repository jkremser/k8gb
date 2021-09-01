/*
Copyright 2021 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/
// Code generated by MockGen. DO NOT EDIT.
// Source: controllers/providers/dns/dns.go

// Package dns is a generated GoMock package.
package dns

import (
	reflect "reflect"

	v1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	gomock "github.com/golang/mock/gomock"
	endpoint "sigs.k8s.io/external-dns/endpoint"
)

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// CreateZoneDelegationForExternalDNS mocks base method.
func (m *MockProvider) CreateZoneDelegationForExternalDNS(arg0 *v1beta1.Gslb) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateZoneDelegationForExternalDNS", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateZoneDelegationForExternalDNS indicates an expected call of CreateZoneDelegationForExternalDNS.
func (mr *MockProviderMockRecorder) CreateZoneDelegationForExternalDNS(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateZoneDelegationForExternalDNS", reflect.TypeOf((*MockProvider)(nil).CreateZoneDelegationForExternalDNS), arg0)
}

// Finalize mocks base method.
func (m *MockProvider) Finalize(arg0 *v1beta1.Gslb) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Finalize", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Finalize indicates an expected call of Finalize.
func (mr *MockProviderMockRecorder) Finalize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finalize", reflect.TypeOf((*MockProvider)(nil).Finalize), arg0)
}

// GetExternalTargets mocks base method.
func (m *MockProvider) GetExternalTargets(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExternalTargets", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetExternalTargets indicates an expected call of GetExternalTargets.
func (mr *MockProviderMockRecorder) GetExternalTargets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExternalTargets", reflect.TypeOf((*MockProvider)(nil).GetExternalTargets), arg0)
}

// GslbIngressExposedIPs mocks base method.
func (m *MockProvider) GslbIngressExposedIPs(arg0 *v1beta1.Gslb) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GslbIngressExposedIPs", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GslbIngressExposedIPs indicates an expected call of GslbIngressExposedIPs.
func (mr *MockProviderMockRecorder) GslbIngressExposedIPs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GslbIngressExposedIPs", reflect.TypeOf((*MockProvider)(nil).GslbIngressExposedIPs), arg0)
}

// SaveDNSEndpoint mocks base method.
func (m *MockProvider) SaveDNSEndpoint(arg0 *v1beta1.Gslb, arg1 *endpoint.DNSEndpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDNSEndpoint", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveDNSEndpoint indicates an expected call of SaveDNSEndpoint.
func (mr *MockProviderMockRecorder) SaveDNSEndpoint(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDNSEndpoint", reflect.TypeOf((*MockProvider)(nil).SaveDNSEndpoint), arg0, arg1)
}

// String mocks base method.
func (m *MockProvider) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockProviderMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockProvider)(nil).String))
}
