package caddybrotli

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/encode"
	"github.com/itchio/go-brotli/enc"
)

func init() {
	caddy.RegisterModule(Brotli{})
}

// Brotli can create brotli encoders.
type Brotli struct {
	Level int `json:"level,omitempty"`
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
	for d.Next() {
		if !d.NextArg() {
			continue
		}
		levelStr := d.Val()
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			return err
		}
		b.Level = level
	}
	return nil
}

// Provision provisions b's configuration.
func (b *Brotli) Provision(ctx caddy.Context) error {
	if b.Level == 0 {
		b.Level = defaultBrotliLevel
	}
	return nil
}

const (
	Bestspeed       = 0
	Bestcompression = 11
)

// Validate validates b's configuration.
func (b Brotli) Validate() error {
	if b.Level < Bestspeed {
		return fmt.Errorf("quality too low; must be >= %d", Bestspeed)
	}
	if b.Level > Bestcompression {
		return fmt.Errorf("quality too high; must be <= %d", Bestcompression)
	}
	return nil
}

// AcceptEncoding returns the name of the encoding as
// used in the Accept-Encoding request headers.
func (Brotli) AcceptEncoding() string { return "br" }

type brotliEncoder struct {
	*enc.BrotliWriter
	*enc.BrotliWriterOptions
}

func (b *brotliEncoder) Reset(writer io.Writer) {
	newWriter := enc.NewBrotliWriter(writer, b.BrotliWriterOptions)
	b.BrotliWriter = newWriter
}

func (b *brotliEncoder) Close() error {
	return b.BrotliWriter.Close()
}

func (b *brotliEncoder) Write(data []byte) (int, error) {
	return b.BrotliWriter.Write(data)
}

// NewEncoder returns a new brotli writer.
func (b Brotli) NewEncoder() encode.Encoder {
	options := enc.BrotliWriterOptions{Quality: b.Level}
	writer := enc.NewBrotliWriter(ioutil.Discard, &options)
	enc := brotliEncoder{
		BrotliWriter:        writer,
		BrotliWriterOptions: &options,
	}

	return &enc
}

var defaultBrotliLevel = 4

// Interface guards
var (
	_ encode.Encoding       = (*Brotli)(nil)
	_ caddy.Provisioner     = (*Brotli)(nil)
	_ caddy.Validator       = (*Brotli)(nil)
	_ caddyfile.Unmarshaler = (*Brotli)(nil)
	_ encode.Encoder        = (*brotliEncoder)(nil)
)
