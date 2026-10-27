package tokens_test

import (
	"testing"

	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

func BenchmarkLexExpression(b *testing.B) {
	cfg := config.New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stream := tokens.LexAll("{{ value | default('fallback') }}", cfg)
		for !stream.End() {
			stream.Next()
		}
	}
}
