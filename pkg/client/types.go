package client

import (
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	Object            = ctrlclient.Object
	ObjectKey         = ctrlclient.ObjectKey
	ObjectList        = ctrlclient.ObjectList
	Patch             = ctrlclient.Patch
	GetOption         = ctrlclient.GetOption
	ListOption        = ctrlclient.ListOption
	CreateOption      = ctrlclient.CreateOption
	UpdateOption      = ctrlclient.UpdateOption
	DeleteOption      = ctrlclient.DeleteOption
	PatchOption       = ctrlclient.PatchOption
	InNamespace       = ctrlclient.InNamespace
	PropagationPolicy = ctrlclient.PropagationPolicy
	MatchingLabels    = ctrlclient.MatchingLabels
)

var RawPatch = ctrlclient.RawPatch
