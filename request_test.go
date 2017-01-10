package main

import "github.com/stretchr/testify/assert"
import "testing"

import "os"
import "strings"
import "net/http/httptest"

import "github.com/fullsailor/pkcs7"

func TestNewRequest(t *testing.T) {
  signature, err := os.Open("./testdata/valid-signature.pem")
  assert.Nil(t, err)
  defer signature.Close()

  r := httptest.NewRequest("POST", "/", signature)
  request, err := NewRequest(r)

  assert.Nil(t, err, "Does not return an error")
  assert.NotNil(t, request.Signature, "Reads a PEM encoded signature from the HTTP request")
}

func TestParse(t *testing.T) {
  verify := CreateVerifier()
  _, err := verify.AddPEMCertificate(AMAZON_PUBLIC_CLOUD)
  assert.Nil(t, err)

  signature, err := os.Open("./testdata/valid-signature.pem")
  assert.Nil(t, err)
  defer signature.Close()

  r := httptest.NewRequest("POST", "/", signature)
  request, err := NewRequest(r)

  err = request.Parse(verify.certificates)

  assert.Nil(t, err, "Does not return an error")
  assert.IsType(t, new(pkcs7.PKCS7), request.P7, "Generates a new PKCS7 instance")
  assert.Equal(t, verify.certificates, request.P7.Certificates, "It sets PKCS7's the Certificates slice")

  r = httptest.NewRequest("POST", "/", strings.NewReader("Not a valid PEM string"))
  request, err = NewRequest(r)

  err = request.Parse(verify.certificates)
  assert.NotNil(t, err, "Returns an error when an invalid PEM string is received")
}

func TestVerify(t *testing.T) {
  verify := CreateVerifier()
  _, err := verify.AddPEMCertificate(AMAZON_PUBLIC_CLOUD)
  assert.Nil(t, err)

  signature, err := os.Open("./testdata/valid-signature.pem")
  assert.Nil(t, err)
  defer signature.Close()

  r := httptest.NewRequest("POST", "/", signature)
  request, err := NewRequest(r)

  err = request.Parse(verify.certificates)
  assert.Nil(t, err)

  err = request.Verify()
  assert.Nil(t, err, "Does not return an error for a valid signature")

  // TODO Needs a negative CA test (e.g. no signer for given signature)
  // and a valid content test (e.g. content does not match signature)
}
