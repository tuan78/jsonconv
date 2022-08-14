package benchmarks

import (
	"testing"

	"github.com/tuan78/jsonconv"
)

func BenchmarkFlattenJsonObject(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := sampleJsonObject()
			jsonconv.FlattenJsonObject(data, &jsonconv.FlattenOption{
				Level: jsonconv.FlattenLevelUnlimited,
				Gap:   "__",
			})
		}
	})
}
