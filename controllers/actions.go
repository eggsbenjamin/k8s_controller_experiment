package controllers

import "k8s.io/apimachinery/pkg/runtime"

// Action defines an interface for performing an arbitrary action
type Action interface {
	Execute() error
}

// ActionIdentifier defines an interface for inspecting a k8s resource to determine if
// there are any differences between it's spec and status. The ActionIdentifier should
// then return the action of the highest priority that needs to be performed.
type ActionIdentifier interface {
	IdentifyAction(runtime.Object) (Action, error)
}
