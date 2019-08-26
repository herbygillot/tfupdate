package tfupdate

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
)

func TestNewProviderUpdater(t *testing.T) {
	cases := []struct {
		name    string
		version string
		want    Updater
		ok      bool
	}{
		{
			name:    "aws",
			version: "2.23.0",
			want: &ProviderUpdater{
				name:    "aws",
				version: "2.23.0",
			},
			ok: true,
		},
		{
			name:    "",
			version: "2.23.0",
			want:    nil,
			ok:      false,
		},
		{
			name:    "aws",
			version: "",
			want:    nil,
			ok:      false,
		},
	}

	for _, tc := range cases {
		got, err := NewProviderUpdater(tc.name, tc.version)
		if tc.ok && err != nil {
			t.Errorf("NewProviderUpdater() with name = %s, version = %s returns unexpected err: %+v", tc.name, tc.version, err)
		}

		if !tc.ok && err == nil {
			t.Errorf("NewProviderUpdater() with name = %s, version = %s expects to return an error, but no error: %+v", tc.name, tc.version, err)
		}

		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("NewProviderUpdater() with name = %s, version = %s returns %#v, but want = %#v", tc.name, tc.version, got, tc.want)
		}
	}
}

func TestUpdateProvider(t *testing.T) {
	cases := []struct {
		src     string
		name    string
		version string
		want    string
		ok      bool
	}{
		{
			src: `
terraform {
  required_version = "0.12.4"
  required_providers {
    null = "2.1.1"
  }
}
`,
			name:    "null",
			version: "2.1.2",
			want: `
terraform {
  required_version = "0.12.4"
  required_providers {
    null = "2.1.2"
  }
}
`,
			ok: true,
		},
		{
			src: `
provider "aws" {
  version = "2.11.0"
  region  = "ap-northeast-1"
}
`,
			name:    "aws",
			version: "2.23.0",
			want: `
provider "aws" {
  version = "2.23.0"
  region  = "ap-northeast-1"
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		u := &ProviderUpdater{
			name:    tc.name,
			version: tc.version,
		}
		f, diags := hclwrite.ParseConfig([]byte(tc.src), "", hcl.Pos{Line: 1, Column: 1})
		if diags.HasErrors() {
			t.Fatalf("unexpected diagnostics: %s", diags)
		}

		err := u.Update(f)
		if tc.ok && err != nil {
			t.Errorf("Update() with src = %s, name = %s, version = %s returns unexpected err: %+v", tc.src, tc.name, tc.version, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("Update() with src = %s, name = %s, version = %s expects to return an error, but no error", tc.src, tc.name, tc.version)
		}

		got := string(hclwrite.Format(f.BuildTokens(nil).Bytes()))
		if got != tc.want {
			t.Errorf("Update() with src = %s, name = %s, version = %s returns %s, but want = %s", tc.src, tc.name, tc.version, got, tc.want)
		}
	}
}
