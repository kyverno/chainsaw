package runner

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TestInfo struct {
	Id         int
	ScenarioId int
	Metadata   metav1.ObjectMeta
}

type StepInfo struct {
	Id int
}

type OperationInfo struct {
	Id         int
	ResourceId int
}
