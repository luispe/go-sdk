# telemetry

The telemetry package provides plumbing primitives for working with 
tracing and metrics to instrument your webapp.

Welcome to the telemetry user guide. Here we will guide you through some
practical examples of how to interact with the API.

# Getting started

```shell
go get github.com/pomelo-la/go-toolkit/telemetry
```

## Using telemetry

Please Configure the following environment variables for you webapp

| Name                              | Value                          |
|-----------------------------------|--------------------------------|
| OTEL_SERVICE_NAME                 | your-service-name              |
| OTEL_EXPORTER_OTLP_HEADERS        | api-key=<newrelic_license_key> |
| OTEL_EXPORTER_OTLP_ENDPOINT       | https://otlp.nr-data.net:4317  |
| OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT | 4095                           |

## Remarks
- Make sure to use your [ingest license key](https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#license-key)
- If your account is based in the EU, set the endpoint to: https://otlp.eu01.nr-data.net:4317

Please visit and contribute to the community examples

* [Usage examples](https://github.com/pomelo-la/go-toolkit-examples)
