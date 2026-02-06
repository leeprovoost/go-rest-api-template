package version

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersionFileValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "VERSION")
	require.NoError(t, os.WriteFile(path, []byte("1.0.0\n"), 0644))

	v, err := ParseVersionFile(path)
	require.NoError(t, err)
	assert.Equal(t, "1.0.0", v)
}

func TestParseVersionFileWithPrefix(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "VERSION")
	require.NoError(t, os.WriteFile(path, []byte("v2.3.4\n"), 0644))

	v, err := ParseVersionFile(path)
	require.NoError(t, err)
	assert.Equal(t, "v2.3.4", v)
}

func TestParseVersionFilePrerelease(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "VERSION")
	require.NoError(t, os.WriteFile(path, []byte("1.0.0-beta.1\n"), 0644))

	v, err := ParseVersionFile(path)
	require.NoError(t, err)
	assert.Equal(t, "1.0.0-beta.1", v)
}

func TestParseVersionFileMissing(t *testing.T) {
	_, err := ParseVersionFile("/nonexistent/path/VERSION")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reading version file")
}

func TestParseVersionFileInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "VERSION")
	require.NoError(t, os.WriteFile(path, []byte("not-a-version\n"), 0644))

	_, err := ParseVersionFile(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a valid version number")
}

func TestParseVersionFileEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "VERSION")
	require.NoError(t, os.WriteFile(path, []byte("\n"), 0644))

	_, err := ParseVersionFile(path)
	assert.Error(t, err)
}
