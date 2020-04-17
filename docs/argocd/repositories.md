# Repositories

For a repository there are two types, with ssh-connection where an url and a certifacte have to be provided and an https-connection where an URL, username and password have to be provided.

| Parameter                          | Description                                                                               | Default                           |
| ---------------------------------- | ----------------------------------------------------------------------------------------- | --------------------------------- |
| `name`                             | Name of the repository in the Argocd-UI                                                   |                                   |
| `url`                              | Used URL for the repository, (starting "git@" or "https://" )                             |                                   |
| `username`                         | Username for connection to repository [here](../common/secret.md)                         |                                   |
| `existingUsernameSecret`           | Username from existing secret for connection to repository [here](../common/secret.md)    |                                   |
| `password`                         | Password for connection to repository [here](../common/secret.md)                         |                                   |
| `existingPasswordSecret`           | Password from existing secret for connection to repository [here](../common/secret.md)    |                                   |
| `certificate`                      | Certificate for connection to repository [here](../common/secret.md)                      |                                   |
| `existingCertificateSecret`        | Certificate from existing secret for connection to repository [here](../common/secret.md) |                                   |
