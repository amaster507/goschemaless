package internal

import (
	"errors"
	"fmt"
	"regexp"
)

type HL7Path struct {
	Segment         string `json:"segment"`
	SegmentIndex    int    `json:"segment_index"`
	Field           int    `json:"field,omitempty"`
	RepetitionIndex int    `json:"repetition_index,omitempty"`
	Component       int    `json:"component,omitempty"`
	Subcomponent    int    `json:"subcomponent,omitempty"`
}

func (p HL7Path) Validate() error {
	// TODO: do advanced validation based on a specific HL7 version and schema.
	// if Segment is "" then the rest must be empty or 0
	if p.Segment == "" {
		if p.SegmentIndex != 0 || p.Field != 0 || p.RepetitionIndex != 0 || p.Component != 0 || p.Subcomponent != 0 {
			return errors.New("if Segment is empty, the rest of the path must be empty or 0")
		}
		return nil
	}
	// if Segment is MSH then SegmentIndex must be 1
	if p.Segment == "MSH" && p.SegmentIndex != 1 {
		return errors.New("if Segment is MSH, SegmentIndex must be 1")
	}
	// if Segment is MSH and Field is 1, then the rest must be empty or 0
	if p.Segment == "MSH" && p.Field == 1 {
		if p.Component != 0 || p.Subcomponent != 0 {
			return errors.New("if Segment is MSH and Field is 1, the rest of the path must be empty or 0")
		}
	}
	// if Field is set, then Segment must be set
	if p.Field != 0 && p.Segment == "" {
		return errors.New("if Field is set, Segment must be set")
	}
	// if RepetitionIndex is set, then Field must be set
	if p.RepetitionIndex != 0 && p.Field == 0 {
		return errors.New("if RepetitionIndex is set, Field must be set")
	}
	// if Field is set, then RepeitionIndex must be at least 1
	if p.Field != 0 && p.RepetitionIndex == 0 {
		return errors.New("if Field is set, RepetitionIndex must be at least 1")
	}
	// if Component is set, then Field must be set
	if p.Component != 0 && p.Field == 0 {
		return errors.New("if Component is set, Field must be set")
	}
	// if Subcomponent is set, then Component must be set
	if p.Subcomponent != 0 && p.Component == 0 {
		return errors.New("if Subcomponent is set, Component must be set")
	}
	return nil
}

func ParsePath(path string) (HL7Path, error) {
	/*
		 * Need to support the following path formats:
		  - Full Path:
			SEGMENT[
				[SEGMENT_INDEX]
					[-FIELD
						[REPETITION_INDEX]
							[-COMPONENT
								[-SUBCOMPONENT]
							]
					]
			]

		  - Support either - or . as separators
		  - Indexes are optional and default to 1 if not provided
		  - Indexes are 1-based, not 0-based


		 * Example Paths:
		  - PID[1]-5[2].3[1] would be PID,1,5,2,3,1
		  - PV1-2 would be PV1,1,2
		  - MSH-10 would be MSH,1,10
		  - OBX[2].5.2 would be OBX,2,5,1,2
	*/

	// seg & segIndex = ([A-Z0-9]{3})(?:\[(\d+)\])?
	// field & repetitionIndex = (?:[-\.](\d+)(?:\[(\d+)\])?)?
	// component = (?:[-\.](\d+))?
	// subcomponent = (?:[-\.](\d+))?
	/*
		full regexp:
		^([A-Z][A-Z0-9]{2})(?:\[(\d+)\])?(?:[-\.](\d+)(?:\[(\d+)\])?(?:[-\.](\d+)(?:[-\.](\d+))?)?)?$

		regexp explanation:
		^ // start of string
		([A-Z][A-Z0-9]{2}) // segment name: 3 characters, first must be a letter, the rest can be letters or digits
		(?:\[(\d+)\])? // optional segment index in square brackets
		(?:
			[-\.] // separator for field either - or .
			(\d+) // field number
			(?:\[(\d+)\])? // optional repetition index in square brackets
			(?:
				[-\.] // separator for component either - or .
				(\d+) // component number
				(?:
					[-\.] // separator for subcomponent either - or .
					(\d+) // subcomponent number
				)? // optional separator and subcomponent
			)? // optional component and subcomponent with separators
		)? // optional field, repetition index, component, and subcomponent with separators
		$ // end of string
	*/
	res := HL7Path{}

	// allow for a empty path
	if path == "" {
		return res, nil
	}

	captureGroups := []string{
		"segment",
		"segmentIndex",
		"field",
		"repetitionIndex",
		"component",
		"subcomponent",
	}
	pathExp := regexp.MustCompile(`^([A-Z][A-Z0-9]{2})(?:\[(\d+)\])?(?:[-\.](\d+)(?:\[(\d+)\])?(?:[-\.](\d+)(?:[-\.](\d+))?)?)?$`)

	match := pathExp.FindStringSubmatch(path)
	if match == nil {
		return res, errors.New("invalid path format")
	}

	// DEBUGGING: pathExp.SubexpNames only returns the names of captured groups.
	// Need to insetad use a manual mapping of group names to indices based on
	// the regex pattern.

	for i, name := range captureGroups {
		data := match[i+1]
		switch name {
		case "segment":
			segment, err := parseSegmentNameOrError(data)
			if err != nil {
				return res, err
			}
			res.Segment = segment
		case "segmentIndex":
			res.SegmentIndex = parseIntOrDefault(data, 1)
		case "field":
			res.Field = parseIntOrDefault(data, 0)
		case "repetitionIndex":
			def := 0
			if res.Field > 0 {
				def = 1
			}
			res.RepetitionIndex = parseIntOrDefault(data, def)
		case "component":
			res.Component = parseIntOrDefault(data, 0)
		case "subcomponent":
			res.Subcomponent = parseIntOrDefault(data, 0)
		}
	}

	return res, nil
}

func parseSegmentNameOrError(s string) (string, error) {
	if len(s) != 3 {
		return "", errors.New("segment name must be 3 characters")
	}
	for i, char := range s {
		if i == 0 {
			if char < 'A' || char > 'Z' {
				return "", errors.New("segment name must begin with an uppercase letter")
			}
		} else if (char < 'A' || char > 'Z') && (char < '0' || char > '9') {
			return "", errors.New("segment name must be uppercase alphanumeric")
		}
	}
	return s, nil
}

func parseIntOrDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	var res int
	_, err := fmt.Sscanf(s, "%d", &res)
	if err != nil {
		return defaultVal
	}
	return res
}
