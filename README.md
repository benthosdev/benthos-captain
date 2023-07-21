![Benthos Captain](icon.png "Benthos Captain")

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

> ⚠️ **This is a work in progress proof of concept** ⚠️  


Benthos Captain is a Kubernetes Operator to orchestrate [Benthos](https://www.benthos.dev/) pipelines.


## Getting Started

Currently, there isn't a stable release of the operator. If you want to install the operator for development purposes, you can follow the [developer guide](./docs/developer-guide.md).

The operator provides a custom resource for managing Benthos pipelines. Once you've got the operator running, you can deploy a `BenthosPipeline` resource to test it out:
```yaml
---
apiVersion: streaming.benthos.dev/v1alpha1
kind: BenthosPipeline
metadata:
  name: benthospipeline-sample
spec:
  replicas: 2
  config: |
    input:
      generate:
        mapping: root = "woof"
        interval: 5s
        count: 0

    pipeline:
      processors:
        - mapping: root = content().uppercase()

    output:
      stdout: {}
```

Once the resource is deployed, you can monitor the state of the resoure:
```bash
kubectl get benthospipelines

NAME                     READY   PHASE     REPLICAS   AVAILABLE   AGE
benthospipeline-sample   true    Running   2          2           62s
```
