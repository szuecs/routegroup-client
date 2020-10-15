// +k8s:deepcopy-gen=package

package v1

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// +k8s:deepcopy-gen=true
// +kubebuilder:resource:shortName=rg
// +kubebuilder:resource:categories="all"
// +kubebuilder:printcolumn:name="Hosts",type=string,JSONPath=`.spec.hosts`,description="Hosts defined for the RouteGroup"
// +kubebuilder:printcolumn:name="Address",type=string,JSONPath=`.status.loadBalancer`,description="Address of the Load Balancer for the RouteGroup"
// +kubebuilder:subresource:status
type RouteGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouteGroupSpec   `json:"spec"`
	Status RouteGroupStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:deepcopy-gen=true
type RouteGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []RouteGroup `json:"items"`
}

// +k8s:deepcopy-gen=true
type RouteGroupSpec struct {
	// List of hostnames for the RouteGroup.
	Hosts []string `json:"hosts,omitempty"`
	// List of backends that can be referenced in the routes.
	Backends []RouteGroupBackend `json:"backends"`
	// DefaultBackends is a list of default backends defined if no explicit
	// backend is defined for a route.
	DefaultBackends []RouteGroupBackendReference `json:"defaultBackends,omitempty"`
	// +kubebuilder:validation:MinItems=1
	Routes []RouteGroupRouteSpec `json:"routes,omitempty"`
}

// RouteGroupBackendType is the type of the route group backend.
// +kubebuilder:validation:Enum=service;shunt;loopback;dynamic;lb;network
type RouteGroupBackendType string

const (
	ServiceRouteGroupBackend  RouteGroupBackendType = "service"
	ShuntRouteGroupBackend    RouteGroupBackendType = "shunt"
	LoopbackRouteGroupBackend RouteGroupBackendType = "loopback"
	DynamicRouteGroupBackend  RouteGroupBackendType = "dynamic"
	LBRouteGroupBackend       RouteGroupBackendType = "lb"
	NetworkRouteGroupBackend  RouteGroupBackendType = "network"
)

// +k8s:deepcopy-gen=true
type RouteGroupBackend struct {
	// Name is the BackendName that can be referenced as RouteGroupBackendReference
	Name string `json:"name"`
	// Type is one of "service|shunt|loopback|dynamic|lb|network"
	Type RouteGroupBackendType `json:"type"`
	// Address is required for Type network
	// +optional
	Address string `json:"address,omitempty"`
	// Algorithm is required for Type lb
	// +optional
	Algorithm string `json:"algorithm,omitempty"`
	// Endpoints is required for Type lb
	// +kubebuilder:validation:MinItems=1
	Endpoints []string `json:"endpoints,omitempty"`
	// ServiceName is required for Type service
	// +optional
	ServiceName string `json:"serviceName,omitempty"`
	// ServicePort is required for Type service
	// +optional
	ServicePort int `json:"servicePort,omitempty"`
}

// +k8s:deepcopy-gen=true
type RouteGroupBackendReference struct {
	// BackendName references the skipperBackend by name
	BackendName string `json:"backendName"`
	// Weight defines the traffic weight, if there are 2 or more
	// default backends
	// +kubebuilder:validation:Minimum=0
	// +optional
	Weight int `json:"weight"`
}

// HTTPMethod is a valid HTTP method in uppercase.
// +kubebuilder:validation:Enum=GET;HEAD;POST;PUT;PATCH;DELETE;CONNECT;OPTIONS;TRACE
type HTTPMethod string

const (
	MethodGet     HTTPMethod = http.MethodGet
	MethodHead    HTTPMethod = http.MethodHead
	MethodPost    HTTPMethod = http.MethodPost
	MethodPut     HTTPMethod = http.MethodPut
	MethodPatch   HTTPMethod = http.MethodPatch
	MethodDelete  HTTPMethod = http.MethodDelete
	MethodConnect HTTPMethod = http.MethodConnect
	MethodOptions HTTPMethod = http.MethodOptions
	MethodTrace   HTTPMethod = http.MethodTrace
)

// +k8s:deepcopy-gen=true
type RouteGroupRouteSpec struct {
	// Path specifies Path predicate, only one of Path or PathSubtree is allowed
	Path string `json:"path,omitempty"`

	// PathSubtree specifies PathSubtree predicate, only one of Path or PathSubtree is allowed
	PathSubtree string `json:"pathSubtree,omitempty"`

	// PathRegexp can be added additionally
	PathRegexp string `json:"pathRegexp,omitempty"`

	// RouteGroupBackendReference specifies the list of backendReference that should
	// be applied to override the defaultBackends
	// +optional
	Backends []RouteGroupBackendReference `json:"backends,omitempty"`

	// Filters specifies the list of filters applied to the routeSpec
	// +optional
	Filters []string `json:"filters,omitempty"`

	// Predicates specifies the list of predicates applied to the routeSpec
	// +optional
	Predicates []string `json:"predicates,omitempty"`

	// Methods defines valid HTTP methods for the specified routeSpec
	// +optional
	Methods []HTTPMethod `json:"methods,omitempty"`
}

// +k8s:deepcopy-gen=true
type RouteGroupStatus struct {
	// LoadBalancer is similar to ingress status, such that
	// external-dns has the same style as in ingress
	LoadBalancer RouteGroupLoadBalancerStatus `json:"loadBalancer"`
}

// +k8s:deepcopy-gen=true
type RouteGroupLoadBalancerStatus struct {
	// RouteGroup is similar to Ingress in ingress status.LoadBalancer.
	RouteGroup []RouteGroupLoadBalancer `json:"routegroup"`
}

// +k8s:deepcopy-gen=true
type RouteGroupLoadBalancer struct {
	// IP is the IP address of the load balancer and is empty if Hostname is set
	IP string `json:"ip,omitempty"`
	// Hostname is the hostname of the load balancer and is empty if IP is set
	Hostname string `json:"hostname,omitempty"`
}
