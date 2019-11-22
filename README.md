# toolsop
Operator for the deployed toolset in a kubernetes cluster

# local

To decrypt the secretdata to run it locally:

```bash
gopass caos-secrets/technical/toolsop/ansible-vault > ansible-vault-secret && \
ansible-vault decrypt --vault-password-file ansible-vault-secret config/manager/secret/id_rsa-toolsop-tools-read && \
rm ansible-vault-secret
```

To encrypt it again:

```bash
gopass caos-secrets/technical/toolsop/ansible-vault > ansible-vault-secret && \
ansible-vault encrypt --vault-password-file ansible-vault-secret config/manager/secret/id_rsa-toolsop-tools-read && \
rm ansible-vault-secret
```

To build it:

```bash
docker build --build-arg ANSIBLEVAULT_SECRET=$(gopass caos-secrets/technical/toolsop/ansible-vault) -t controller:latest .
```


# internal logic

## folder structure

The toolsop will extend the existing tools folder with different subfolders for the necessary applications.
Each application will have 3 folders fetchers, templators and results. The fetchers will fetch the necessary chart for the local availability, the templators will use de generated values.yaml to compress all resulting yamls into the results/results.yaml file which gets applied into the clsuter.

Like this:

* tools
  * logging-operator
    * fetcherts
      * *crd-name*
        * fetcher.yaml
        * kustomization.yaml
    * templators
      * *crd-name*
        * templator.yaml
        * kustomization.yaml
        * values.yaml
    * results
      * *crd-name*
        * results.yaml
  * start.sh
  * *helm*
  * *charts*
  * kustomize

also are there the differnt tools for templating, the charts folder consists of all fetchet charts localy, the kustomize folder has the necessary shell scripts for the templators and fetchers and the helm folder is the helm-home folder.

To start the different steps:

```bash
# fetch chart for local
./start.sh *application* fetchers/*crd-name*
# template chart with values.yaml
./start.sh *application* templators/*crd-name*
# apply results to cluster
kubectl apply -f *application*/results/*crd-name*/results.yaml
```

## toolsets

To add any new toolset or change existing ones look into the toolsets folder.
The structure in this folder is *important* as it is as follows:

* toolsets
  * *toolset-name*
    * *version*
      * *application-name*.yaml

It is *important* as the toolsop has logic which works over this structe to build the knowledge which toolsets are existing and out of which applications do they consist.

## used tools

The following cli-tools are used from the toolsop:

* helm
* kubectl
* kustomize

As they are used, they also have to be installed into the image during the docker build.