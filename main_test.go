package main

import "github.com/stretchr/testify/assert"
import "testing"

import "errors"

func TestFatal(t *testing.T) {
	defer func() {
		recovered := recover()
		assert.Equal(t, "This is an error", recovered, "Panics with an error")
	}()

	fatal(errors.New("This is an error"))
}

func TestLoadCertificates(t *testing.T) {
  handler := CreateVerifier()
  err := loadCertificates("./testdata/fake-certificate.pem,./testdata/fake-certificate.pem,", handler)

  assert.Nil(t, err, "Does not return an error")
  assert.Equal(t, 2, len(handler.certificates), "Adds certificates to the verifier")

  err = loadCertificates("", handler)
  assert.Equal(t, 3, len(handler.certificates), "Adds default certificate to the verifier")
}
