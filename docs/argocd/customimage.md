# Argocd Customimage

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `enabled`                          | Flag if custom argocd-image should get used with gopass                         | false                             |
| `imagePullSecret`                  | Name of used imagePullSecret to pull customImage                                |                                   |
| `gopassStores`                     | List of gopass-stores synced to argocd                                          | nil                               |


## Argocd Gopass-stores

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `gpgKey`                           | Config to mount gpg key into repo-server pod  [here](../common/secret.md)                 |                                   |
| `existingGpgKeySecret`             | Config to mount existing gpg key into repo-server pod  [here](../common/secret.md)        |                                   |
| `sshKey`                           | Config to mount ssh key into repo-server pod  [here](../common/secret.md)                 |                                   |
| `existingSshKeySecret`             | Config to mount existing ssh key into repo-server pod  [here](../common/secret.md)        |                                   |
| `directory`                        | SSH-URL to Repository which is used as gopass secret store                      |                                   |
| `storeName`                        | Name of the gopass secret store                                                 |                                   |
