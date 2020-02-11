# Examples to use SSO

This yaml-parts are only examples and there are alot of additional configurations possible, but they should desplay the most used cases.

# Grafana

All configuration for SSO is under the "auth"-attribute, whereas the domain has to be set correctly so that the redirect works correctly:

```yaml  
apiVersion: boom.caos.ch/v1beta1
kind: Toolset
metadata:
  name: caos
  namespace: caos-system
spec:
  grafana:
    deploy: true
    domain: example.caos.ch
    auth:
```

## Google

The use google as IDP there is the possbility to limit the allowed hosted-domains:

```yaml
      google:
        secretName: google-auth
        clientIDKey: client_id
        clientSecretKey: client_secret
        allowedDomains:
        - caos.ch
```

## Gitlab

The use google as IDP there is the possbility to limit the allowed groups:

```yaml
      gitlab:
        secretName: gitlab-auth
        clientIDKey: client_id
        clientSecretKey: client_secret
        allowedGroups:
        - caos
```

## Github

The use google as IDP there is the possbility to limit the allowed organizations:

```yaml
      github:
        secretName: github-auth-monitoring
        clientIDKey: client_id
        clientSecretKey: client_secret
        allowedOrganizations:
        - caos
```

# Argocd

All configuration for SSO is under the "auth"-attribute, whereas the rootUrl is used for the redirect urls:

```yaml  
apiVersion: boom.caos.ch/v1beta1
kind: Toolset
metadata:
  name: caos
  namespace: caos-system
spec:
  Argocd:
    deploy: true
    auth:
      rootUrl: https://argocd.caos.ch
```

## Google

The use google as IDP there is the possbility to limit the allowed hosted-domains:

```yaml
      google:
        id: google
        name: google
        config:
          secretName: google-auth
          clientIDKey: client_id
          clientSecretKey: client_secret
          hostedDomains:
          - caos.ch
```

## Gitlab

The use google as IDP there is the possbility to limit the allowed groups:

```yaml
      gitlab:
        id: gitlab
        name: gitlab
        config:
          secretName: gitlab-auth
          clientIDKey: client_id
          clientSecretKey: client_secret
          groups:
          - caos
```

## Github

The use google as IDP there is the possbility to limit the allowed organizations:

```yaml
      github:
        id: github
        name: github
        config:
          secretName: github-auth-reconciling
          clientIDKey: client_id
          clientSecretKey: client_secret
          orgs:
          - name: caos
```
