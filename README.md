![Benthos Captain](icon.png "Benthos Captain")

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

Benthos Captain is a Kubernetes Operator to orchestrate [Benthos](https://www.benthos.dev/) pipelines.

This operator was created with [Operator SDK](https://sdk.operatorframework.io/)

## Build

```
make docker-build docker-push IMG=<some-registry>/benthos-captain:0.1.0
```

## Run

```
make install
make deploy IMG=<some-registry>/benthos-captain:0.1.0
```

Create a sample Pipeline:
```
kubectl apply -f config/samples/benthos_v1beta1_pipeline.yaml -n default
```

See the operator's logs:
```
kubectl logs deployment.apps/benthos-captain-controller-manager -n benthos-captain-operator-system -c manager
```
