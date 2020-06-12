# Snippetbox

## Creating Certificates

We want to serve files via https also for devlopment and therefore we need to
create a self signed certificate. Luckily Go comes with a tool for that out
of the box which we can use. In order to create these certicficates, create a
tls directory in the root and execute the following command:

```bash
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

If you've installed Go via homebrew, you can find the `generate_cert.go` file under
`/usr/local/Cellar/go/<version>/libexec/src/crypto/tls/generate_cert.go`
