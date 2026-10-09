package anomaly

type Level int

const (
	ReadUncommitted Level = iota
	ReadCommitted
	SnapshotIsolation
	Serializable
)

func (l Level) String() string {
	switch l {
	case ReadCommitted:
		return "read committed"
	case SnapshotIsolation:
		return "snapshot isolation"
	case Serializable:
		return "serializable"
	default:
		return "read uncommitted"
	}
}

func Requires(a Anomaly) Level {
	switch a {
	case G0:
		return ReadUncommitted
	case G1c, IncompatibleOrder:
		return ReadCommitted
	case GSingle, Internal, LostUpdate:
		return SnapshotIsolation
	case G2:
		return Serializable
	default:
		return ReadUncommitted
	}
}
