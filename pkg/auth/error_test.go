package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDomain = "test"

type testError struct {
	src []byte
}

func (e *testError) Error() string {
	return "reconnect not be allowed code -1016"
}
func (e *testError) Domain() string {
	return testDomain
}
func (e *testError) Encoding() (string, error) {
	return base64.StdEncoding.EncodeToString(e.src), nil
}

func decoding(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func TestUnregisterDecoder(t *testing.T) {
	text := "hello world"
	err := &testError{
		src: []byte(text),
	}
	gErr := ToGRPCErr(err)
	metadata, e := ParseGRPCErr(gErr)
	assert.Nil(t, metadata)
	assert.EqualError(t, e, fmt.Sprintf("unknown decoder type %s ", testDomain))
}

func TestToGRPCErrAndParse(t *testing.T) {
	RegisterErrorDecoder(testDomain, decoding)
	text := "hello world"
	err := &testError{
		src: []byte(text),
	}
	gErr := ToGRPCErr(err)
	metadata, e := ParseGRPCErr(gErr)
	assert.Nil(t, e)
	assert.EqualValues(t, []byte(text), metadata)
}

func TestUnknownError(t *testing.T) {
	RegisterErrorDecoder(testDomain, decoding)
	err := errors.New("test unknown error")
	gErr := ToGRPCErr(err)
	metadata, e := ParseGRPCErr(gErr)
	assert.Nil(t, metadata)
	assert.Equal(t, err, e)
}
