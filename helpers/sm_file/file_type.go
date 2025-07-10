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

package smfile

// FileType represents file type format used in Smartling API.
type FileType string

// Android and next are types that are supported by Smartling API.
const (
	FileTypeUnknown        FileType = ""
	FileTypeAndroid        FileType = "android"
	FileTypeIOS            FileType = "ios"
	FileTypeGettext        FileType = "gettext"
	FileTypeHTML           FileType = "html"
	FileTypeJavaProperties FileType = "javaProperties"
	FileTypeYAML           FileType = "yaml"
	FileTypeXLIFF          FileType = "xliff"
	FileTypeXML            FileType = "xml"
	FileTypeJSON           FileType = "json"
	FileTypeDOCX           FileType = "docx"
	FileTypePPTX           FileType = "pptx"
	FileTypeXLSX           FileType = "xlsx"
	FileTypeIDML           FileType = "idml"
	FileTypeQt             FileType = "qt"
	FileTypeResx           FileType = "resx"
	FileTypePlaintext      FileType = "plaintext"
	FileTypeCSV            FileType = "csv"
	FileTypeStringsdict    FileType = "stringsdict"
)
