# Secret

Different tools need the secrets provided differently, 
for every secret there is the possibility to provide the secret inside the kubernetes cluster or integrate the secret encrypted with the orbctl. 

## Integrated

It is recommended to only add integrated secrets with orbctl.

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `encryption`                       | Used encryption for the value                                                   |                                   |
| `encoding`                         | Used encoding for the value                                                     |                                   |
| `value`                            | Encrypted and encoded value                                                     |                                   |

## Existing

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `name`                             | Name of the existent secret which contains the secret                           |                                   |
| `key`                              | Key in the existent secret which contains the secret                            |                                   |
| `internalName`                     | Internal name used to mount (only necessary in some instances)                  |                                   |

## ClientID+ClientSecret or Username+Password (Existing)

Some tools need the provided values always combined into one secret.

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `name`                             | Name of the existent secret which contains the secret                           |                                   |
| `idKey`                            | Key in the existent secret which contains the id                                |                                   |
| `secretKey`                        | Key in the existent secret which contains the secret                            |                                   |
| `internalName`                     | Internal name used to mount (only necessary in some instances)                  |                                   |
