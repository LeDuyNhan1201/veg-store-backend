package context

type SecurityContext struct {
	Identity string
	Roles    []string
}
