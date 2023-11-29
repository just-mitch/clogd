package log

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "index_test")

	require.NoError(t, err)

	defer os.Remove(f.Name())

	c := Config{}
	c.Segment.MaxIndexBytes = 1024

	idx, err := newIndex(f, c)
	require.NoError(t, err)

	_, _, err = idx.Read(-1)
	require.Equal(t, io.EOF, err)

	require.Equal(t, f.Name(), idx.Name())

	entries := []struct {
		offset   uint32
		position uint64
	}{
		{0, 0},
		{1, 10},
	}

	for _, want := range entries {
		err = idx.Write(want.offset, want.position)
		require.NoError(t, err)

		offset, position, err := idx.Read(int64(want.offset))
		require.NoError(t, err)
		require.Equal(t, want.offset, offset)
		require.Equal(t, want.position, position)
	}

	_, _, err = idx.Read(int64(len(entries)))
	require.Equal(t, io.EOF, err)
	_ = idx.Close()

	f, _ = os.OpenFile(f.Name(), os.O_RDWR, 0600)
	idx, err = newIndex(f, c)
	require.NoError(t, err)

	off, pos, err := idx.Read(-1)
	require.NoError(t, err)

	require.Equal(t, entries[1].offset, off)
	require.Equal(t, entries[1].position, pos)

}
