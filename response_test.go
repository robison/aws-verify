package main

import "github.com/stretchr/testify/assert"
import "testing"

import "encoding/json"
import "errors"
import "net/http"
import "net/http/httptest"

func TestNewResponse(t *testing.T) {
  response := NewResponse(http.StatusOK, true)

  assert.IsType(t, new(Response), response, "Returns an instance of Response")
  assert.Equal(t, http.StatusOK, response.Code, "Sets the Code property")
  assert.True(t, response.Success, "Sets the Success property")
}

func TestAddDocument(t *testing.T) {
  response := NewResponse(http.StatusOK, true)
  response.AddDocument([]byte(`{"valid": "json", "string":["meh"]}`))

  assert.Equal(t, "json", response.Document["valid"], "Parses JSON into a map")
}

func TestAddError(t *testing.T) {
  response := NewResponse(http.StatusBadRequest, false)
  response.AddError(errors.New("This is an error"))

  assert.Equal(t, "This is an error", response.Error, "Sets Error property to error message")
}

func TestSend(t *testing.T) {
  response := NewResponse(http.StatusOK, true)
  w := httptest.NewRecorder()

  response.AddDocument([]byte(`{"valid": "json", "string":["meh"]}`))

  err := response.Send(w)
  body := new(Response)

  assert.Nil(t, json.NewDecoder(w.Body).Decode(body))

  assert.Nil(t, err, "Does not return an error")
  assert.Equal(t, http.StatusOK, w.Code, "Sets the correct status code")
  assert.Equal(t, http.StatusOK, body.Code, "Sets the code property in the response")
  assert.True(t, body.Success, "Sets the success property in the response")
  assert.Equal(t, "json", response.Document["valid"], "Marshals the signed document back into the response body")
}
