!!! note "Certificate input"

    The argument must be a PEM-encoded x.509 certificate string.

```
# decode a PEM certificate and access its subject
x509_decode('-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----').subject.common_name
```

```
# decode a base64-encoded certificate using base64_decode
x509_decode(base64_decode('LS0tLS1CRUdJTi...')).not_after
```
