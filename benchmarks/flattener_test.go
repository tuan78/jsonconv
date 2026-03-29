package benchmarks

import (
	"testing"

	"github.com/tuan78/jsonconv/v2"
)

func BenchmarkFlatten(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := sampleObject()
			jsonconv.Flatten(data, &jsonconv.FlattenOption{
				Level: jsonconv.FlattenLevelUnlimited,
				Gap:   "__",
			})
		}
	})
}
