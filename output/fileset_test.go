// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package output_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/andreychh/tgen/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// partial represents a view that writes content and then returns err, as a
// template stopped midway by an assert does. A nil err renders content in full.
type partial struct {
	content string
	err     error
}

// Render implements [output.View].
func (p partial) Render(w io.Writer) error {
	_, err := io.WriteString(w, p.content)
	if err != nil {
		return err
	}
	return p.err
}

func TestFileset_Emit(t *testing.T) {
	broken := errors.New("assertion failed")
	seed := func(t *testing.T, path, content string) {
		t.Helper()
		require.NoError(
			t,
			os.MkdirAll(filepath.Dir(path), 0o750),
			"directory of an existing file must be prepared before Emit runs",
		)
		require.NoError(
			t,
			os.WriteFile(path, []byte(content), 0o600),
			"existing file must be prepared before Emit runs",
		)
	}
	t.Run("writes every artifact exactly as it rendered", func(t *testing.T) {
		dir := t.TempDir()
		err := output.NewFileset(output.Artifacts{
			"api.go":   partial{content: "// ① api\n"},
			"types.py": partial{content: "# ② types\n"},
		}).Emit(dir)
		require.NoError(t, err, "Emit must succeed when every artifact renders")
		api, err := os.ReadFile(filepath.Join(dir, "api.go"))
		require.NoError(t, err, "a rendered artifact must be written")
		types, err := os.ReadFile(filepath.Join(dir, "types.py"))
		require.NoError(t, err, "a rendered artifact must be written")
		assert.Equal(
			t,
			[]string{"// ① api\n", "# ② types\n"},
			[]string{string(api), string(types)},
			"Emit must write each artifact exactly as it rendered",
		)
	})
	t.Run("creates the output directory and its missing parents", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "out", "deeper")
		err := output.NewFileset(output.Artifacts{
			"api.go": partial{content: "// api\n"},
		}).Emit(dir)
		require.NoError(t, err, "Emit must succeed when every artifact renders")
		assert.FileExists(
			t,
			filepath.Join(dir, "api.go"),
			"Emit must create the output directory and its missing parents",
		)
	})
	t.Run("returns the error a view failed with", func(t *testing.T) {
		err := output.NewFileset(output.Artifacts{
			"api.go":    partial{content: "// new\n"},
			"client.go": partial{content: "// cut off", err: broken},
		}).Emit(t.TempDir())
		assert.ErrorIs(t, err, broken, "Emit must report the error the failing view returned")
	})
	t.Run(
		"leaves an existing file byte-identical when another artifact fails to render",
		func(t *testing.T) {
			dir := t.TempDir()
			seed(t, filepath.Join(dir, "api.go"), "// old — previous generation\n")
			_ = output.NewFileset(output.Artifacts{
				"api.go":    partial{content: "// new\n"},
				"client.go": partial{content: "// cut off", err: broken},
			}).Emit(dir)
			data, err := os.ReadFile(filepath.Join(dir, "api.go"))
			require.NoError(t, err, "an existing file must still be readable")
			assert.Equal(
				t,
				"// old — previous generation\n",
				string(data),
				"Emit must not touch an existing file when any artifact fails to render",
			)
		},
	)
	t.Run("leaves an existing file whole when its own view fails to render", func(t *testing.T) {
		dir := t.TempDir()
		seed(t, filepath.Join(dir, "client.go"), "package api\n\n// ✓ whole file\n")
		_ = output.NewFileset(output.Artifacts{
			"client.go": partial{content: "// cut off", err: broken},
		}).Emit(dir)
		data, err := os.ReadFile(filepath.Join(dir, "client.go"))
		require.NoError(t, err, "an existing file must still be readable")
		assert.Equal(
			t,
			"package api\n\n// ✓ whole file\n",
			string(data),
			"Emit must not truncate a file whose view fails to render",
		)
	})
	t.Run("writes no artifact when another one fails to render", func(t *testing.T) {
		dir := t.TempDir()
		_ = output.NewFileset(output.Artifacts{
			"api.go":    partial{content: "// new\n"},
			"client.go": partial{content: "// cut off", err: broken},
		}).Emit(dir)
		assert.NoFileExists(
			t,
			filepath.Join(dir, "api.go"),
			"Emit must write no artifact unless every one of them rendered",
		)
	})
	t.Run("creates no output directory when an artifact fails to render", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "out")
		_ = output.NewFileset(output.Artifacts{
			"api.go": partial{content: "// cut off", err: broken},
		}).Emit(dir)
		assert.NoDirExists(
			t,
			dir,
			"Emit must not create directories when an artifact fails to render",
		)
	})
	t.Run("returns an error when the output directory cannot be created", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "taken")
		seed(t, file, "not a directory")
		err := output.NewFileset(output.Artifacts{
			"api.go": partial{content: "// api\n"},
		}).Emit(filepath.Join(file, "out"))
		assert.Error(t, err, "Emit must report a directory it cannot create")
	})
}
