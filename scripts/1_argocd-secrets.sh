#!/usr/bin/env bash
set -e
[[ `uname` = "Linux" ]] && ENCODE="base64 --wrap=0" || ENCODE="base64"

# argocd gpg key secrets
ARGOCD_GPG_KEY=$(gopass caos-secrets/technical/boom/test/gpg-private-key )
ARGOCD_SSH_KEY=$(gopass caos-secrets/technical/boom/test/ssh-private-key )
NAMESPACE="caos-system"

cat <<EOL
apiVersion: v1
data:
  gpg-boom-test: $ARGOCD_GPG_KEY
kind: Secret
metadata:
  name: argocd-gpg-keys
  namespace: $NAMESPACE
---
apiVersion: v1
data:
  ssh-boom-test: $ARGOCD_SSH_KEY
kind: Secret
metadata:
  name: argocd-ssh-config
  namespace: $NAMESPACE
EOL
