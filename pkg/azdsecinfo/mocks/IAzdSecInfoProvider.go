// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	contracts "github.com/Azure/AzureDefender-K8S-InClusterDefense/pkg/azdsecinfo/contracts"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mock "github.com/stretchr/testify/mock"

	v1 "k8s.io/api/core/v1"
)

// IAzdSecInfoProvider is an autogenerated mock type for the IAzdSecInfoProvider type
type IAzdSecInfoProvider struct {
	mock.Mock
}

// GetContainersVulnerabilityScanInfo provides a mock function with given fields: podSpec, resourceMetadata, resourceKind
func (_m *IAzdSecInfoProvider) GetContainersVulnerabilityScanInfo(podSpec *v1.PodSpec, resourceMetadata *metav1.ObjectMeta, resourceKind *metav1.TypeMeta) ([]*contracts.ContainerVulnerabilityScanInfo, error) {
	ret := _m.Called(podSpec, resourceMetadata, resourceKind)

	var r0 []*contracts.ContainerVulnerabilityScanInfo
	if rf, ok := ret.Get(0).(func(*v1.PodSpec, *metav1.ObjectMeta, *metav1.TypeMeta) []*contracts.ContainerVulnerabilityScanInfo); ok {
		r0 = rf(podSpec, resourceMetadata, resourceKind)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*contracts.ContainerVulnerabilityScanInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1.PodSpec, *metav1.ObjectMeta, *metav1.TypeMeta) error); ok {
		r1 = rf(podSpec, resourceMetadata, resourceKind)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
