package main

import (
	"encoding/json"
	"log"

	"github.com/davecgh/go-spew/spew"
	rgclient "github.com/szuecs/routegroup-client"
	rgv1 "github.com/szuecs/routegroup-client/apis/zalando.org/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	cli, err := rgclient.CreateUnified()
	if err != nil {
		log.Fatalf("Failed to create unified RouteGroup client: %v", err)
	}
	log.Println("have cli")

	// example kubernetes.Interface access
	ings, err := cli.ExtensionsV1beta1().Ingresses("").List(metav1.ListOptions{})
	//ings, err := cli.NetworkingV1().Ingress("").List(metav1.ListOptions{})
	for _, ing := range ings.Items {
		log.Printf("ing NAmespace/Name: %s/%s", ing.Namespace, ing.Name)
	}
	log.Printf("have ing %d", len(ings.Items))

	// example RouteGroups access
	l, err := cli.ZalandoV1().RouteGroups("").List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to get RouteGroup list: %v", err)
	} else {
		for _, rg := range l.Items {
			log.Printf("rg Namespace/Name: %s/%s", rg.Namespace, rg.Name)
			log.Printf("status: %+v", rg.Status)
			log.Printf("spec: %+v", rg.Spec)
		}
	}

	// test create
	name := "myrg"
	namespace := "default"
	hostname := "myrg.teapot-e2e.zalan.do"
	port := 83
	newRg := &rgv1.RouteGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: rgv1.RouteGroupSpec{
			Hosts: []string{hostname},
			Backends: []rgv1.RouteGroupBackend{
				{
					Name:        name,
					Type:        "service",
					ServiceName: name,
					ServicePort: port,
				},
				{
					Name: "router",
					Type: "shunt",
				},
			},
			DefaultBackends: []rgv1.RouteGroupBackendReference{
				{
					BackendName: name,
					Weight:      1,
				},
			},
			Routes: []rgv1.RouteGroupRouteSpec{
				{
					PathSubtree: "/",
				},
				{
					PathSubtree: "/router-response",
					Filters: []string{
						`status(418)`,
						`inlineContent("I am a teapot")`,
					},
					Backends: []rgv1.RouteGroupBackendReference{
						{
							BackendName: "router",
							Weight:      1,
						},
					},
				},
			},
		},
	}

	b, err := json.Marshal(newRg)
	log.Printf("b: %s", string(b))
	if err != nil {
		log.Fatalf("json marshal failed: %v", err)
	}

	//rg, err := cli.ZalandoV1().RouteGroups(namespace).Create(newRg, metav1.CreateOptions{})
	rg, err := cli.ZalandoV1().RouteGroups(namespace).Create(newRg)
	if err != nil {
		spew.Dump(newRg)
		log.Fatalf("Failed to create routegroup: %v", err)
	}
	//spew.Dump(rg)
	log.Printf("Created %s/%s", namespace, rg.Name)
}
