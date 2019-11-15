export TOOLS_HOME="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

function kustomizeIt {
  XDG_CONFIG_HOME=$TOOLS_HOME \
  kustomize build --enable_alpha_plugins \
    $TOOLS_HOME/$1/$2
}

case $1 in
    'logging-operator')
        kustomizeIt $1 $2
    ;;
    'prometheus-operator')
        kustomizeIt $1 $2
    ;;
    'prometheus-node-exporter')
        kustomizeIt $1 $2
    ;;
    'grafana')
        kustomizeIt $1 $2
    ;;
    'prometheus')
        kustomizeIt $1 $2
    ;;
    'cert-manager')
        kustomizeIt $1 $2
    ;;
    'ambassador')
        kustomizeIt $1 $2
    ;;
    *)
        echo "please select a valid tool:

            logging-operator
            prometheus-operator
            prometheus-node-exporter
            grafana
            prometheus
            cert-manager
            ambassador
        "
        exit 1
esac