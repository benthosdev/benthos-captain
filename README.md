![Benthos Captain](icon.png "Benthos Captain")

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

> ⚠️ **This is a work in progress proof of concept** ⚠️  


Benthos Captain is a Kubernetes Operator to orchestrate [Benthos](https://www.benthos.dev/) pipelines.

This operator has been created with [Operator SDK](https://sdk.operatorframework.io/)

## Build

```
make docker-build docker-push IMG=<some-registry>/benthos-captain:0.1.0
```

## Install the Pipelines CRD

```
make install
```

## Check the newly created CRD

```
kubectl get crd pipelines.captain.benthos.dev -oyaml
```

## Deploy Benthos-Captain operator

```
make deploy IMG=<some-registry>/benthos-captain:0.1.0
```

## Create a sample Pipeline:
```
kubectl apply -f config/samples/benthoscaptain_v1alpha1_pipeline.yaml -n default
```

## Check the newly created Pipeline:
```
kubectl get pipelines -n default
```

## See benthos-captain controller manager logs:
```
kubectl logs -f deploy/benthos-captain-controller-manager -n benthos-captain-system -c manager
```

## See sample benthos pipeline logs:
```
kubectl logs -f deploy/pipeline-sample -n default
```

# Helm Chart

If you want to deploy Benthos-Captain with Helm, see [here](chart/benthos-captain/README.md)
