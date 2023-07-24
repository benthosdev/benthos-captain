# Benthos Captain

Benthos Captain is a tool that automatically ensures the orchestration of benthos pipelines.

## Introduction

This chart bootstraps a [Benthos Captain](https://github.com/mfamador/benthos-captain) deployment on
a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

Kubernetes >= v1.11


Add the Benthos repo:

```sh
helm repo add benthos https://charts.benthos.dev
```

#### Install the chart with the release name `benthos-captain`

1. Create the benthos namespace:

   ```sh
   kubectl create namespace benthos-captain-system
   ```
   
1. Run helm install:

   ```sh
   helm upgrade -i benthos-captain benthos/benthos-captain \
   --namespace benthos-captain-system
   ```


### Uninstalling the Chart

To uninstall/delete the `benthos-captain` deployment:

```sh
helm delete benthos-captain
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

### Configuration

The following tables lists the configurable parameters of the Benthos Captain chart and their default values.

| Parameter                                         | Default                                              | Description
| -----------------------------------------------   | ---------------------------------------------------- | ---
| `image.repository`                                | `docker.io/marcoamador/benthos-captain`              | Image repository
| `image.tag`                                       | `<VERSION>`                                          | Image tag
| `replicaCount`                                    | `1`                                                  | Number of Benthos-Captain pods to deploy.
| `image.pullPolicy`                                | `IfNotPresent`                                       | Image pull policy
| `image.pullSecret`                                | `None`                                               | Image pull secret
| `logFormat`                                       | `fmt`                                                | Log format (fmt or json)
| `resources.requests.cpu`                          | `50m`                                                | CPU resource requests for the Benthos Captain deployment
| `resources.requests.memory`                       | `64Mi`                                               | Memory resource requests for the Benthos Captain deployment
| `resources.limits`                                | `None`                                               | CPU/memory resource limits for the Benthos Captain deployment
| `nodeSelector`                                    | `{}`                                                 | Node Selector properties for the Benthos Captain deployment
| `tolerations`                                     | `[]`                                                 | Tolerations properties for the Benthos Captain deployment
| `affinity`                                        | `{}`                                                 | Affinity properties for the Benthos Captain deployment
| `extraVolumeMounts`                               | `[]`                                                 | Extra volumes mounts
| `extraVolumes`                                    | `[]`                                                 | Extra volumes
| `dnsPolicy`                                       | ``                                                   | Pod DNS policy
| `dnsConfig`                                       | ``                                                   | Pod DNS config
| `extraEnvs`                                       | `[]`                                                 | Extra environment variables for the Benthos Captain pod(s)
| `env.secretName`                                  | ``                                                   | Name of the secret that contains environment variables which should be defined in the Benthos-Captain container (using `envFrom`)
| `rbac.create`                                     | `true`                                               | If `true`, create and use RBAC resources
| `rbac.pspEnabled`                                 | `false`                                              | If `true`, create and use a restricted pod security policy for Benthos-Captain pod(s)
| `allowedNamespaces`                               | `[]`                                                 | Allow benthos-captain to manage resources in the specified namespaces. The namespace benthos deployed in will always be included
| `serviceAccount.create`                           | `true`                                               | If `true`, create a new service account
| `serviceAccount.name`                             | `benthos-captain`                                    | Service account to be used
| `serviceAccount.annotations`                      | ``                                                   | Additional Service Account annotations
| `clusterRole.create`                              | `true`                                               | If `false`, Benthos Captain will be restricted to the namespaces given in `allowedNamespaces` and the namespace where it is deployed
| `service.type`                                    | `ClusterIP`                                          | Service type to be used (exposing the Benthos Captain API outside of the cluster is not advised)
| `service.port`                                    | `3030`                                               | Service port to be used
| `sync.state`                                      | `git`                                                | Where to keep sync state; either a tag in the upstream repo (`git`), or as an annotation on the SSH secret (`secret`)
| `sync.timeout`                                    | `None`                                               | Duration after which sync operations time out (defaults to `1m`)
| `sync.interval`                                   | `<git.pollInterval>`                                 | Controls how often Benthos Captain will apply what’s in git, to the cluster, absent new commits (defaults to `git.pollInterval`)
| `podLabels`                                       | `{}`                                                 | Additional labels for the Benthos-Captain pod
| `prometheus.enabled`                              | `false`                                              | If enabled, adds prometheus annotations to Benthos Captain and helmOperator pod(s)
| `prometheus.serviceMonitor.create`                | `false`                                              | Set to true if using the Prometheus Operator
| `prometheus.serviceMonitor.interval`              | ``                                                   | Interval at which metrics should be scraped
| `prometheus.serviceMonitor.namespace`             | ``                                                   | The namespace where the ServiceMonitor is deployed
| `prometheus.serviceMonitor.additionalLabels`      | `{}`                                                 | Additional labels to add to the ServiceMonitor
| `hostAliases`                                     | `{}`                                                 | Additional hostAliases to add to the Benthos-Captain pod(s). See <https://kubernetes.io/docs/concepts/services-networking/add-entries-to-pod-etc-hosts-with-host-aliases/>
| `dashboards.enabled`                              | `false`                                              | If enabled, benthos-captain will create a configmap with a dashboard in json that's going to be picked up by grafana (see [sidecar.dashboards.enabled](https://github.com/helm/charts/tree/master/stable/grafana#configuration)). Also remember to set `prometheus.enabled=true` to expose the metrics.
| `dashboards.namespace`                            | ``                                                   | The namespace where the dashboard is deployed, defaults to the installation namespace
| `dashboards.nameprefix`                           | `benthos-captain-dashboards`                         | The prefix of the generated configmaps

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```sh
helm upgrade -i benthos-captain benthos/benthos-captain \
--set replicaCount=2 \
--namespace benthos-captain-system \
benthos/benthos-captain
```

### Upgrade

Update Benthos Captain version with:

```sh
helm upgrade --reuse-values benthos/benthos-captain \
--set image.tag==.!.=
```
