package main

import "errors"
import "io/ioutil"
import "net/http"

import "encoding/pem"
import "crypto/x509"

import "github.com/fullsailor/pkcs7"

// NewRequest instantiates a Request and reads the body of the provided
// http.Request instance
func NewRequest(r *http.Request) (*Request, error) {
	request := &Request{}

	return request, request.Read(r)
}

// Request encapsulates an incoming HTTP request's signature and the resulting PKCS7 object
type Request struct {
	P7        *pkcs7.PKCS7
	Signature []byte
}

// Read loads the incoming request's body into the Signature property
func (request *Request) Read(r *http.Request) (err error) {
	body, err := ioutil.ReadAll(r.Body)

	if err == nil {
		request.Signature = body
	}

	return err
}

// Parse converts the Request's PEM encoded Signature property into a PKCS7 instance
// and attaches a set of x509 signing candidates for later verification
func (request *Request) Parse(certificates []*x509.Certificate) error {
	// Extract ANS1 from the PEM block
	block, _ := pem.Decode(request.Signature)
	if block == nil {
		return errors.New("invalid PEM data. Unable to parse certificate")
	}

	// Parse the ASN1 object into a PKCS7 instance
	p7, err := pkcs7.Parse(block.Bytes)

	// Attach a certificate set for later verification
	if err == nil {
		p7.Certificates = certificates
		request.P7 = p7
	}

	return err
}

// Verify the signature's signer integrity against signing candidates and
// validate the signed content against the provided request document
func (request *Request) Verify() error {
	return request.P7.Verify()
}
