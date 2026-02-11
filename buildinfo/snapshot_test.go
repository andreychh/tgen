// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package buildinfo_test

import (
	"testing"
	"time"

	"github.com/andreychh/tgen/buildinfo"
	"github.com/stretchr/testify/assert"
)

func TestSnapshot_GenerationDate(t *testing.T) {
	snapshot := buildinfo.NewSnapshotAt(
		&buildinfo.Info{Version: "v1.0.0"},
		time.Date(2026, 5, 25, 15, 4, 5, 0, time.UTC),
	)
	assert.Equal(t, "2026-05-25", snapshot.GenerationDate())
}

func TestSnapshot_Version(t *testing.T) {
	snapshot := buildinfo.NewSnapshot(&buildinfo.Info{Version: "v2.5.0-beta.1"})
	assert.Equal(t, "v2.5.0-beta.1", snapshot.Version())
}
