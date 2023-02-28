package auth

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorDetail interface {
	Domain() string
	Encoding() (string, error)
}

func ToGRPCErr(err error) error {
	switch e := err.(type) {
	case ErrorDetail:
		var data string
		s := status.New(codes.Unauthenticated, "")
		data, err = e.Encoding()
		if err != nil {
			return err
		}
		s, err = s.WithDetails(&errdetails.ErrorInfo{
			Reason: "DEVICE_REJECT",
			Domain: e.Domain(),
			Metadata: map[string]string{
				"resp_data": data,
			},
		})
		if err != nil {
			return err
		}
		return s.Err()
	default:
		return err
	}
}

var errorDecoder = make(map[string]DecodeFunc)

type DecodeFunc func(src string) ([]byte, error)

func RegisterErrorDecoder(name string, exec DecodeFunc) {
	errorDecoder[name] = exec
}

func ParseGRPCErr(err error) ([]byte, error) {
	errStatus, ok := status.FromError(err)
	if !ok {
		return nil, err
	}
	if errStatus.Code() == codes.Unauthenticated {
		for _, detail := range errStatus.Details() {
			if d, ok := detail.(*errdetails.ErrorInfo); ok && d.GetReason() == "DEVICE_REJECT" {
				decoder, ok := errorDecoder[d.GetDomain()]
				if !ok || decoder == nil {
					return nil, fmt.Errorf("unknown decoder type %s ", d.GetDomain())
				}
				return decoder(d.Metadata["resp_data"])
			}
		}
	}
	return nil, err
}
