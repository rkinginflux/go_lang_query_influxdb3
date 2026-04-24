### Manifest yaml file for install
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: influxdb3-query-client
  namespace: query-client
spec:
  replicas: 1
  selector:
    matchLabels:
      app: influxdb3-query-client
  template:
    metadata:
      labels:
        app: influxdb3-query-client
    spec:
      containers:
        - name: influxdb3-query-client
          image: influxdb3-query-client:latest
          imagePullPolicy: Never
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: influxdb3-query-client
  namespace: query-client
spec:
  type: NodePort
  selector:
    app: influxdb3-query-client
  ports:
    - name: http
      port: 80
      targetPort: 8080
      nodePort: 30080
```

## Install
```bash
kubectl apply -f manifest.yaml -n query-client
```

## Patch for external access
```bash
kubectl patch svc -n query-client influxdb3-query-client -p '{"spec":{"externalIPs":["192.168.x.xx"]}}'
```


