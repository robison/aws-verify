package main

import "github.com/stretchr/testify/assert"
import "testing"

import "net/http"
import "os"
import "time"

func TestCreateServer(t *testing.T) {
  handler := CreateVerifier()
  server := CreateServer("./testdata/listener.sock", 0700, handler)

  assert.IsType(t, new(Server), server, "Returns an instance of Server")
  assert.IsType(t, new(http.Server), server.server, "Creates an http.Server instance")
  assert.Equal(t, handler, server.server.Handler, "Sets the http.Server's Handler")
  assert.Equal(t, "./testdata/listener.sock", server.socket, "Sets the socket path property")
  assert.Equal(t, os.FileMode(0700), server.mode, "Sets the file mode property")
}

func TestListenAndClose(t *testing.T) {
  server := CreateServer("./testdata/listener.sock", 0700, CreateVerifier())

  go func() {
    // Wait a bit for the listener to initialize
    time.Sleep(time.Millisecond * 250)

    info, err := os.Stat("./testdata/listener.sock")
    assert.Nil(t, err, "Does not return an error")
    assert.Equal(t, os.FileMode(0700), info.Mode() & 0xFFF, "Opens socket with the correct file-mode")

    server.Close()
  }()

  err := server.Listen()
  assert.Nil(t, err, "Does not return an error")

  _, err = os.Stat("./testdata/listener.sock")
  assert.True(t, os.IsNotExist(err), "Cleans up socket handle after closing")

  server = CreateServer("./testdata-foo-does-not-exist/listener.sock", 0700, CreateVerifier())
  err = server.Listen()

  assert.NotNil(t, err, "Fails to open a socket in a non-existent directory")
}
