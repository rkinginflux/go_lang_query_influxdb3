## command
```
kubectl -n query-client create secret generic db3-token \ 
--from-literal=influxToken="apiv3_xxxxxxxx" \
--dry-run=client -o yaml | kubectl apply -f -
```
