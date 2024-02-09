// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofu

import (
	"fmt"
	"testing"

	"github.com/opentofu/opentofu/internal/lang/marks"
	"github.com/zclconf/go-cty/cty"
)

func TestMarksEqual(t *testing.T) {
	for i, tc := range []struct {
		a, b  []cty.PathValueMarks
		equal bool
	}{
		{
			[]cty.PathValueMarks{
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "a"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
			},
			[]cty.PathValueMarks{
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "a"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
			},
			true,
		},
		{
			[]cty.PathValueMarks{
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "a"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
			},
			[]cty.PathValueMarks{
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "A"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
			},
			false,
		},
		{
			[]cty.PathValueMarks{
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "a"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "b"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "c"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
			},
			[]cty.PathValueMarks{
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "b"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "c"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "a"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
			},
			true,
		},
		{
			[]cty.PathValueMarks{
				cty.PathValueMarks{
					Path:  cty.Path{cty.GetAttrStep{Name: "a"}, cty.GetAttrStep{Name: "b"}},
					Marks: cty.NewValueMarks(marks.Sensitive),
				},
				cty.PathValueMarks{
					Path:  cty.Path{cty.GetAttrStep{Name: "a"}, cty.GetAttrStep{Name: "c"}},
					Marks: cty.NewValueMarks(marks.Sensitive),
				},
			},
			[]cty.PathValueMarks{
				cty.PathValueMarks{
					Path:  cty.Path{cty.GetAttrStep{Name: "a"}, cty.GetAttrStep{Name: "c"}},
					Marks: cty.NewValueMarks(marks.Sensitive),
				},
				cty.PathValueMarks{
					Path:  cty.Path{cty.GetAttrStep{Name: "a"}, cty.GetAttrStep{Name: "b"}},
					Marks: cty.NewValueMarks(marks.Sensitive),
				},
			},
			true,
		},
		{
			[]cty.PathValueMarks{
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "a"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
			},
			[]cty.PathValueMarks{
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "b"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
			},
			false,
		},
		{
			nil,
			nil,
			true,
		},
		{
			[]cty.PathValueMarks{
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "a"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
			},
			nil,
			false,
		},
		{
			nil,
			[]cty.PathValueMarks{
				cty.PathValueMarks{Path: cty.Path{cty.GetAttrStep{Name: "a"}}, Marks: cty.NewValueMarks(marks.Sensitive)},
			},
			false,
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if marksEqual(tc.a, tc.b) != tc.equal {
				t.Fatalf("marksEqual(\n%#v,\n%#v,\n) != %t\n", tc.a, tc.b, tc.equal)
			}
		})
	}
}
