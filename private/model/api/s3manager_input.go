// +build codegen

package api

import (
	"bytes"
	"fmt"
	"text/template"
)

// S3ManagerUploadInputGoCode returns the Go code for the S3 Upload Manager's
// input structure.
func S3ManagerUploadInputGoCode(a *API) string {
	if v := a.PackageName(); v != "s3" {
		panic(fmt.Sprintf("unexpected API model %s", v))
	}
	s, ok := a.Shapes["PutObjectInput"]
	if !ok {
		panic(fmt.Sprintf("unable to find PutObjectInput shape in S3 model"))
	}
	a.resetImports()
	delete(a.imports, "github.com/aws/aws-sdk-go-v2/aws")
	a.imports["io"] = true
	a.imports["time"] = true
	a.imports["github.com/aws/aws-sdk-go-v2/service/s3"] = true
	var w bytes.Buffer
	if err := s3managerUploadInputTmpl.Execute(&w, s); err != nil {
		panic(fmt.Sprintf("failed to execute %s template, %v",
			s3managerUploadInputTmpl.Name(), err))
	}
	return a.importsGoCode() + w.String()
}

var s3managerUploadInputTmpl = template.Must(
	template.New("s3managerUploadInputTmpl").
		Parse(s3managerUploadInputTmplDef),
)

const s3managerUploadInputTmplDef = `
// UploadInput provides the input parameters for uploading a stream or buffer
// to an object in an Amazon S3 bucket. This type is similar to the s3
// package's PutObjectInput with the exception that the Body member is an
// io.Reader instead of an io.ReadSeeker.
type UploadInput struct {
	_ struct{} {{ .GoTags true false }}
 	{{ range $name, $ref := $.MemberRefs -}}
		{{ if eq $name "Body" }}
			// The readable body payload to send to S3.
			Body io.Reader
		{{ else if eq $name "ContentLength" }}
			{{/* S3 Upload Manager does not use modeled content length */}}
		{{ else }}
			{{ $isBlob := $.WillRefBeBase64Encoded $name -}}
			{{ $isRequired := $.IsRequired $name -}}
			{{ $doc := $ref.Docstring -}}
 			{{ if $doc -}}
				{{ $doc }}
			{{ end -}}
			{{ if $isBlob -}}
				{{ if $doc -}}
					//
				{{ end -}}
				// {{ $name }} is automatically base64 encoded/decoded by the SDK.
			{{ end -}}
			{{ if $isRequired -}}
				{{ if or $doc $isBlob -}}
					//
				{{ end -}}
				// {{ $name }} is a required field
			{{ end -}}
			{{ $name }} {{ $ref.GoTypeWithPkgName }} {{ $ref.GoTags false $isRequired }}
		{{ end }}
	{{ end }}
}
`
