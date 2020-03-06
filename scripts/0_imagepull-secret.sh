#!/usr/bin/env bash
set -e
[[ `uname` = "Linux" ]] && ENCODE="base64 --wrap=0" || ENCODE="base64"

# argocd gpg key secrets
IMAGEPULL_KEY=$(gopass caos-secrets/technical/boom/git-read-secret )
NAMESPACE="caos-system"

cat <<EOL
apiVersion: v1
data:
  id_rsa: $IMAGEPULL_KEY
kind: Secret
metadata:
  name: privaterepo-secret
  namespace: $NAMESPACE
EOL
