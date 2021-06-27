package tsdb

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/prometheus/prometheus/pkg/labels"
)

func BenchmarkHeadMemoryUsage(b *testing.B) {
	chunkDir, err := ioutil.TempDir("", "chunk_dir")
	require.NoError(b, err)
	defer func() {
		require.NoError(b, os.RemoveAll(chunkDir))
	}()
	opts := DefaultHeadOptions()
	opts.ChunkRange = 1000
	opts.ChunkDirRoot = chunkDir
	h, err := NewHead(nil, nil, nil, opts, nil)
	require.NoError(b, err)
	defer h.Close()

	stats := runtime.MemStats{}
	runtime.GC()
	runtime.ReadMemStats(&stats)
	fmt.Printf("Before: %v\n", stats.HeapInuse)

	for i := 0; i < b.N; i += 2 {
		si := strings.Repeat(strconv.Itoa(i), 4)
		var ss []string
		for _, c := range []string{"a", "b", "c", "d", "e", "f", "g", "i"} {
			ss = append(ss, c+si)
		}
		h.getOrCreate(uint64(i), labels.FromStrings(ss...))
		rand.Shuffle(len(ss), func(i, j int) { ss[i], ss[j] = ss[j], ss[i] })
		h.getOrCreate(uint64(i+1), labels.FromStrings(ss...))
	}

	runtime.GC()
	runtime.ReadMemStats(&stats)
	fmt.Printf("After: %v\n", stats.HeapInuse)
}
