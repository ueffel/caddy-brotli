# Brotli for Caddy

This package implements a brotli encoder for [Caddy](https://caddyserver.com/).

Requires Caddy 2+.

Uses the pure go brotli implementation <https://github.com/andybalholm/brotli>

This implementation is NOT high performance, so it is not recommended to use this encoding as
primary compression algorithm. Use gzip or zstd instead.

The new matchfinder algorithm implementation in andybalholm/brotli is also supported, use with `v2`.

> Currently they [new compression algorithms] give better results than the old implementation (at least for compressing
> my test file, Newton’s Opticks) on levels 2 to 7.

quoted from the README of <https://github.com/andybalholm/brotli>.

## Installation

```sh
xcaddy build --with github.com/ueffel/caddy-brotli
```

## Syntax

There will be the new encoding `br` available within the
[encode directive](https://caddyserver.com/docs/caddyfile/directives/encode)

```caddyfile
encode [<matcher>] <formats...> {
    br [[<level>] [v2]]
}
```

* `level` controls the compression level (ranges from 0 to 11), default is 4.  
  for the new algorithm (`v2`) the levels range from 2 to 7, default is also 4.
* `v2` if set the new brotli algorithm is used

Example usages could look like this:

```caddyfile
encode br
```

```caddyfile
encode {
    br 4
}
```

```caddyfile
encode {
    br v2
}
```

```caddyfile
encode {
    br 4 v2
}
```

or together with other encodings like gzip

```caddyfile
encode gzip br
```

```caddyfile
encode {
    gzip 5
    br 4
}
```

## Remarks

Update 2: From Caddy v2.4.0 onwards preferred order is implied by definition order.

Update: Since Caddy v2.4.0-beta.2 the preferred order of encodings can be set via `prefer` setting.

> There is currently no way to set a prefered order of content-encodings via
> caddy's configuration. The content-encoding is determined by the clients
> preference. In most cases that means a response is encoded with the first
> accepted encoding in the `Accept-Encoding` header of the request that the caddy
> also supports.
>
> Example:
>
> Caddyfile
>
> ```caddyfile
> encode gzip br
> ```
>
> * Request:
>
>   ```plain
>   [...]
>   Accept-Encoding: deflate, gzip, br
>   [...]
>   ```
>
>   Response will be:
>
>   ```plain
>   [...]
>   Content-Encoding: gzip
>   [...]
>   ```
>
> * Request: (different order of encodings)
>
>   ```plain
>   [...]
>   Accept-Encoding: deflate, br, gzip
>   [...]
>   ```
>
>   Response will be:
>
>   ```plain
>   [...]
>   Content-Encoding: br
>   [...]
>   ```
