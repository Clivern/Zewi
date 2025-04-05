### Using kubectl

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.14.1/deploy/static/provider/cloud/deploy.yaml
```

### Verification

1. Check that the `NGINX` Ingress Controller pods are running:
```bash
kubectl get pods -n ingress-nginx
```

You should see at least one pod in `Running` state:
```
NAME                                        READY   STATUS    RESTARTS   AGE
ingress-nginx-controller-xxxxxxxxxx-xxxxx   1/1     Running   0          2m
```

2. Check the ingress controller service:
```bash
kubectl get svc -n ingress-nginx
```

For `LoadBalancer` type, wait for the external IP to be assigned:
```
NAME                                 TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   10.96.xxx.xxx   <pending>       80:3xxxx/TCP,443:3xxxx/TCP   2m
```

3. Verify the ingress class exists:
```bash
kubectl get ingressclass
```

You should see:
```
NAME    CONTROLLER             PARAMETERS   AGE
nginx   k8s.io/ingress-nginx   <none>       2m
```
