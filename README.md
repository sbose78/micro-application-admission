# MicroApplication Mutating Webhook

When a new [MicroApplication](https://github.com/sbose78/micro-application) is created, this webhook mutates the object to add/update the creator of the resource in an annotation.

```
apiVersion: argoproj.io/v1alpha1
kind: MicroApplication
metadata:
  annotations:
    foo: bar
    generated-creator: 'kube:admin'
```


## Install


1. Create the `MicroApplication` CRD 

```
$ kubectl apply -f https://github.com/sbose78/micro-application/blob/main/config/crd/bases/argoproj.io_microapplications.yaml
```

2. Create the namespace and install the webhook.

On OpenShift,

```
$ kubectl apply -f https://github.com/sbose78/micro-application-admission/blob/main/manifests/openshift/install.yaml
```

On other distributions of Kubernetes,

```
$ kubectl apply -f https://github.com/sbose78/micro-application-admission/blob/main/manifests/kubernetes/install.yaml
```




## Verify

1. The `webhook-server` pod in the `microapplication-webhook` namespace should be running:
```
$ kubectl -n microapplication-webhook get pods
NAME                             READY     STATUS    RESTARTS   AGE
webhook-server-6f976f7bf-hssc9   1/1       Running   0          35m
```

2. A `MutatingWebhookConfiguration` named `microapplication-webhook` should exist:
```
$ kubectl get mutatingwebhookconfigurations
NAME           AGE
microapplication-webhook   36m
```

3. Create a `MicroApplication`

```
apiVersion: argoproj.io/v1alpha1
kind: MicroApplication
metadata:
  annotations:
    foo: bar
    generated-creator: 'kube:admin'
  name: example
  namespace: any-namespace
spec:
  repoURL: https://github.com/sbose78/gitops-samples
  path: developer/new-app
```

Fetch the recently created `MicroApplication` named `example` in namespace `any-namespace`. Notice that the annotation `generated-creator` has been filled in.

```
apiVersion: argoproj.io/v1alpha1
kind: MicroApplication
metadata:
  annotations:
    foo: bar
    generated-creator: 'john'
  name: example
  namespace: any-namespace
spec:
  repoURL: https://github.com/sbose78/gitops-samples
  path: developer/new-app
```

## Build the Image from Sources (optional)

An image can be built by running `make`.
If you want to modify the webhook server for testing purposes, be sure to set and export
the shell environment variable `IMAGE` to an image tag for which you have push access. You can then
build and push the image by running `make push-image`. Also make sure to change the image tag
in `deployment/deployment.yaml.template`, and if necessary, add image pull secrets.


## Credits

This code is based on https://github.com/stackrox/admission-controller-webhook-demo

