// Package plugin/service provides a common way for host and
// plugins to communicate in a IPC-like way.
//
// In this framework, the host serves as the server while the
// plugins serves as clients. The plugins requests services
// provided by server through namespace, name and arguments.
// To simplify the problem, all data will be marshaled and
// transferred in json.
//
// The host-plugin communication will be multiplexed over a
// pair of pipes which are identified by URLs. The scheme of
// the URL identifies the kind of opener it is, while the host
// and path portion specifies how to reach the file. And the
// communiction pipes can be named pipes, sockets and even
// shared memories, as long as valid opener is defined for it.
package service

import (
	"encoding/json"
	"reflect"
)

var typeError = reflect.TypeOf((*error)(nil)).Elem()

// serviceType identifies the possible fields that could appear
// as the serviceRequest.Type field.
type serviceType string

const (
	serviceTypeCall         = serviceType("call")
	serviceTypeHasNamespace = serviceType("hasNamespace")
	serviceTypeGetManifest  = serviceType("getManifest")
	serviceTypeListServices = serviceType("listServices")
)

// Service is service function's general interface.
//
// Each service will be a function with multiple arguments and
// results. While running, the arguments will be attempted to
// be marshaled into json, while the results will be marshaled
// from json.
//
// There's a special restriction that the service provider
// should not use error interface as one of the arguments or
// results, while the service consumer must have the error
// interface as the last argument to receive interface error.
//
// Say, we have a service whose prototype looks like
// func(A, B) C, and thus its service consumer must be in the
// form of func(A, B) (C, error).
type Service interface{}

type serviceRequest struct {
	Sequence  uint64          `json:"sequence"`
	Type      serviceType     `json:"type"`
	Namespace string          `json:"namespace"`
	Name      string          `json:"name"`
	Args      json.RawMessage `json:"args"`
}

type serviceResponse struct {
	Sequence uint64 `json:"sequence"`

	// Services used in serviceTypeListServices.
	Services []string `json:"services,omitempty"`

	// Ok used in serviceTypeHasNamespace.
	Ok bool `json:"ok,omitempty"`

	// Reply used in serviceTypeCall and serviceTypeGetManifest.
	Reply json.RawMessage `json:"reply,omitempty"`

	// ErrMsg stores API errors only, and user should define
	// their own error type as a marshalable reply type.
	ErrMsg *string `json:"error,omitempty"`
}
