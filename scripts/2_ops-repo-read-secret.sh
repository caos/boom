#!/usr/bin/env bash
set -e
[[ `uname` = "Linux" ]] && ENCODE="base64 --wrap=0" || ENCODE="base64"

# argocd gpg key secrets
BOOM_TEST_OPS_SSH=$(gopass caos-secrets/technical/boom/test/ops-ssh-private-key )
NAMESPACE="caos-system"

cat <<EOL
apiVersion: v1
data:
  ops-ssh: $BOOM_TEST_OPS_SSH
kind: Secret
metadata:
  name: boom-test-ops
  namespace: $NAMESPACE
EOL