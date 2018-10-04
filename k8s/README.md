To test the Kubernetes interface,

```
export KUBECONFIG=<location of your KUBECONFIG>
go test -v
```

All resources are deployed in the "sched-ops-test" namespace. If test fails midway, you can do `kubectl delete ns sched-ops-test` to terminate all resources.
