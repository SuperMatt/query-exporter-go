# query-exporter-go

This is a very work in progress tool for getting metrics back from a Prometheus query. This eventually should be used for end to end testing of Prometheus infrastructure to ensure that endpoints can be queried. You can use these metrics to then alert on reasons why metrics cannot be queried.

Possible errors that this exporter should be able to offer

* No data can be queried
* No historical data can be queries
* Query endpoint port is not open
* DNS lookups are failing
* Query endpoint is not correct

## Config

The config for this project can be specified with the `-config` argument, or it will load `config.yaml` or `config.yml` in that order of priority.

A config file looks like this:

```
---
endpoints:
  - name: scrape_name
    address: http://prometheus:9090/
    headers:
      - name: X-Scope-OrgID
        value: fake
    query_offsets:
      - 0      # now
      - 3600   # 1 hour ago
      - 14400  # 4 hours ago

server:
    port: 8080
    metrics_endpoint: /metrics

debug: false
```

The `server` adnd `debug` keys have default values, but endpoints must be specified fully.

