AWS EC2 Identity PKCS7 Verifier
===============================

A minimal service for verifying and extracting EC2 instance identity data from PKCS7 objects.

All EC2 instances have access to a signed, PEM-encoded PKCS7 object containing their unique identifiers via the EC2 meta-data service: http://169.254.169.254/latest/dynamic/instance-identity/pkcs7. This endpoint is unique for every instance and provides them with a means to cryptographically prove their identity, as EC2 instances, to non-AWS resources (e.g things that aren't AWS-IAM integrated) that are concerned with such things. [The Docs](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html) explain how fetch and prepare an instance's PKCS7 signature object, and how one can use the `openssl` CLI to validate it.

## The Problem

Unfortunately, PKCS7 is a somewhat esoteric ANS1 structure these days. OpenSSL implements it, and some languages provide native bindings to OpenSSL's interface, but many simply have no support for PKCS7, or achieve it by shelling out to the `openssl` CLI. Even for languages that use OpenSSL's C implementation, either via bindings or shell-outs, the problem remains that those code-lines in OpenSSL are not likely heavily utilized or reviewed/tested for the usual suite of vulnerabilities common to C memory management.

`aws-verify` is a tiny HTTP server that validates EC2 PKCS7 documents and returns their now-verified contents to the client. It uses https://github.com/fullsailor/pkcs7, a pure-golang implementation of PKCS7 on top of Go's own `crypto` and `encoding/ans1` packages. It accepts requests with a PEM-encoded PKCS7 object as the body, and responds with a JSON object indicating whether the signature is valid and authenticated by one of the configured signing-certificates, and if valid, includes the JSON document that was originally signed.

### Tradeoff

Using `fullsailor/pkcs7` removes direct dependencies upon OpenSSL and C memory-management, but introduces a new, lightly-adopted code-line. The tradeoff is that, while there likely aren't any low-level memory-access violations, the number of eyes verifying the correctness of the implementation is small. We have to rely upon positive and negative testing herein to be reasonably confident that the library is behaving as expected.

## Interface

A PEM-encoded PKCS7 object can be sent using any method to any path. The server passes all requests to the verification handler:

```
POST / HTTP/1.1
Content-Length: 1088

-----BEGIN PKCS7-----
MIAGCSqGSIb3DQEHAqCAMIACAQExCzAJBgUrDgMCGgUAMIAGCSqGSIb3DQEHAaCAJIAEggGoewog
ICJwcml2YXRlSXAiIDogIjE3Mi4xOC4xMzAuMjE3IiwKICAiZGV2cGF5UHJvZHVjdENvZGVzIiA6
IG51bGwsCiAgImF2YWlsYWJpbGl0eVpvbmUiIDogInVzLWVhc3QtMWEiLAogICJ2ZXJzaW9uIiA6
ICIyMDEwLTA4LTMxIiwKICAiaW5zdGFuY2VJZCIgOiAiaS1iMzU1ZTdmNiIsCiAgImJpbGxpbmdQ
cm9kdWN0cyIgOiBudWxsLAogICJpbnN0YW5jZVR5cGUiIDogInQyLnNtYWxsIiwKICAiYWNjb3Vu
dElkIiA6ICIwNDIyOTM5NjQzODEiLAogICJpbWFnZUlkIiA6ICJhbWktNjk2MmEzMDQiLAogICJw
ZW5kaW5nVGltZSIgOiAiMjAxNi0wNi0yMVQxNzoyODoyM1oiLAogICJhcmNoaXRlY3R1cmUiIDog
Ing4Nl82NCIsCiAgImtlcm5lbElkIiA6IG51bGwsCiAgInJhbWRpc2tJZCIgOiBudWxsLAogICJy
ZWdpb24iIDogInVzLWVhc3QtMSIKfQAAAAAAADGCARcwggETAgEBMGkwXDELMAkGA1UEBhMCVVMx
GTAXBgNVBAgTEFdhc2hpbmd0b24gU3RhdGUxEDAOBgNVBAcTB1NlYXR0bGUxIDAeBgNVBAoTF0Ft
YXpvbiBXZWIgU2VydmljZXMgTExDAgkAlrpI2eVeGmcwCQYFKw4DAhoFAKBdMBgGCSqGSIb3DQEJ
AzELBgkqhkiG9w0BBwEwHAYJKoZIhvcNAQkFMQ8XDTE2MDYyMTE3MjgzNlowIwYJKoZIhvcNAQkE
MRYEFDv689ev9in7oviZYY/W5wz5rZf6MAkGByqGSM44BAMELjAsAhRug2rKGCmxa4BuYKo3UBff
Gu+CMQIUbEchq2hiuRajESVAdi2xvCaCmDsAAAAAAAA=
-----END PKCS7-----

```

For a valid request, the server responds with `success: true` and a `document` key containing the parsed JSON included in the PKCS7 object:

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 06 Jan 2017 16:34:41 GMT
Content-Length: 414

{
  "v": "1.0",
  "code": 200,
  "success": true,
  "error": "",
  "document": {
    "accountId": "042293964381",
    "architecture": "x86_64",
    "availabilityZone": "us-east-1a",
    "billingProducts": null,
    "devpayProductCodes": null,
    "imageId": "ami-6962a304",
    "instanceId": "i-b355e7f6",
    "instanceType": "t2.small",
    "kernelId": null,
    "pendingTime": "2016-06-21T17:28:23Z",
    "privateIp": "172.18.130.217",
    "ramdiskId": null,
    "region": "us-east-1",
    "version": "2010-08-31"
  }
}
```

An invalid request will receive a response `success: false` and an `error` key with some explanation of the failure:

```
{
  "v": "1.0",
  "code": 403,
  "success": false,
  "error": "pkcs7: No certificate for signer",
  "document": null
}
```

## Configuration

The server is configured via CLI flags:

```
$ ./aws-verify -h
Usage of ./aws-verify:
  -certificates string
    	Comma-separated list of paths to certificates used to verify signatures
  -socket string
    	UNIX socket that verifier will listen on (default "./aws-verify.sock")
```

By default, the server will use Amazon AWS's public key, published at http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html, to verify signatures.


The server only listened on UNIX domain sockets. It's not meant to be a general-purpose interface, but rather a hidden-backend consumed by some other API service.

## Caveats

The `aws-verify` interface can be configured to handle any PKCS7 objects from any signer via the `certificates` flag, **but** the signed data must be a valid JSON string.
