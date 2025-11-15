# PostgreSQL:
## SSL - verify-full with mutual TLS
- Need a RootCA cert
- Generate Postgres Server cert signed by RootCA cert:
```bash
# Server key
openssl genrsa -out server.key 4096

# Server CSR (Certificate Signing Request)
openssl req -new -key server.key -out server.csr -subj "/CN=postgres"

# Sign server CSR with CA
openssl x509 -req -in server.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out server.crt -days 500 -sha256
```