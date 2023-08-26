package compressor

import "testing"

var testData = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")

func BenchmarkCompressData(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompressData(testData)
		if err != nil {
			b.Fatalf("Error compressing data: %s", err)
		}
	}
}
