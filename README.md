# toolsop
Operator for the deployed toolset in a kubernetes cluster

# local

To decrypt the secretdata to run it locally:

```bash
gopass caos-secrets/technical/toolsop/ansible-vault > ansible-vault-secret && \
ansible-vault decrypt --vault-password-file ansible-vault-secret config/manager/secret/id_rsa-toolsop-tools-read && \
rm ansible-vault-secret
```

To encrypt it again:

```bash
gopass caos-secrets/technical/toolsop/ansible-vault > ansible-vault-secret && \
ansible-vault encrypt --vault-password-file ansible-vault-secret config/manager/secret/id_rsa-toolsop-tools-read && \
rm ansible-vault-secret
```

To build it:

```bash
docker build --build-arg ANSIBLEVAULT_SECRET=$(gopass caos-secrets/technical/toolsop/ansible-vault) -t controller:latest .
```
