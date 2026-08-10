// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT
package api_test

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

// requireForm reads the request body apart as a multipart form, failing the
// test when it cannot be read as one.
func requireForm(t *testing.T, request *http.Request) Form {
	t.Helper()
	form, err := NewMultipartBody(request).Form()
	require.NoError(t, err)
	return form
}

// MultipartBody is the body of a captured request, taken as the multipart form
// the request declares it to be.
type MultipartBody struct {
	request *http.Request
}

// NewMultipartBody creates a MultipartBody over the request whose body it reads.
func NewMultipartBody(request *http.Request) MultipartBody {
	return MultipartBody{request: request}
}

// Form reads the body apart into the values and the files it carries. It fails
// when the request declares no multipart body, or the body is malformed.
func (b MultipartBody) Form() (Form, error) {
	_, params, err := mime.ParseMediaType(b.request.Header.Get("Content-Type"))
	if err != nil {
		return Form{}, fmt.Errorf("parsing content type: %w", err)
	}
	form := Form{fields: map[string]string{}, files: map[string]FormFile{}}
	reader := multipart.NewReader(b.request.Body, params["boundary"])
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			return form, nil
		}
		if err != nil {
			return Form{}, fmt.Errorf("reading part: %w", err)
		}
		content, err := io.ReadAll(part)
		if err != nil {
			return Form{}, fmt.Errorf("reading part %q: %w", part.FormName(), err)
		}
		if part.FileName() == "" {
			form.fields[part.FormName()] = string(content)
			continue
		}
		form.files[part.FormName()] = FormFile{Name: part.FileName(), Content: string(content)}
	}
}

// Form is one multipart body read apart: the values riding in it under their
// own keys, and the files attached to it under theirs.
type Form struct {
	fields map[string]string
	files  map[string]FormFile
}

// Field returns the value the form carries under key, false when it carries
// none.
func (f Form) Field(key string) (string, bool) {
	value, ok := f.fields[key]
	return value, ok
}

// File returns the file attached to the form under key, false when none is
// attached under it.
func (f Form) File(key string) (FormFile, bool) {
	file, ok := f.files[key]
	return file, ok
}

// Attached returns every key the form carries a file under, in lexical order,
// so that what a request attached can be stated as a whole.
func (f Form) Attached() []string {
	keys := make([]string, 0, len(f.files))
	for key := range f.files {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

// FormFile is one file attached to a multipart form: the name it travels under
// and the bytes it carries.
type FormFile struct {
	Name    string
	Content string
}
