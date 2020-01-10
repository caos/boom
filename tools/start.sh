#!/bin/bash

export TOOLS_HOME="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

function kustomizeIt {
  XDG_CONFIG_HOME=$TOOLS_HOME \
  kustomize build --enable_alpha_plugins \
    $TOOLS_HOME/$1/$2
}

kustomizeIt $1 $2