// Copyright (c) 2020 Smartling
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package smartling

type TranslationState string

const (
	TranslationStatePublished       TranslationState = "PUBLISHED"
	TranslationStatePostTranslation TranslationState = "POST_TRANSLATION"
)

type ImportRequest struct {
	FileURIRequest

	File             []byte
	FileType         FileType
	TranslationState TranslationState
	Overwrite        bool
}

func (request *ImportRequest) GetForm() (*Form, error) {
	form, err := request.FileURIRequest.GetForm()
	if err != nil {
		return nil, err
	}

	file, err := form.Writer.CreateFormFile("file", request.FileURI)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(request.File)
	if err != nil {
		return nil, err
	}

	err = form.Writer.WriteField("fileType", string(request.FileType))
	if err != nil {
		return nil, err
	}

	err = form.Writer.WriteField(
		"translationState",
		string(request.TranslationState),
	)
	if err != nil {
		return nil, err
	}

	if request.Overwrite {
		err = form.Writer.WriteField("overwrite", "true")
		if err != nil {
			return nil, err
		}
	}

	return form, nil
}
