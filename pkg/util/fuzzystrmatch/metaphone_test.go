// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package fuzzystrmatch

import "testing"

func testMetaphone(t *testing.T) {
	tt := []struct {
		Source    string
		MaxLength int
		Expected  string
	}{
		{
			Source:    "GUMBO",
			MaxLength: 4,
			Expected:  "KM",
		},
	}

	for _, tc := range tt {
		got := Metaphone(tc.Source, tc.MaxLength)
		if tc.Expected != got {
			t.Fatalf("error convert string to its Metaphone with source=%q max_length=%d"+
				" expected %s got %s", tc.Source, tc.MaxLength, tc.Expected, got)
		}
	}
}
