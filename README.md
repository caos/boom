# boom: the base tooling operator

![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)
![release](https://github.com/caos/boom/workflows/Release/badge.svg)
[![GitHub license](https://img.shields.io/github/license/caos/boom)](https://github.com/caos/boom/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/caos/boom)](https://GitHub.com/caos/boom/releases/)

> This project is in alpha state. The API will continue breaking until version 1.0.0 is released

## What is it

`boom` is designed to ensure that someone can create a reproducable "platform" with tools which are tested for their interoperability.

Currently we include the following tools:

- Ambassador Edge Stack
- Prometheus Operator
- Grafana
- logging-operator
- kube-state-metrics
- prometheus-node-exporter
- loki
- ArgoCD

Upcoming tools:

- Flux

## How does it work

The operator works by reading a configuration (crd) located in a GIT Repository. Alternativly this `crd` can be read from the k8s api.
In our default setup our "cluster lifecycle" tool `orbiter`, shares the repository and secrets with `boom`. This because `orbiter` deploys `boom` in a newly created `k8s` cluster.

```yaml
apiVersion: boom.caos.ch/v1beta1
kind: Toolset
metadata:
  name: caos
  namespace: caos-system
spec:
  preApply:
    deploy: true
    folder: preapply
  postApply:
    deploy: true
    folder: postapply
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
    service:
      type: LoadBalancer
  kube-state-metrics:
    deploy: true
  argocd:
    deploy: false
    customImage:
      enabled: false
      imagePullSecret: github-image
      gopassGPGKey: "gpg"
      gopassSSHKey: "ssh"
      gopassStores:
      - directory: "directory"
        storeName: "store"
  prometheus:
    deploy: true
    storage:
      size: 5Gi
      storageClass: standard
  loki:
    deploy: true
    storage:
      size: 5Gi
      storageClass: standard
```

## How to use it

> Due to the github restriciton that even public images need to be authenticated, you need to make sure that you have `pull secret`. The used `personal access token` has to have the `repo` and `read:packages` permissions.

```bash
kubectl -n caos-system create secret docker-registry boomregistry --docker-server=docker.pkg.github.com --docker-username=${GITHUB_USERNAME} --docker-password=${GITHUB_ACCESS_TOKEN}
```

### GitOps Mode

#### Demo with a public crd repository

To easy test the example we have created a `demo crd repo`, located here [demo-orbiter-boom](https://github.com/caos/demo-orbiter-boom). It holds a `boom.yml` which can be applied to your cluster.

Apply `Boom` to your cluster:

```bash
kustomize build examples/gitops/publicrepo | kubectl apply -f -
```

#### Example with a private repository

Your first have to create an ssh-key which is added as deploy key to your git repository and then save the private key as secret in examples/gitops/privaterepo/secret.
Change the name of the key in the examples/gitops/privaterepo/kustomization.yaml with the filename of the saved key.

Apply `Boom` to your cluster:

```bash
kustomize build examples/gitops/privaterepo | kubectl apply -f -
```

#### k8s API Mode

example coming soon

## License

As usual Apache-2.0 see [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
