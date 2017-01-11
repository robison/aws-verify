package main

import "github.com/stretchr/testify/assert"
import "testing"

import "errors"
import "os"

import "net/http"
import "net/http/httptest"

import "crypto/x509"
import "encoding/json"

func TestCreateVerifier(t *testing.T) {
	verify := CreateVerifier()

	assert.IsType(t, new(Verifier), verify, "CreateVerifier creates a Verifier")
}

func TestAddPEMCertificate(t *testing.T) {
	verify := CreateVerifier()
	certificate, err := verify.AddPEMCertificate(AmazonAWSCloudSigner)

	assert.Equal(t, 1, len(verify.certificates), "Adds a certificate to the verifier")
	assert.Equal(t, "Amazon Web Services LLC", certificate.Subject.Organization[0], "Read the correct certificate")

	certificate, err = verify.AddPEMCertificate([]byte("definitely not a PEM block"))

	assert.NotNil(t, err, "Returns an error if the argument isn't a valid PEM object")
}

func TestReadPEMCertificate(t *testing.T) {
	verify := CreateVerifier()
	certificate, err := verify.ReadPEMCertificate("./testdata/fake-certificate.pem")

	assert.Nil(t, err, "Does not return an error")
	assert.IsType(t, new(x509.Certificate), certificate, "Returns a new instance of x509.Certificate")
	assert.Equal(t, 1, len(verify.certificates), "Adds certificate to the verifier")
	assert.Equal(t, "Internet Widgits Pty Ltd", certificate.Subject.Organization[0], "Read the correct certificate")

	certificate, err = verify.ReadPEMCertificate("./testdata/not-a-real-file")

	assert.NotNil(t, err, "Returns an error if the file does not exist")

	certificate, err = verify.ReadPEMCertificate("./testdata/valid-signature.pem")

	assert.NotNil(t, err, "Returns an error if the file is not an x509 certificate")
}

func TestOK(t *testing.T) {
	verify := CreateVerifier()
	w := httptest.NewRecorder()
	err := errors.New("This is a test error")

	result := verify.OK(err, http.StatusBadRequest, w)
	body := new(Response)

	assert.False(t, result, "Returns false when an error is handled")
	assert.Equal(t, http.StatusBadRequest, w.Code, "Sets the response status code")

	assert.Nil(t, json.NewDecoder(w.Body).Decode(body))
	assert.Equal(t, "1.0", body.V, "Sets a version string in the response body")
	assert.Equal(t, http.StatusBadRequest, body.Code, "Sets code parameter to response code")
	assert.False(t, body.Success, "Sets success parameter to true")
	assert.Equal(t, err.Error(), body.Error, "Sets error parameter to error message")

	w = httptest.NewRecorder()
	result = verify.OK(nil, http.StatusBadRequest, w)

	assert.True(t, result, "Returns true when argument is nil")
	assert.Equal(t, http.StatusOK, w.Code, "Does not set the status code")
	assert.Equal(t, 0, w.Body.Len(), "Does not write a response body")
}

func TestServeHTTP(t *testing.T) {
	verify := CreateVerifier()
	_, err := verify.AddPEMCertificate(AmazonAWSCloudSigner)
	assert.Nil(t, err)

	signature, err := os.Open("./testdata/valid-signature.pem")
	assert.Nil(t, err)
	defer signature.Close()

	r := httptest.NewRequest("POST", "/", signature)
	w := httptest.NewRecorder()
	body := new(Response)

	verify.ServeHTTP(w, r)

	assert.Nil(t, json.NewDecoder(w.Body).Decode(body))
	assert.Equal(t, http.StatusOK, w.Code, "Sets the response status code")

	assert.Equal(t, "1.0", body.V, "Sets a version string in the response body")
	assert.Equal(t, http.StatusOK, body.Code, "Sets code parameter to response code")
	assert.True(t, body.Success, "Sets success parameter to false")
	assert.Empty(t, body.Error, "Sets error parameter to error message")
}
