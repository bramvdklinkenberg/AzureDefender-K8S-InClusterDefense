package annotations

import (
	"encoding/json"
	"github.com/Azure/AzureDefender-K8S-InClusterDefense/pkg/azdsecinfo/contracts"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

const (
	_expectedTestAddPatchOperation   = "add"
	_expectedTestAnnotationPatchPath = "/metadata/annotations"
	_annotationTestKeyOne = "cluster-autoscaler.kubernetes.io/safe-to-evict"
	_annotationTestValueOne = "true"
	_annotationTestKeyTwo = "container.seccomp.security.alpha.kubernetes.io/manager"
	_annotationTestValueTwo = "runtime/default"
)

type TestSuite struct {
	suite.Suite
	containersScanInfo []*contracts.ContainerVulnerabilityScanInfo
	podNoAnnotations *corev1.Pod
	podWithAnnotations *corev1.Pod
}

func (suite *TestSuite) SetupSuite() {
	suite.containersScanInfo = []*contracts.ContainerVulnerabilityScanInfo{
		{
			Name: "container1",
			Image: &contracts.Image{
				Name:   "imageTest1",
				Digest: "imageDigest1",
			},
			ScanStatus: contracts.UnhealthyScan,
			ScanFindings: []*contracts.ScanFinding{
				{
					Id:        "11",
					Severity:  "High",
					Patchable: true,
				},
			},
		},
		{
			Name: "container2",
			Image: &contracts.Image{
				Name:   "imageTest2",
				Digest: "imageDigest2",
			},
			ScanStatus: contracts.UnhealthyScan,
			ScanFindings: []*contracts.ScanFinding{
				{
					Id:        "22",
					Severity:  "Medium",
					Patchable: true,
				},
			},
		},
	}
	suite.podNoAnnotations = &corev1.Pod{}
	suite.podWithAnnotations = createPodWithAnnotationsForTest()
}

func (suite *TestSuite) Test_CreateContainersVulnerabilityScanAnnotationPatchAdd_TwoContainersScanInfo_AnnotationsGeneratedAsExpected() {
	suite.checkContainersVulnerabilityScanAnnotation(1, suite.podNoAnnotations)
}

func (suite *TestSuite) Test_CreateContainersVulnerabilityScanAnnotationPatchAdd_PodWithAnnotations_AnnotationsGeneratedAsExpected() {

	// check containers vulnerability scan annotations
	mapAnnotations := suite.checkContainersVulnerabilityScanAnnotation(3, suite.podWithAnnotations)

	// check no override of existing annotations
	suite.checkNoOverrideOfExistingAnnotations(mapAnnotations, _annotationTestKeyOne, _annotationTestValueOne)
	suite.checkNoOverrideOfExistingAnnotations(mapAnnotations, _annotationTestKeyTwo, _annotationTestValueTwo)
}

func (suite *TestSuite)checkContainersVulnerabilityScanAnnotation(patchLen int, pod *corev1.Pod) map[string]string{
	result, err := CreateContainersVulnerabilityScanAnnotationPatchAdd(suite.containersScanInfo, pod)
	suite.Nil(err)
	suite.Equal(_expectedTestAddPatchOperation, result.Operation)
	suite.Equal(_expectedTestAnnotationPatchPath, result.Path)

	mapAnnotations, ok := result.Value.(map[string]string)
	suite.True(ok)

	suite.Equal(patchLen, len(mapAnnotations))
	strContainersVulnerabilityScanValue, ok := mapAnnotations[contracts.ContainersVulnerabilityScanInfoAnnotationName]
	suite.True(ok)

	// Unmarshal
	scanInfoList := new(contracts.ContainerVulnerabilityScanInfoList)
	err = json.Unmarshal([]byte(strContainersVulnerabilityScanValue), scanInfoList)
	suite.Nil(err)

	// Verify timestamp
	diff := time.Now().UTC().Sub(scanInfoList.GeneratedTimestamp)
	suite.True((diff >= 0 && diff < time.Second))
	suite.Equal(time.UTC, scanInfoList.GeneratedTimestamp.Location())
	return mapAnnotations
}

func (suite *TestSuite)checkNoOverrideOfExistingAnnotations(mapAnnotations map[string]string, expectedKey string, expectedVal string){
	strAnnotationField1, ok := mapAnnotations[expectedKey]
	suite.True(ok)
	suite.Equal(expectedVal, strAnnotationField1)
}

func TestCreateContainersVulnerabilityScanAnnotationPatchAdd(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func createPodWithAnnotationsForTest() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "podTest",
			Annotations: map[string]string{
				_annotationTestKeyOne : _annotationTestValueOne,
				_annotationTestKeyTwo : _annotationTestValueTwo,
			},
		},
		TypeMeta: metav1.TypeMeta{},
		Spec: corev1.PodSpec{},
	}
}
