package plugin

const (
	URL int = iota
	IP
	IPPORT
)

type Target struct {
	Type int
	Raw  string
}
