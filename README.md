# sidecar-injector

A contrived reference example illustrating the use of [operator-sdk](https://sdk.operatorframework.io) in developing a basic custom controller that operates on Deployments.

`sidecar-injector` is intended as the companion repository to [What is the Kubernetes Controller Pattern?](https://mikeball.info/blog/what-is-the-kubernetes-controller-pattern/)

## Getting Started

You'll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.

**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster

1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/sidecar-injector:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/sidecar-injector:tag
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

### How it works

This project uses a [Controller](https://kubernetes.io/docs/concepts/architecture/controller/) which provides a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out

Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

### Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)
