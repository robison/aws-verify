package main

import "flag"
import "log"
import "os"
import "strings"

/**
 * Helper to panic on a returned error during setup
 */
func fatal(err error) {
	if err != nil {
		log.Panic(err)
	}
}

/**
 * Parse the certificates command flag and add to the handler
 */
func loadCertificates(flag string, handler *Verifier) (err error) {
	// Use the default Amazon certificate if none are specified in command arguments
	if len(flag) == 0 {
		log.Printf("Using default Amazon AWS signing certificate")

		_, err = handler.AddPEMCertificate(AmazonAWSCloudSigner)
		return err
	}

	// Load certificates from specified files
	for _, path := range strings.Split(flag, ",") {
		if len(path) == 0 {
			continue
		}

		_, err = handler.ReadPEMCertificate(path)

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var certs = flag.String("certificates", "", "Comma-separated list of paths to certificates used to verify signatures")
	var socket = flag.String("socket", "./aws-verify.sock", "UNIX socket that verifier will listen on")
	var mode = flag.Uint("mode", 0700, "File mode for listener socket")

	flag.Parse()

	handler := CreateVerifier()

	fatal(
		loadCertificates(*certs, handler),
	)

	fatal(
		CreateServer(*socket, os.FileMode(*mode), handler).Listen(),
	)
}
