package v1beta1

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DynaKubeStatus defines the observed state of DynaKube
// +k8s:openapi-gen=true
type DynaKubeStatus struct {
	// Defines the current state (Running, Updating, Error, ...)
	Phase DynaKubePhaseType `json:"phase,omitempty"`

	// UpdatedTimestamp indicates when the instance was last updated
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Last Updated"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:text"
	UpdatedTimestamp metav1.Time `json:"updatedTimestamp,omitempty"`

	// Deprecated: use DynatraceApiStatus.LastTokenScopeRequest instead
	// LastTokenProbeTimestamp tracks when the last request for the API token validity was sent
	LastTokenProbeTimestamp *metav1.Time `json:"lastTokenProbeTimestamp,omitempty"`

	// KubeSystemUUID contains the UUID of the current Kubernetes cluster
	KubeSystemUUID string `json:"kubeSystemUUID,omitempty"`

	// Conditions includes status about the current state of the instance
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	ActiveGate   ActiveGateStatus   `json:"activeGate,omitempty"`
	OneAgent     OneAgentStatus     `json:"oneAgent,omitempty"`
	CodeModules  CodeModulesStatus  `json:"codeModules,omitempty"`
	Synthetic    SyntheticStatus    `json:"synthetic,omitempty"`
	DynatraceApi DynatraceApiStatus `json:"dynatraceApi,omitempty"`
}

type DynatraceApiStatus struct {
	LastTokenScopeRequest               metav1.Time `json:"lastTokenScopeRequest,omitempty"`
	LastOneAgentConnectionInfoRequest   metav1.Time `json:"lastOneAgentConnectionInfoRequest,omitempty"`
	LastActiveGateConnectionInfoRequest metav1.Time `json:"lastActiveGateConnectionInfoRequest,omitempty"`
}

func (dynatraceApiStatus *DynatraceApiStatus) ResetCachedTimestamps() {
	dynatraceApiStatus.LastTokenScopeRequest = metav1.Time{}
	dynatraceApiStatus.LastOneAgentConnectionInfoRequest = metav1.Time{}
	dynatraceApiStatus.LastActiveGateConnectionInfoRequest = metav1.Time{}
}

func GetCacheValidMessage(functionName string, lastRequestTimestamp metav1.Time, timeout time.Duration) string {
	remaining := timeout - time.Since(lastRequestTimestamp.Time)

	return fmt.Sprintf("skipping %s, last request was made less than %d minutes ago, %d minutes remaining until next request",
		functionName,
		int(timeout.Minutes()),
		int(remaining.Minutes()))
}

type ConnectionInfoStatus struct {
	CommunicationHosts              []CommunicationHostStatus `json:"communicationHosts,omitempty"`
	TenantUUID                      string                    `json:"tenantUUID,omitempty"`
	FormattedCommunicationEndpoints string                    `json:"formattedCommunicationEndpoints,omitempty"`
}

type CommunicationHostStatus struct {
	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     uint32 `json:"port,omitempty"`
}

type VersionSource string

const (
	TenantRegistryVersionSource VersionSource = "tenant-registry"
	CustomImageVersionSource    VersionSource = "custom-image"
	CustomVersionVersionSource  VersionSource = "custom-version"
	PublicRegistryVersionSource VersionSource = "public-registry"
)

type VersionStatus struct {
	Source             VersionSource `json:"source,omitempty"`
	ImageHash          string        `json:"imageHash,omitempty"`
	ImageRepository    string        `json:"imageRepository,omitempty"`
	ImageTag           string        `json:"imageTag,omitempty"`
	Version            string        `json:"version,omitempty"`
	LastProbeTimestamp *metav1.Time  `json:"lastProbeTimestamp,omitempty"`
}

type ActiveGateStatus struct {
	VersionStatus `json:",inline"`
}

type CodeModulesStatus struct {
	VersionStatus `json:",inline"`
}

type OneAgentStatus struct {
	VersionStatus `json:",inline"`

	Instances map[string]OneAgentInstance `json:"instances,omitempty"`

	LastInstanceStatusUpdate *metav1.Time `json:"lastInstanceStatusUpdate,omitempty"`
}

type OneAgentInstance struct {
	PodName   string `json:"podName,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
}

type SyntheticStatus struct {
	VersionStatus `json:",inline"`
}

type DynaKubePhaseType string

const (
	Running   DynaKubePhaseType = "Running"
	Deploying DynaKubePhaseType = "Deploying"
	Error     DynaKubePhaseType = "Error"
)

// SetPhase sets the status phase on the DynaKube object
func (dk *DynaKubeStatus) SetPhase(phase DynaKubePhaseType) bool {
	upd := phase != dk.Phase
	dk.Phase = phase
	return upd
}

// SetPhaseOnError fills the phase with the Error value in case of any error
func (dk *DynaKubeStatus) SetPhaseOnError(err error) bool {
	if err != nil {
		return dk.SetPhase(Error)
	}
	return false
}
