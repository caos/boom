# boom: the base tooling operator

![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)
![release](https://github.com/caos/boom/workflows/Release/badge.svg)
[![GitHub license](https://img.shields.io/github/license/caos/boom)](https://github.com/caos/boom/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/caos/boom)](https://GitHub.com/caos/boom/releases/)

> This project is in alpha state. The API will continue breaking until version 1.0.0 is released

## What is it

`boom` is designed to ensure that someone can create a reproducable "platform" with tools which are tested for their interoperability.

Currently we include the following tools:

- Ambassador
- Cert-Manager
  - Will be removed when the `Ambassador Edge Stack` reaches GA
- Prometheus Operator
- Grafana
- logging-operator
- kube-state-metrics
- prometheus-node-exporter

Upcoming tools:

- Loki
- ArgoCD
- Flux

## How does it work

The operator works by reading a configuration (crd) located in a GIT Repository. Alternativly this `crd` can be read from the k8s api.
In our default setup our "cluster lifecycle" tool `orbiter`, shares the repository and secrets with `boom`. This because `orbiter` deploys `boom` in a newly creadted `k8s` cluster.

```yaml
---
apiVersion: boom.caos.ch/v1beta1
kind: Toolset
metadata:
  name: caos
  namespace: caos-system
spec:
  kubeVersion: v1.17.0
  prometheus-operator:
    deploy: true
  logging-operator:
    deploy: true
  prometheus-node-exporter:
    deploy: true
  grafana:
    deploy: true
  ambassador:
    deploy: true
  kube-state-metrics:
    deploy: true
  argocd:
    deploy: true
  prometheus:
    deploy: true
```

## How to use it

> Due to the github restriciton that even public images need to be authenticated, you need to make sure that you have `pull secret`

### GitOps

ImagePullSecret create for github-packages as basic-auth is required.
The used personal access token has to have the "repo" and "read:packages" permissions. 
```bash
kubectl -n caos-system create secret docker-registry boomregistry --docker-server=docker.pkg.github.com --docker-username=${GITHUB_USERNAME} --docker-password=${GITHUB_ACCESS_TOKEN}
```

#### example with a public repository:
```bash
kustomize build examples/gitops/publicrepo | kubectl apply -f -
```

#### example with a private repository:

Your first have to create an ssh-key with is added as deploy key to your git repository and then save the private key as secret in examples/gitops/privaterepo/secret.
Change the name of the key in the examples/gitops/privaterepo/kustomization.yaml with the filename of the saved key.

```bash
kustomize build examples/gitops/publicrepo | kubectl apply -f -
```


### k8s API

example coming soon

## License

As usual Apache-2.0 see [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

## Inspiration

### Name

The inspiration for the name `boom` originates from the `Orbiter Boom Sensor System` used by `NASA`. `OBSS` or in short `booom` is a package of tools which are used to inspect the `orbiter` for damages in the thermal shielding. As our application automation tool is called `orbiter` we thaught this is a name which fits great.

More information regarding the `boom` can be found here ![Wikipedia OBSS](https://en.wikipedia.org/wiki/Orbiter_Boom_Sensor_System)

Our project `orbiter` is located here ![CAOS orbiter](https://github.com/caos/orbiter)

### Further Inspiration

Thanks to `rancher` with `rio` as well for the inspiration.
And also to ex. `coreos` with `tectonic`.

Both platform influenced this project.
