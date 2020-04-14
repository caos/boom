#!/usr/bin/env bash
# getter for gopass and secret yaml creation
set -e
[[ `uname` = "Linux" ]] && ENCODE="base64 --wrap=0" || ENCODE="base64"

# apply via: secrets.yaml.sh | kubectl apply -f -

GITHUB_IMAGE_PULL_SECRET=$(cat ~/.docker/config.json | $ENCODE)
ORBCONFIG=$(cat ~/.orb/config | $ENCODE)

NAMESPACE=caos-system

cat <<EOL
---
apiVersion: v1
data:
  .dockerconfigjson: $GITHUB_IMAGE_PULL_SECRET
kind: Secret
metadata:
  name: local-docker-login
  namespace: $NAMESPACE
type: kubernetes.io/dockerconfigjson
---
apiVersion: v1
data:
  orbconfig: $ORBCONFIG
kind: Secret
metadata:
  name: caos
  namespace: $NAMESPACE
type: Opaque
---
EOL
