// +k8s:deepcopy-gen=package

package v1

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KindRouteGroup     = "RouteGroup"
	KindRouteGroupList = "RouteGroupList"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// +k8s:deepcopy-gen=true
// +kubebuilder:resource:categories="all",shortName=rg;rgs
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
	// List of hostnames for the RouteGroup
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:items:MaxLength=255
	// +kubebuilder:validation:items:Pattern="^[a-z0-9]([-a-z0-9]*[a-z0-9])?([.][a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"
	// +listType=set
	Hosts []string `json:"hosts,omitempty"`
	// List of backends that can be referenced in the routes
	Backends []RouteGroupBackend `json:"backends"`
	// DefaultBackends is a list of default backends defined if no explicit
	// backend is defined for a route
	DefaultBackends []RouteGroupBackendReference `json:"defaultBackends,omitempty"`
	// Routes describe how a matching HTTP request is handled and where it is forwarded to
	// +kubebuilder:validation:MinItems=1
	Routes []RouteGroupRouteSpec `json:"routes,omitempty"`
	// TLS defines which Kubernetes secret will be used to terminate the connection
	// based on the matching hostnames
	// +optional
	TLS []RouteGroupTLSSpec `json:"tls,omitempty"`
}

// RouteGroupBackendType is the type of the route group backend.
type RouteGroupBackendType string

const (
	ServiceRouteGroupBackend  RouteGroupBackendType = "service"
	ShuntRouteGroupBackend    RouteGroupBackendType = "shunt"
	LoopbackRouteGroupBackend RouteGroupBackendType = "loopback"
	DynamicRouteGroupBackend  RouteGroupBackendType = "dynamic"
	LBRouteGroupBackend       RouteGroupBackendType = "lb"
	NetworkRouteGroupBackend  RouteGroupBackendType = "network"
	ForwardRouteGroupBackend  RouteGroupBackendType = "forward"
)

// BackendAlgorithmType is the type of algorithm used for load balancing
// traffic to a backend. This is only valid for backend type lb|service.
type BackendAlgorithmType string

const (
	RoundRobinBackendAlgorithm            BackendAlgorithmType = "roundRobin"
	RandomBackendAlgorithm                BackendAlgorithmType = "random"
	ConsistentHashBackendAlgorithm        BackendAlgorithmType = "consistentHash"
	PowerOfRandomNChoicesBackendAlgorithm BackendAlgorithmType = "powerOfRandomNChoices"
)

// +k8s:deepcopy-gen=true
type RouteGroupBackend struct {
	// Name is the BackendName that can be referenced as RouteGroupBackendReference
	Name string `json:"name"`
	// Type of the backend.
	// `service`- resolve Kubernetes service to the available Endpoints belonging to the Service, and generate load balanced routes using them.
	// `shunt` - reply directly from the proxy itself. This can be used to shortcut, for example have a default that replies with 404 or use skipper as a backend serving static content in demos.
	// `loopback` - lookup again the routing table to a better matching route after processing the current route. Like this you can add some headers or change the request path for some specific matching requests.
	// `dynamic` - use the backend provided by filters. This allows skipper as library users to do proxy calls to a certain target from their own implementation dynamically looked up by their filters.
	// `lb` - balance the load across multiple network endpoints using specified algorithm. If algorithm is not specified it will use the default algorithm set by Skipper at start.
	// `network` - use arbitrary HTTP or HTTPS URL.
	// `forward` - replaced by a network backend chosen by skipper -forward-backend-url.
	// +kubebuilder:validation:Enum=service;shunt;loopback;dynamic;lb;network;forward
	Type RouteGroupBackendType `json:"type"`
	// Address is required for type `network`
	// +optional
	Address string `json:"address,omitempty"`
	// Algorithm is required for type `lb`.
	// `roundRobin` - backend is chosen by the round robin algorithm, starting with a random selected backend to spread across all backends from the beginning.
	// `random` - backend is chosen at random.
	// `consistentHash` - backend is chosen by [consistent hashing](https://en.wikipedia.org/wiki/Consistent_hashing) algorithm based on the request key. The request key is derived from `X-Forwarded-For` header or request remote IP address as the fallback. Use [`consistentHashKey`](filters.md#consistenthashkey) filter to set the request key. Use [`consistentHashBalanceFactor`](filters.md#consistenthashbalancefactor) to prevent popular keys from overloading a single backend endpoint.
	// `powerOfRandomNChoices` - backend is chosen by selecting N random endpoints and picking the one with least outstanding requests from them (see http://www.eecs.harvard.edu/~michaelm/postscripts/handbook2001.pdf).
	// +kubebuilder:validation:Enum=roundRobin;random;consistentHash;powerOfRandomNChoices
	// +optional
	Algorithm BackendAlgorithmType `json:"algorithm,omitempty"`
	// Endpoints is required for type `lb`
	// +kubebuilder:validation:MinItems=1
	Endpoints []string `json:"endpoints,omitempty"`
	// ServiceName is required for type `service`
	// +optional
	ServiceName string `json:"serviceName,omitempty"`
	// ServicePort is required for type `service`
	// +optional
	ServicePort int `json:"servicePort,omitempty"`
}

// +k8s:deepcopy-gen=true
type RouteGroupBackendReference struct {
	// BackendName references backend by name
	BackendName string `json:"backendName"`
	// Weight defines a portion of traffic for the referenced backend.
	// It equals to weight divided by the sum of all backend weights.
	// When all references have zero (or unspecified) weight then traffic is split equally between them.
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
type RouteGroupTLSSpec struct {
	// TLS hosts specify the list of hosts included in the TLS secret.
	// The values in this list must match the host name(s) used for
	// the RouteGroup in order to terminate TLS for the host(s).
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:items:MaxLength=255
	// +kubebuilder:validation:items:Pattern="^[a-z0-9]([-a-z0-9]*[a-z0-9])?([.][a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"
	// +listType=set
	Hosts []string `json:"hosts"`

	// SecretName is the name of the secret used to terminate TLS traffic.
	// Secret should reside in the same namespace as the RouteGroup.
	SecretName string `json:"secretName"`
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
