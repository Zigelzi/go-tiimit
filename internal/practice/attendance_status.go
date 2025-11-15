package practice

type AttendanceStatus int

const (
	AttendanceIn AttendanceStatus = iota
	AttendanceOut
	AttendanceUnknown
	AttendanceInvalid
)

var attendanceName = map[AttendanceStatus]string{
	AttendanceIn:      "in",
	AttendanceOut:     "out",
	AttendanceUnknown: "unknown",
	AttendanceInvalid: "invalid",
}

func (as AttendanceStatus) String() string {
	return attendanceName[as]
}

func determineStatus(status string) AttendanceStatus {
	var attendanceMap = map[string]AttendanceStatus{
		"Osallistuu":   AttendanceIn,
		"Ei osallistu": AttendanceOut,
		"Ei vastausta": AttendanceUnknown,
	}

	result, exists := attendanceMap[status]
	if !exists {
		return AttendanceInvalid
	}
	return result
}
