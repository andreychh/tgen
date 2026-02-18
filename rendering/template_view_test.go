// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package rendering_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/andreychh/tgen/rendering"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateView_Render(t *testing.T) {
	tests := []struct {
		name       string
		tmplText   string
		tmplName   string
		data       any
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "successfully renders template with valid data",
			tmplText:   `{{- define "valid" -}}Hello, {{.Name}}!{{- end -}}`,
			tmplName:   "valid",
			data:       map[string]string{"Name": "Telegram"},
			wantOutput: "Hello, Telegram!",
		},
		{
			name:     "returns error when template name does not exist",
			tmplText: `{{- define "valid" -}}Hello{{- end -}}`,
			tmplName: "non_existent",
			data:     map[string]string{},
			wantErr:  true,
		},
		{
			name:     "returns error when template execution fails due to missing key",
			tmplText: `{{- define "bad_data" -}}Hello, {{.MissingKey}}!{{- end -}}`,
			tmplName: "bad_data",
			data:     map[string]string{"WrongKey": "Telegram"},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New("").Option("missingkey=error").Parse(tt.tmplText)
			require.NoError(t, err, "template fixture does not parse correctly")
			view := rendering.NewTemplateView(tmpl, tt.tmplName, tt.data)
			var buf bytes.Buffer
			err = view.Render(&buf)
			if tt.wantErr {
				assert.Error(t, err, "view did not return expected error for invalid input")
				return
			}
			require.NoError(t, err, "unexpected error returned during template execution")
			assert.Equal(
				t,
				tt.wantOutput,
				buf.String(),
				"rendered output does not match expectation",
			)
		})
	}
}
