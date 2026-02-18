package hl7

import "testing"

// MSH|^~\\&|HIS|RIH|EKG|EKG|20060529090131||ADT^A01|MSG00001|P|2.5
// PID|||555-44-4444^^^^SSN~123^^^^MRN||EVERYWOMAN^EVE^E^^^^L~QUE^SUZY^^^^^N||19610615|F||C|2222 HOMES TREET^^GREENSBORO^NC^27401||(919)379-1212|(919)271-3434||S||555-55-5555
// PV1||I|2000^2012^01||||004777^LEBAUER^JAMES^A^^^^MD|||||||||||V
// OBX|1|ST|^Body Height||1.80|m|1.50-2.00|N|||F
// OBX|2|ST|^Body Weight||79|kg|50-100|N|||F
// ZZZ||This is~a^custom&segment&with^custom&fields
// ZZZ||foo|bar|baz

var message = "MSH|^~\\&|HIS|RIH|EKG|EKG|20060529090131||ADT^A01|MSG00001|P|2.5\rPID|||555-44-4444^^^^SSN~123^^^^MRN||EVERYWOMAN^EVE^E^^^^L~QUE^SUZY^^^^^N||19610615|F||C|2222 HOMES TREET^^GREENSBORO^NC^27401||(919)379-1212|(919)271-3434||S||555-55-5555\rPV1||I|2000^2012^01||||004777^LEBAUER^JAMES^A^^^^MD|||||||||||V\rOBX|1|ST|^Body Height||1.80|m|1.50-2.00|N|||F\rOBX|2|ST|^Body Weight||79|kg|50-100|N|||F\rZZZ||This is~a^custom&segment&with^custom&fields\rZZZ||foo|bar|baz"

func TestAbstractHL7(t *testing.T) {
	var err1, err2 error
	var resp string
	var path HL7Path

	path, err1 = ParsePath("")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, message, resp, err1, err2)

	path, err1 = ParsePath("MSH.1")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "|", resp, err1, err2)

	path, err1 = ParsePath("MSH.2")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "^~\\&", resp, err1, err2)

	path, err1 = ParsePath("MSH.3")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "HIS", resp, err1, err2)

	path, err1 = ParsePath("MSH[1].1")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "|", resp, err1, err2)

	path, err1 = ParsePath("MSH[2].1")
	resp, err2 = AbstractHL7(message, path)
	expectError(t, err2, "if Segment is MSH, SegmentIndex must be 1")

	path, err1 = ParsePath("PID.3")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "555-44-4444^^^^SSN", resp, err1, err2)

	path, err1 = ParsePath("PID.3[2]")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "123^^^^MRN", resp, err1, err2)

	path, err1 = ParsePath("PID-3.1")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "555-44-4444", resp, err1, err2)

	path, err1 = ParsePath("PID-3.5")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "SSN", resp, err1, err2)

	path, err1 = ParsePath("PID-3[2].1")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "123", resp, err1, err2)

	path, err1 = ParsePath("PID-3[2].5")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "MRN", resp, err1, err2)

	path, err1 = ParsePath("OBX[1].1")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "1", resp, err1, err2)

	path, err1 = ParsePath("OBX[2].1")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "2", resp, err1, err2)

	// non-existent repetition should return empty string, not error
	path, err1 = ParsePath("OBX[3].1")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "", resp, err1, err2)

	// ZZZ||This is~a^custom&segment&with^custom&fields
	// ZZZ||foo|bar|baz

	path, err1 = ParsePath("ZZZ-2.1")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "This is", resp, err1, err2)

	path, err1 = ParsePath("ZZZ-2[2].2.2")
	resp, err2 = AbstractHL7(message, path)
	expectValue(t, "segment", resp, err1, err2)

}
