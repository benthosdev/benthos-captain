# Developer Guide

Benthos Operator is a Kubernetes Operator built with [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder). There are ways to get a local development environment up and running. The easiest is using [Tilt](https://tilt.dev/). You can still work locally using a makefile and manually build the image if you'd prefer.

## Provision a Cluster

Before we do anything, we'll need a local cluster to work on. The easiest way to create one is to use [kind](https://kind.sigs.k8s.io/docs/user/quick-start/).

```bash
kind create cluster
```

**Note: if you already manage Kubernetes clusters you might want to check your context before proceeding.**

## Using Tilt

Once you've installed Tilt, you'll simply be able to run `tilt up` and it will do everything for you, it will also reload the deployment and build the docker image when you make any code changes.

```bash
tilt up
```

## Manually

If you don't want to use Tilt, you can use the makefile targets generated by kubebuilder.

### Running on the cluster

1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/benthos-captain:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/benthos-captain:tag
```

### Uninstall CRDs

To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller

UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out

1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)
