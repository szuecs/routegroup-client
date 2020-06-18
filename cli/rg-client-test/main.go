package main

import (
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
	log.Printf("cli: %#v", cli)

	// example kubernetes.Interface access
	ings, err := cli.ExtensionsV1beta1().Ingresses("").List(metav1.ListOptions{})
	//ings, err := cli.NetworkingV1().Ingress("").List(metav1.ListOptions{})
	for _, ing := range ings.Items {
		log.Printf("ing NAmespace/Name: %s/%s", ing.Namespace, ing.Name)
	}

	// example RouteGroups access
	l, err := cli.ZalandoV1().RouteGroups("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get RouteGroup list: %v", err)
	}
	log.Printf("%#v", l)
	for _, rg := range l.Items {
		log.Printf("rg Namespace/Name: %s/%s", rg.Namespace, rg.Name)
		log.Printf("status: %+v", rg.Status)
		log.Printf("spec: %+v", rg.Spec)
	}

	zcli, err := rgclient.Create()
	if err != nil {
		log.Fatalf("Failed to create RouteGroup client: %v", err)
	}
	ls, err := zcli.ZalandoV1().RouteGroups("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list routegroups: %v", err)
	}
	for _, rg := range ls.Items {
		log.Printf("rg Namespace/Name: %s/%s", rg.Namespace, rg.Name)
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
						`status(418) -> inlineContent("I am a teapot")`,
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
	rg, err := zcli.ZalandoV1().RouteGroups(namespace).Create(newRg, metav1.CreateOptions{})
	if err != nil {
		spew.Dump(newRg)
		log.Printf("newRg: %#v", newRg)
		log.Fatalf("Failed to create routegroup: %v", err)
	}
	spew.Dump(rg)
}
