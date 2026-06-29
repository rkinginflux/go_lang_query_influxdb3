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
          env:
            - name: INFLUX_HOST
              value: "http://db3-influxdb3-enterprise-querier.influxdb3.svc.cluster.local:8181"
            - name: INFLUX_TOKEN
              valueFrom:
                secretKeyRef:
                  name: db3-token
                  key: influxToken
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


