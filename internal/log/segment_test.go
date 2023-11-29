package log

import (
	"io"
	"os"
	"testing"

	api "github.com/just-mitch/clogd/api/v1"
	"github.com/stretchr/testify/require"
)

func TestNearestMultiple(t *testing.T) {
	require.Equal(t, uint64(8), nearestMultiple(9, 4))
}

func TestSegment(t *testing.T) {
	dir, _ := os.MkdirTemp(os.TempDir(), "segment-test")
	defer os.RemoveAll(dir)

	want := &api.Record{Value: []byte("hello world")}

	c := Config{}
	c.Segment.MaxStoreBytes = 1024
	c.Segment.MaxIndexBytes = entWidth * 3

	seg, err := newSegment(dir, 16, c)
	require.NoError(t, err)
	require.Equal(t, uint64(16), seg.nextOffset, seg.nextOffset)
	require.False(t, seg.IsMaxed())

	for i := uint64(0); i < 3; i++ {
		off, err := seg.Append(want)
		require.NoError(t, err)
		require.Equal(t, 16+i, off)

		got, err := seg.Read(off)
		require.NoError(t, err)
		require.Equal(t, want.Value, got.Value)
	}

	_, err = seg.Append(want)
	require.Equal(t, io.EOF, err)
	require.True(t, seg.IsMaxed())

	c.Segment.MaxStoreBytes = uint64(len(want.Value) * 3)
	c.Segment.MaxIndexBytes = 1024

	seg, err = newSegment(dir, 16, c)
	require.NoError(t, err)
	require.True(t, seg.IsMaxed())

	err = seg.Remove()
	require.NoError(t, err)
	s, err := newSegment(dir, 16, c)
	require.NoError(t, err)
	require.False(t, s.IsMaxed())

}
