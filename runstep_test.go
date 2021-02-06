// Copyright (c) 2021 Pierce Bartine. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package atlantis

import (
	"os"
	"reflect"
	"testing"
)

func TestNewRunStepFails(t *testing.T) {
	tests := []struct {
		name    string
		want    *RunStep
		wantErr bool
	}{
		{
			name:    "EnvironmentVariablesMissing",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRunStep()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRunStep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRunStep() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRunStepSucceeds(t *testing.T) {
	for key, value := range map[string]string{
		"WORKSPACE":                  "default",
		"ATLANTIS_TERRAFORM_VERSION": "0.11.0",
		"DIR":                        "/root/dir1/dir2",
		"PLANFILE":                   "/root/test.tfplan",
		"BASE_REPO_NAME":             "atlantis",
		"BASE_REPO_OWNER":            "runatlantis",
		"HEAD_REPO_NAME":             "atlantis",
		"HEAD_REPO_OWNER":            "acme-corp",
		"HEAD_BRANCH_NAME":           "feature",
		"BASE_BRANCH_NAME":           "main",
		"PROJECT_NAME":               "",
		"PULL_NUM":                   "2",
		"PULL_AUTHOR":                "acme-user",
		"REPO_REL_DIR":               "dir1/dir2",
		"USER_NAME":                  "atlantis",
		"COMMENT_ARGS":               `\a\r\g\1,\a\r\g\2`,
	} {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	tests := []struct {
		name    string
		want    *RunStep
		wantErr bool
	}{
		{
			name: "EnvironmentVariablesExist",
			want: &RunStep{
				TerraformWorkspace: "default",
				TerraformVersion:   "0.11.0",
				AbsoluteProjectDir: "/root/dir1/dir2",
				Planfile:           "/root/test.tfplan",
				BaseRepoName:       "atlantis",
				BaseRepoOwner:      "runatlantis",
				HeadRepoName:       "atlantis",
				HeadRepoOwner:      "acme-corp",
				HeadBranchName:     "feature",
				BaseBranchName:     "main",
				ProjectName:        "",
				PullNum:            2,
				PullAuthor:         "acme-user",
				RelativeProjectDir: "dir1/dir2",
				Username:           "atlantis",
				CommentArgs:        []string{"arg1", "arg2"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRunStep()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRunStep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRunStep() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentArgs_Decode(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		a       CommentArgs
		args    args
		want    CommentArgs
		wantErr bool
	}{
		{
			name:    "ImproperEscapeFails",
			a:       CommentArgs{},
			args:    args{value: `abc`},
			want:    []string{},
			wantErr: true,
		},
		{
			name:    "AmbiguousInputFails",
			a:       CommentArgs{},
			args:    args{value: `\,`},
			want:    []string{},
			wantErr: true,
		},
		{
			name:    "NoArgs",
			a:       CommentArgs{},
			args:    args{value: ``},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "SingleBackslash",
			a:       CommentArgs{},
			args:    args{value: `\\`},
			want:    []string{`\`},
			wantErr: false,
		},
		{
			name:    "DoubleBackslash",
			a:       CommentArgs{},
			args:    args{value: `\\\\`},
			want:    []string{`\\`},
			wantErr: false,
		},
		{
			name:    "TwoArgsSimple",
			a:       CommentArgs{},
			args:    args{value: `\a\r\g\1,\a\r\g\2`},
			want:    []string{"arg1", "arg2"},
			wantErr: false,
		},
		{
			name:    "ThreeArgsWithComma",
			a:       CommentArgs{},
			args:    args{value: `\a\r\g\1,\f\o\o\,\b\a\r\,,\a\r\g\2`},
			want:    []string{"arg1", "foo,bar,", "arg2"},
			wantErr: false,
		},
		{
			name:    "TwoArgsWithJSON",
			a:       CommentArgs{},
			args:    args{value: `\{\"\h\e\l\l\o\"\:\"\f\o\o\\\n\b\a\r\\\n\b\a\z\\\n\"\,\"\w\o\r\l\d\"\:\"\x\\\n\y\\\n\z\"\},\{\"\h\e\l\l\o\"\:\"\f\o\o\\\n\b\a\r\\\n\b\a\z\\\n\"\}`},
			want:    []string{`{"hello":"foo\nbar\nbaz\n","world":"x\ny\nz"}`, `{"hello":"foo\nbar\nbaz\n"}`},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.Decode(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.a, tt.want) {
				t.Errorf("NewRunStep() got = %v, want %v", tt.a, tt.want)
			}
		})
	}
}
