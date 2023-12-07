package auth

import (
	"fmt"

	"github.com/casbin/casbin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Authorizer struct {
	enforcer *casbin.Enforcer
}

func New(enforcer *casbin.Enforcer) *Authorizer {
	return &Authorizer{enforcer: enforcer}
}

func (a *Authorizer) Authorize(subject, object, action string) error {
	if !a.enforcer.Enforce(subject, object, action) {
		msg := fmt.Sprintf("forbidden: %s, %s, %s", subject, object, action)
		st := status.New(codes.PermissionDenied, msg)
		return st.Err()
	}
	return nil
}