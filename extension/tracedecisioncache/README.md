# Cache extension

The cache extension can be used by compatible processors to cache sampling decisions. This is useful when data may be received after a sampling decision has been made. 

## Example configuration

```yaml
extensions:
  zipkin_encoding:
    protocol: zipkin_proto
    version: v2

receivers:
  kafka:
    encoding: zipkin_encoding
    # ... other configuration values
```
