# Docker build, and Containerd commands
```bash
docker build --no-cache -t influxdb3-query-client:latest .
```

```bash
docker save influxdb3-query-client:latest -o influxdb3-query-client.tar
```

```bash
sudo ctr -n k8s.io images import influxdb3-query-client.tar
```

## Check it
```bash
sudo crictl image ls | grep influxdb3-query-client
```

Should see something like this ...
```
IMAGE                                           TAG                 IMAGE ID            SIZE
docker.io/library/influxdb3-query-client        latest              3e2db20ba76a9       35.2MB
```
