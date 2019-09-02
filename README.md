# Demo application for Kubernetes

![screenshot](docs/images/screenshot.png)

```
kubectl create -f .\kubernetes-manifests\1-rbac.yml
kubectl create -f .\kubernetes-manifests\2-deploy.yaml
kubectl logs -l app=kubedemo
```
