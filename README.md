# routegroup-client
client-go implementation of
[routegroup](https://opensource.zalando.com/skipper/kubernetes/routegroups/).

There is example code in ./cli/rg-client-test to show how to use
it. If you want to build the example you can run `make`, which will do
code generation and build this example client application, that you
can run from build/ directory with `./build/rg-example`. You can use
`kubectl proxy` to run the example against the local API endpoint.

Unified Client does composition to have one client for all Kubernetes
API access. In case you want to do composition yourself you can use
the Zalando Client.

## Unified Client

A unified client has all the functions from `kubernetes.Interface` and
the `ZalandoV1()` available in one client-go client.

Sample code:

```go
import (
	"context"
	"log"

	rgclient "github.com/szuecs/routegroup-client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	cli, err := rgclient.CreateUnified()
	if err != nil {
		log.Fatalf("Failed to create unified RouteGroup client: %v", err)
	}

	// example kubernetes.Interface access with unified client
	ings, err := cli.ExtensionsV1beta1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	for _, ing := range ings.Items {
		log.Printf("ing Namespace/Name: %s/%s", ing.Namespace, ing.Name)
	}

	// example RouteGroups access with unified client
	l, err := cli.ZalandoV1().RouteGroups("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get RouteGroup list: %v", err)
	}

	for _, rg := range l.Items {
		log.Printf("rg Namespace/Name: %s/%s", rg.Namespace, rg.Name)
	}
```

## Zalando Client

Zalando client has access to `ZalandoV1()`, but does not have access
to `kubernetes.Interface`, such that you can't list Ingress for
example with it.

Sample code:

```go
import (
	"context"
	"log"

	rgclient "github.com/szuecs/routegroup-client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	zcli, err := rgclient.Create()
	if err != nil {
		log.Fatalf("Failed to create RouteGroup client: %v", err)
	}

	ls, err := zcli.ZalandoV1().RouteGroups("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list routegroups: %v", err)
	}

	for _, rg := range l.Items {
		log.Printf("rg Namespace/Name: %s/%s", rg.Namespace, rg.Name)
	}
```
