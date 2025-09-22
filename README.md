# caddy-form-parse
Caddy v2 module for parsing form data in request body, forked from [caddy-json-parse](https://github.com/abiosoft/caddy-json-parse)

## Installation

```
xcaddy build v2.0.0 \
    --with github.com/lemniskett/caddy-form-parse
```

## Usage

`form_parse` parses the request body as form data for reference as [placeholders](https://caddyserver.com/docs/caddyfile/concepts#placeholders).

### Caddyfile

Simply use the directive anywhere in a route. You can define formvalue multiple times to parse multiple keys.
```
form_parse <formvalue...>
```

And reference variables via `{form.*}` placeholders.

## License

Apache 2
