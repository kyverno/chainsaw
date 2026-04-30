# x_metrics_decode

## Signature

`x_metrics_decode(string)`

## Description

Decodes metrics in the Prometheus text format.

## Examples

```
# decode Prometheus text format metrics and access a specific metric value
x_metrics_decode('http_requests_total{method="GET",status="200"} 1234\n').http_requests_total[0].value == '1234'
```
