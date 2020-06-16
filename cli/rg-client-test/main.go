package main

import (
	"log"

	rgclient "github.com/szuecs/routegroup-client"
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
}
