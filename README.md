# Overview

This is an educational resource demonstrating a relational database running in a Kubernetes cluster that has SQL injection vulnerabilities. Here are a few suggestions for how you might use this repository:

1. Practice your hacking skills by trying to retrieve, insert, or destroy sensitive information.

2. Read the source code found in `main.go` and the solution guide found in `solution/README.md` to better-understand what SQL injection vulnerabilities look like, how to exploit them, and how to prevent them.

3. Use this as a guide/inspiration for building your own applications with vulnerabilities.

The instructions and solutions were written assuming you are using some kind of Linux distribution (sorry Windows :grimacing:), whether Ubuntu, MacOS, or another Debian-based OS.

## \*\*\*\*\**DISCLAIMER*\*\*\*\*\*

This is an application with a built-in security vulnerability. Please don't deploy the Helm chart into a production environment. There are also instructions showing how to exploit SQL injection vulnerabilities, so please don't use this to break any laws :grin:.

# Usage

## Requirements

- Latest version of [Helm](https://helm.sh/docs/intro/install/)
- Latest version of [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- A fully-compliant Kubernetes distribution (i.e. microk8s, k3s, k3d) that is running on Linux/amd64, and is using containerd or Docker as the runtime.

## Deploying to Kubernetes

Add the Helm chart repository:

```sh
helm repo add kube-hack https://kube-hack.github.io/charts
```

Update the charts in your Helm repository:

```sh
helm repo update
```

Deploy the chart to your Kubernetes cluster:

```sh
helm install sql-injection kube-hack/sql-injection
```

## Interacting with the application

### Port-forward the application

```sh
kubectl port-forward svc/client-sql-injection 3000:3000
```

After the application is port-forwarded (accessible via localhost), you can run the below actions:

### Retrieve messages

```sh
# Get all messages for user 1
curl localhost:3000/messages/1
```

### Send messages

```sh
# Send a message to user 1
curl \
--request POST \
--data-binary "Hello, user 1" \
localhost:3000/message/1
```

