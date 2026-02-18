package hl7

import (
	"errors"
	"slices"
	"strings"
)

func AbstractHL7(message string, path HL7Path) (string, error) {
	// just do a check before wasting time parsing the message if the path is invalid
	if err := path.Validate(); err != nil {
		return "", err
	}
	/**
	* This function will take an HL7 message and a path, and return the value at that path in the message.
	* First we need to check that the message is mostly valid and extract the separators
	* - It must begin with MSH
	* - It must have a field separator at the 4th character
	* - It must have separators between the first field separator and the 2nd
	*   field separator in the MSH segment. Like MSH|...| not MSH||
	* - The separators must not be reused. Like MSH|^~\&| not MSH|1111|
	* - There can be a max of 5 separators but the 5th one is not really
	*   supported here.
	* - By default, the separators are ^~\&# but they can be interchanged
	*   dynamically in the MSH segment.
	 */
	// if the path is 0 value, return the whole message
	if path == (HL7Path{}) {
		return message, nil
	}

	// validate message begins with MSH
	if len(message) < 3 || message[:3] != "MSH" {
		return "", errors.New("invalid HL7 message: must begin with MSH")
	}
	// get the next 6 characters after MSH which should be the separators
	// if there are not 6 characters after MSH, it's an error because the separators must be defined
	if len(message) < 10 {
		return "", errors.New("invalid HL7 message: message too short to contain separators and meaningful data")
	}
	separators := message[3:10]
	fieldSeparator := separators[0]
	if path.Segment == "MSH" && path.Field == 1 {
		// MSH-1 is the field separator itself, so return that if requested
		return string(fieldSeparator), nil
	}
	componentSeparator := separators[1]
	if componentSeparator == fieldSeparator {
		return "", errors.New("missing component separator")
	}
	repetitionSeparator := separators[2]
	if repetitionSeparator == fieldSeparator {
		return "", errors.New("missing repetition separator")
	}
	escapeCharacter := separators[3]
	// if escapeCharacter is the same as the fieldSeparator then it is missing
	if escapeCharacter == fieldSeparator {
		return "", errors.New("missing escape character")
	}
	subcomponentSeparator := separators[4]
	if subcomponentSeparator == fieldSeparator {
		return "", errors.New("missing subcomponent separator")
	}
	// there could be a 5th separator we don't care about...
	// but the separators must end with the field separator again.
	if separators[5] != fieldSeparator && separators[6] != fieldSeparator {
		return "", errors.New("unexpected extra separators")
	}

	// check that all separators are unique
	separatorsSet := []byte{fieldSeparator, componentSeparator, repetitionSeparator, escapeCharacter, subcomponentSeparator}
	seen := make(map[byte]bool)
	for _, sep := range separatorsSet {
		if seen[sep] {
			return "", errors.New("separators must be unique")
		}
		seen[sep] = true
	}

	// if we made it here, the message is valid enough to parse the path and
	// extract the value.

	// split the message into segments by the segment separator which could be
	// any of \r, \n, or \r\n.
	segments := splitByAnyOf(message, []string{"\r\n", "\r", "\n"})

	// loop over the segments and find the one that starts with the segment name
	// in the path, if the segment index is greater than 1, we need to find the
	// nth occurrence of the segment.
	segmentCount := 0
	for _, segment := range segments {
		if strings.HasPrefix(segment, path.Segment) {
			segmentCount++
			if segmentCount == path.SegmentIndex {
				// we found the target segment!
				// if field is 0, we want the whole segment returned
				if path.Field == 0 {
					return segment, nil
				}
				// split the segment into fields by the field separator
				fields := strings.Split(segment, string(fieldSeparator))
				// MSH segment has edge case where MSH-1 is the field separator handled above.
				// to handle this edge case here and basically reindex MSH
				// fields, insert "|" at fields[1]
				if path.Segment == "MSH" {
					fields = append(fields[:1], append([]string{string(fieldSeparator)}, fields[1:]...)...)
				}
				// if the field index is greater than the number of fields,
				// return empty string
				if path.Field >= len(fields) {
					return "", nil
				} else {
					field := fields[path.Field]
					// we found the target field!
					// split by repetition unless MSH-2 which is the encoding
					// characters field and does not use repetition, so we will
					// only split by repetition if the path is not MSH-2
					var repetitions []string
					if !(path.Segment == "MSH" && path.Field == 2) {
						repetitions = strings.Split(field, string(repetitionSeparator))
					} else {
						repetitions = []string{field}
					}
					if path.RepetitionIndex > len(repetitions) {
						return "", nil
					} else {
						repetition := repetitions[path.RepetitionIndex-1]
						// we found the target repetition!
						// if component is 0, we want the whole repetition
						// returned
						if path.Component == 0 {
							return repetition, nil
						}
						// split by component...
						components := strings.Split(repetition, string(componentSeparator))
						if path.Component > len(components) {
							return "", nil
						} else {
							component := components[path.Component-1]
							// we found the target component!
							if path.Subcomponent == 0 {
								return component, nil
							}
							// split by subcomponent...
							subcomponents := strings.Split(component, string(subcomponentSeparator))
							if path.Subcomponent > len(subcomponents) {
								return "", nil
							} else {
								return subcomponents[path.Subcomponent-1], nil
							}
						}
					}
				}
			}
		}
	}

	return "", nil
}

func splitByAnyOf(s string, separators []string) []string {
	if len(separators) == 0 {
		return []string{s}
	}
	slices.SortFunc(separators, func(a, b string) int {
		return len(b) - len(a)
	})
	for i, sep := range separators {
		// if the first separator, skip for now because we will use it to split
		// the message into segments later
		if i == 0 {
			continue
		}
		s = strings.ReplaceAll(s, sep, separators[0])
	}

	return strings.Split(s, separators[0])
}
