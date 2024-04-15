package accessor

const (
	AccessGranted AccessStatus = iota
	AccessDenied
	AccessUnknown
)

type AccessStatus uint8 // для передачи состояния доступа

type IAccessor interface {
	CheckAccess(role string, res, act string) (AccessStatus, error)
}
