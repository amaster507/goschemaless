package internal

import "testing"

type testCase struct {
	name     string
	path     string
	expected HL7Path
}

var tests = []testCase{
	{
		name: "full path",
		path: "PID[2]-3[4].5.6",
		expected: HL7Path{
			Segment:         "PID",
			SegmentIndex:    2,
			Field:           3,
			RepetitionIndex: 4,
			Component:       5,
			Subcomponent:    6,
		},
	},
	{
		name: "missing indices",
		path: "PV1-3.4",
		expected: HL7Path{
			Segment:         "PV1",
			SegmentIndex:    1,
			Field:           3,
			RepetitionIndex: 1,
			Component:       4,
			Subcomponent:    0,
		},
	},
}

func (test *testCase) Run(t *testing.T) {
	result, err := ParsePath(test.path)
	expectValue(t, result, test.expected, err)
}

func TestParsePath(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.Run(t)
		})
	}
}

func expectValue(t *testing.T, expected any, received any, errors ...error) {
	t.Helper()
	for _, err := range errors {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if received != expected {
		t.Errorf("\nExpected: %+v\nReceived: %+v", expected, received)
	}
}

func expectError(t *testing.T, err error, expectedError ...string) {
	t.Helper()
	if len(expectedError) > 1 {
		t.Fatalf("improper test case definition: more than one expected error provided")
	}
	if err == nil {
		t.Fatalf("expected error: %s but received none", expectedError)
	}
	if len(expectedError) > 0 && err.Error() != expectedError[0] {
		t.Errorf("\nExpected error: %s\nReceived error: %s", expectedError[0], err.Error())
	}
}

func TestParseSegmentNameOrError(t *testing.T) {
	v1 := "PID"
	result, err := parseSegmentNameOrError(v1)
	expectValue(t, result, v1, err)

	v2 := "INVALID"
	_, err = parseSegmentNameOrError(v2)
	expectError(t, err, "segment name must be 3 characters")

	v3 := "123"
	_, err = parseSegmentNameOrError(v3)
	expectError(t, err, "segment name must begin with an uppercase letter")

	v4 := "PV_"
	_, err = parseSegmentNameOrError(v4)
	expectError(t, err, "segment name must be uppercase alphanumeric")
}
