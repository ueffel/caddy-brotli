package caddybrotli

import (
	"fmt"
	"strconv"

	"github.com/andybalholm/brotli"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/encode"
)

func init() {
	caddy.RegisterModule(Brotli{})
}

// Brotli can create brotli encoders.
type Brotli struct {
	Level int  `json:"level,omitempty"`
	UseV2 bool `json:"use_v2,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (Brotli) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.encoders.br",
		New: func() caddy.Module { return new(Brotli) },
	}
}

// UnmarshalCaddyfile sets up the encoder from Caddyfile tokens.
func (b *Brotli) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	b.Level = -1
	for d.Next() {
		switch d.CountRemainingArgs() {
		case 0:
			continue
		case 1:
			d.NextArg()
			if d.Val() == "v2" {
				b.UseV2 = true
				continue
			}

			level, err := strconv.Atoi(d.Val())
			if err != nil {
				return err
			}
			b.Level = level
		case 2:
			d.NextArg()
			level, err := strconv.Atoi(d.Val())
			if err != nil {
				return err
			}
			b.Level = level

			d.NextArg()
			if d.Val() != "v2" {
				return d.Errf("invalid argument: %s", d.Val())
			}
			b.UseV2 = true
		default:
			return d.Errf("too many arguments (%d)", d.CountRemainingArgs())
		}
	}
	return nil
}

// Provision provisions b's configuration.
func (b *Brotli) Provision(ctx caddy.Context) error {
	if b.Level == -1 {
		b.Level = defaultBrotliLevel
	}
	return nil
}

// Validate validates b's configuration.
func (b Brotli) Validate() error {
	if b.UseV2 {
		if b.Level < brotliV2MinLevel {
			return fmt.Errorf("quality too low; must be >= %d for the new algorithm", brotliV2MinLevel)
		}
		if b.Level > brotliV2MaxLevel {
			return fmt.Errorf("quality too high; must be <= %d for the new algorithm", brotliV2MaxLevel)
		}
	} else {
		if b.Level < brotli.BestSpeed {
			return fmt.Errorf("quality too low; must be >= %d", brotli.BestSpeed)
		}
		if b.Level > brotli.BestCompression {
			return fmt.Errorf("quality too high; must be <= %d", brotli.BestCompression)
		}
	}
	return nil
}

// AcceptEncoding returns the name of the encoding as
// used in the Accept-Encoding request headers.
func (Brotli) AcceptEncoding() string { return "br" }

// NewEncoder returns a new brotli writer.
func (b *Brotli) NewEncoder() encode.Encoder {
	if b.UseV2 {
		return brotli.NewWriterV2(nil, b.Level)
	}
	return brotli.NewWriterLevel(nil, b.Level)
}

const (
	defaultBrotliLevel = 4
	brotliV2MinLevel   = 2
	brotliV2MaxLevel   = 7
)

// Interface guards.
var (
	_ encode.Encoding       = (*Brotli)(nil)
	_ caddy.Provisioner     = (*Brotli)(nil)
	_ caddy.Validator       = (*Brotli)(nil)
	_ caddyfile.Unmarshaler = (*Brotli)(nil)
)
