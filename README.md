# TLS Client Certificate Testing Server
A Golang HTTPS server to test client certificates.
It prints the details (issuer and subject DN, chain in PEM format) of the client certificate to stdout.

## Usage

```shell
tls-clientcert-tester <address> <optional server certificate PEM file>
```

Arguments:
- `<address>` is the address to listen on, e.g. `:8080`.
- `<optional server certificate PEM file>` is the path to a PEM file containing the server certificate.
  If not provided, the server will use a self-signed certificate matching the hostname.

Using Docker: 

```shell
docker run -p 8080:8080 reinkrul/tls-clientcert-test:latest :8080
```

or

```shell
docker run -p 8080:8080 -v certificate.pem:/cert.pem reinkrul/tls-clientcert-test:latest :8080 /cert.pem
```

## Building

```shell
docker buildx build --platform linux/amd64,linux/arm64 -t reinkrul/tls-clientcert-test:latest --push .
```
