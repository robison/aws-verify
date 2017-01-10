package main

import "flag"
import "log"
import "strings"

/**
 * Helper to panic on a returned error during setup
 */
func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var socket = flag.String("socket", "./aws-verify.sock", "UNIX socket that verifier will listen on")
	var certificates = flag.String("certificates", "", "Comma-separated list of paths to certificates used to verify signatures")
	flag.Parse()

	handler := CreateVerifier()

	// Use the default Amazon certificate if none are specified in command arguments
	if len(*certificates) == 0 {
		log.Printf("Using default Amazon AWS signing certificate")
		_, err := handler.AddPEMCertificate(AMAZON_PUBLIC_CLOUD)

		fatal(err)
	}

	// Load certificates from specified file
	for _, path := range strings.Split(*certificates, ",") {
		if len(path) > 0 {
			_, err := handler.ReadPEMCertificate(path)

			fatal(err)
		}
	}

	CreateServer(*socket, handler).Listen()
}
