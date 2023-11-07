# webapp

The `webapp` package provides different primitives for working with 
applications within the web environment.

Welcome to the webapp user guide. Here we will guide you through some
practical examples of how to interact with the API.

# Getting started

```shell
go get github.com/pomelo-la/go-toolkit/webapp
```

## Using webapp

## Application

A `webapp.Application` is what most people will end up using, this struct 
is just container of other components exposed by the `go-toolkit` module, 
constructed with sane defaults, compliant with the behavior that's expected 
from an application working in the web environment.

### Requirements

Please Configure the following environment variables for you webapp

| Name                              | Value                          | Mandatory | Default                       |
|-----------------------------------|--------------------------------|-----------|-------------------------------|
| OTEL_SERVICE_NAME                 | your-service-name              | yes       |                               |
| OTEL_EXPORTER_OTLP_HEADERS        | api-key=<newrelic_license_key> | yes       |                               |
| OTEL_EXPORTER_OTLP_ENDPOINT       | https://otlp.nr-data.net:4317  | yes       | https://otlp.nr-data.net:4317 |
| OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT | 4095                           | no        | 4095                          |
| RUNTIME                           |                                | no        | local                         |
| LOG_LEVEL                         |                                | no        | info                          |

## Remarks
- Make sure to use your [ingest license key](https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#license-key)
- If your account is based in the EU, set the endpoint to: https://otlp.eu01.nr-data.net:4317

Please visit and contribute to the community examples

* [Usage examples](https://github.com/pomelo-la/go-toolkit-examples)
