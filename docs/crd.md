# CRD boom.caos.ch/v1beta1

## Structure

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `currentStatePath`                 | Relative folder path where the currentstate is written to                       |                                   |
| `forceApply`                       | Flag if --force should be used by apply of resources                            |                                   |
| `preApply`                         | Spec for the yaml-files applied before applications                             |                                   |
| `postApply`                        | Spec for the yaml-files applied after applications                              |                                   |
| `prometheus-operator`              | Spec for the Prometheus-Operator                                                |                                   |
| `logging-operator`                 | Spec for the Banzaicloud Logging-Operator                                       |                                   |
| `prometheus-node-exporter`         | Spec for the Prometheus-Node-Exporter                                           |                                   |
| `grafana`                          | Spec for the Grafana                                                            |                                   |
| `ambassador`                       | Spec for the Ambassador                                                         |                                   |
| `kube-state-metrics`               | Spec for the Kube-State-Metrics                                                 |                                   |
| `argocd`                           | Spec for the Argo-CD                                                            |                                   |
| `prometheus`                       | Spec for the Prometheus instance                                                |                                   |
| `loki`                             | Spec for the Loki instance                                                      |                                   |

### Pre-Apply

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |
| `folder`                           | Relative path of folder in cloned git repository which should be applied        |                                   |

### Post-Apply

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |
| `folder`                           | Relative path of folder in cloned git repository which should be applied        |                                   |

### Prometheus-Operator

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |

### Logging-Operator

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |
| `fluentdStorage`                   | Spec to define how the persistency should be handled [here](common/storage.md)  | nil                               |

### Prometheus-Node-Exporter

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |

### Grafana

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |
| `admin`                            | Spec for the definition of the admin account [here](grafana/admin.md)            | nil                               |
| `datasources`                      | Spec for additional datasources                                                 | nil                               |
| `datasources.name`                 | Name of the datasource                                                          |                                   |
| `datasources.type`                 | Type of the datasource (for example prometheus)                                 |                                   |
| `datasources.url`                  | URL to the datasource                                                           |                                   |
| `datasources.access`               | Access defintion of the datasource                                              |                                   |
| `datasources.isDefault`            | Boolean if datasource should be used as default                                 |                                   |
| `dashboardproviders`               | Spec for additional Dashboardproviders                                          | nil                               |
| `dashboardproviders.configMaps`    | ConfigMaps in which the dashboards are stored                                   |                                   |
| `dashboardproviders.folder`        | Local folder in which the dashboards are mounted                                |                                   |
| `storage`                          | Spec to define how the persistency should be handled                            | nil                               |
| `storage.size`                     | Defined size of the PVC                                                         |                                   |
| `storage.storageClass`             | Storageclass used by the PVC                                                    |                                   |
| `storage.accessModes`              | Accessmodes used by the PVC                                                     |                                   |
| `network`                          | Network configuration, [here](common/network.md)                                |                                   |
| `auth`                             | Authorization and Authentication configuration for SSO, [here](sso-examples.md) |                                   |

### Ambassador

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |
| `replicaCount`                     | Number of replicas used for deployment                                          | 1                                 |
| `service`                          | Service definition for ambassador                                               | nil                               |
| `service.type`                     | Type for the service                                                            | NodePort                          |
| `service.loadBalancerIP`           | Used IP for loadbalancing for ambassador if loadbalancer is used                | nil                               |

### Kube-State-Metrics

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |
| `replicaCount`                     | Number of replicas used for deployment                                          | 1                                 |

### Argo-CD

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |
| `customImage`                      | Custom argocd-image                                                             | nil                               |
| `rbacConfig`                       | Config for RBAC in argocd [here](argocd/rbacconfig.md)                          | nil                               |
| `network`                          | Network configuration, [here](common/network.md)                                       | nil                               |
| `auth`                             | Authorization and Authentication configuration for SSO, [here](sso-examples.md) | nil                               |
| `repositories`                     | Repositories used by argocd, [here](argocd/repositories.md)              | nil                               |
| `knownHosts`                       | List of known_hosts as strings for argocd                                       | nil                               |

### Prometheus

When the metrics spec is nil all metrics will get scraped.

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |
| `metrics`                          | Spec to define which metrics should get scraped [here](prometheus/metrics.md)   | nil                               |
| `storage`                          | Spec to define how the persistency should be handled [here](common/storage.md)  | nil                               |

### Loki

When the logs spec is nil all logs will get persisted in loki.

| Parameter                          | Description                                                                     | Default                           |
| ---------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| `deploy`                           | Flag if tool should be deployed                                                 | false                             |
| `logs`                             | Spec to define which logs will get persisted [here](loki/logs.md)               | nil                               |
| `storage`                          | Spec to define how the persistency should be handled [here](common/storage.md)  | nil                               |
