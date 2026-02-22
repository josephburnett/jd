package jd

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkRenderStringDiff(b *testing.B) {
	for _, size := range []int{100, 500, 1000, 2000, 5000} {
		b.Run(fmt.Sprintf("len_%d", size), func(b *testing.B) {
			oldStr := strings.Repeat("a", size)
			newStr := strings.Repeat("a", size-1) + "b"
			a, _ := ReadJsonString(fmt.Sprintf(`{"key":"%s"}`, oldStr))
			bNode, _ := ReadJsonString(fmt.Sprintf(`{"key":"%s"}`, newStr))
			d := a.Diff(bNode)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				d.Render()
			}
		})
	}
}

func BenchmarkRenderStringDiffColor(b *testing.B) {
	for _, size := range []int{100, 500, 1000, 2000, 5000} {
		b.Run(fmt.Sprintf("len_%d", size), func(b *testing.B) {
			oldStr := strings.Repeat("a", size)
			newStr := strings.Repeat("a", size-1) + "b"
			a, _ := ReadJsonString(fmt.Sprintf(`{"key":"%s"}`, oldStr))
			bNode, _ := ReadJsonString(fmt.Sprintf(`{"key":"%s"}`, newStr))
			d := a.Diff(bNode)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				d.Render(COLOR)
			}
		})
	}
}

func BenchmarkRenderStringDiffColorWords(b *testing.B) {
	for _, size := range []int{100, 500, 1000} {
		b.Run(fmt.Sprintf("len_%d", size), func(b *testing.B) {
			oldStr := strings.Repeat("a", size)
			newStr := strings.Repeat("a", size-1) + "b"
			a, _ := ReadJsonString(fmt.Sprintf(`{"key":"%s"}`, oldStr))
			bNode, _ := ReadJsonString(fmt.Sprintf(`{"key":"%s"}`, newStr))
			d := a.Diff(bNode)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				d.Render(COLOR_WORDS)
			}
		})
	}
}
