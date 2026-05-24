# Brotli for Caddy

This package implements a brotli encoder for [Caddy](https://caddyserver.com/).

Requires Caddy 2+.

Uses the pure Go Brotli implementation <https://github.com/molecule-man/go-brrr>, this seems to have better performance
than the previously used implementation <https://github.com/andybalholm/brotli>.

To quote from [it's README](https://github.com/molecule-man/go-brrr#when-to-use-go-brrr) for when to use brotli:

> For on-the-fly compression, brotli q5-6 is a strong choice if you're already using zstd at its highest level: q5 is
> often faster with a better ratio, and q6 is only slightly slower with an even better ratio. At lower compression
> levels, zstd is significantly faster - if throughput is your priority and you don't need the best ratio, zstd is the
> better tool for the job.

Within caddy that means, that you probably shouldn't use brotli encoding as primary compression algorithm. Zstd and gzip
are the better choice to maximize throughput.

## Installation

```sh
xcaddy build --with github.com/ueffel/caddy-brotli
```

## Syntax

There will be the new encoding `br` available within the
[encode directive](https://caddyserver.com/docs/caddyfile/directives/encode)

```caddyfile
encode [<matcher>] <formats...> {
    br [<level>]
}
```

`level` controls the compression level (ranges from 0 to 11), default is 4.

Example usages could look like this:

```caddyfile
encode br
```

```caddyfile
encode {
    br 4
}
```

or together with gzip

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
