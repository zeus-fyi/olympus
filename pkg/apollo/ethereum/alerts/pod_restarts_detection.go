package apollo_ethereum_alerts

const restartsPromQL = "avg(rate(kube_pod_container_status_restarts_total{}[1h]) * 3600)"
