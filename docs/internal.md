# Internal logic

## Folder structure

The boom will extend the existing tools folder with different subfolders for the necessary applications.
For each crd there will be an subfolder for each application where the generated values.yaml and kustomization.yaml are stored.
As a result of this files there will be an results.yaml under the subfolder results/*crd-name*/.

Like this:

* tools
  * logging-operator
    * *crd-name*
      * templator.yaml
      * kustomization.yaml
      * values.yaml
    * results
      * *crd-name*
        * results.yaml
  * start.sh
  * fetch-all.sh
  * *helm*
  * *charts*
  * kustomize

also are there the differnt tools for templating, the charts folder consists of all fetchet charts localy, the kustomize folder has the necessary shell scripts for the templators and the helm folder is the helm-home folder.
The charts will get fetched during the docker build phase with running of the fetch-all.sh.

To start the different steps:

```bash
# fetch chart for local
./fetch-all.sh *toolset*
# template chart with values.yaml
./start.sh *application* *crd-name*
# apply results to cluster
kubectl apply -f *application*/results/*crd-name*/results.yaml
```

## toolsets

To add any new toolset or change existing ones look into the toolsets folder.
The structure in this folder is *important* as it is as follows:

* tools
  * toolsets
    * *toolset-name*
      * *application-name*.yaml

It is *important* as the boom has logic which works over this structure to build the knowledge which toolsets are existing and out of which applications do they consist.

## used tools

The following cli-tools are used from the boom:

* helm
* kubectl
* kustomize

As they are used, they also have to be installed into the image during the docker build.

## To let it run

### locally

Before you can run locally you have to fetch all charts:

```bash
./tools/fetch-all.sh *toolset*
```

To decrypt the secretdata to run it locally:

```bash
gopass caos-secrets/technical/boom/ansible-vault > ansible-vault-secret && \
ansible-vault decrypt --vault-password-file ansible-vault-secret config/manager/secret/id_rsa-boom-tools-read && \
rm ansible-vault-secret
```

To encrypt it again:

```bash
gopass caos-secrets/technical/boom/ansible-vault > ansible-vault-secret && \
ansible-vault encrypt --vault-password-file ansible-vault-secret config/manager/secret/id_rsa-boom-tools-read && \
rm ansible-vault-secret
```

To build it:

```bash
docker build -t controller:latest .
```

### cluster

To deploy the boom to a cluster:

```bash
cd config/manager && kustomize edit set image controller=docker.pkg.github.com/caos/boom/boom:latest && cd ../..
kustomize build config/default | kubectl apply -f -
```
