#!/bin/bash

gopass sync -s caos-secrets
./scripts/0_imagepull-secrets.sh | kubectl apply -f -
./scripts/1_argocd-secrets.sh | kubectl apply -f -

skaffold run -f ../build/skaffold/skaffold.yaml