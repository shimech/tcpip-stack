package icmp

type Type uint8

const (
	ICMP_TYPE_ECHOREPLY      Type = 0
	ICMP_TYPE_DEST_UNREACH   Type = 3
	ICMP_TYPE_SOURCE_QUENCH  Type = 4
	ICMP_TYPE_REDIRECT       Type = 5
	ICMP_TYPE_ECHO           Type = 8
	ICMP_TYPE_TIME_EXCEEDED  Type = 11
	ICMP_TYPE_PARAM_PROBLEM  Type = 12
	ICMP_TYPE_TIMESTAMP      Type = 13
	ICMP_TYPE_TIMESTAMPREPLY Type = 14
	ICMP_TYPE_INFO_REQUEST   Type = 15
	ICMP_TYPE_INFO_REPLY     Type = 16
)

func (t *Type) String() string {
	switch *t {
	case ICMP_TYPE_ECHOREPLY:
		return "EchoReply"
	case ICMP_TYPE_DEST_UNREACH:
		return "DestinationUnreachable"
	case ICMP_TYPE_SOURCE_QUENCH:
		return "SourceQuench"
	case ICMP_TYPE_REDIRECT:
		return "Redirect"
	case ICMP_TYPE_ECHO:
		return "Echo"
	case ICMP_TYPE_TIME_EXCEEDED:
		return "TimeExceeded"
	case ICMP_TYPE_PARAM_PROBLEM:
		return "ParameterProblem"
	case ICMP_TYPE_TIMESTAMP:
		return "Timestamp"
	case ICMP_TYPE_TIMESTAMPREPLY:
		return "TimestampReply"
	case ICMP_TYPE_INFO_REQUEST:
		return "InformationRequest"
	case ICMP_TYPE_INFO_REPLY:
		return "InformationReply"
	default:
		return "Unknown"
	}
}
