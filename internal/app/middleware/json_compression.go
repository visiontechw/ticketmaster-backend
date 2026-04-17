package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/compress/zstd"
)

func CompressionMiddleware() gin.HandlerFunc {
	encoder, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedDefault))
	return func(c *gin.Context) {
		// Verifica se o cliente aceita compressão
		acceptEnc := c.GetHeader("Accept-Encoding")

		if strings.Contains(acceptEnc, "zstd") {
			c.Header("Content-Encoding", "zstd")
			c.Header("Vary", "Accept-Encoding")

			writer := &compressionWriter{
				ResponseWriter: c.Writer,
				encoder:        encoder,
			}
			c.Writer = writer

			defer writer.Close()
		} else if strings.Contains(acceptEnc, "gzip") {
			c.Header("Content-Encoding", "gzip")
			c.Header("Vary", "Accept-Encoding")
			// Usa middleware do Gin para gzip
		}

		c.Next()
	}
}

type compressionWriter struct {
	gin.ResponseWriter
	encoder *zstd.Encoder
	closed  bool
}

func (cw *compressionWriter) Write(b []byte) (int, error) {
	if cw.closed {
		return 0, errors.New("writer already closed")
	}

	compressed := cw.encoder.EncodeAll(b, nil)
	return cw.ResponseWriter.Write(compressed)
}

func (cw *compressionWriter) Close() {
	if !cw.closed {
		cw.encoder.Close()
		cw.closed = true
	}
}
