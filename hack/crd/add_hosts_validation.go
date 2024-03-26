// This program adds host validation to CRD yaml.
//
// # Why
//
// controller-gen does not support validating internal list items on list types,
// see https://github.com/kubernetes-sigs/controller-tools/issues/342
package main

import (
	"log"
	"os"

	"sigs.k8s.io/yaml"
)

const (
	hostPattern   = `^[a-z0-9]([-a-z0-9]*[a-z0-9])?([.][a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	hostMaxLength = 255 // https://datatracker.ietf.org/doc/html/rfc1035#section-2.3.4
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mustGet(v []byte, err error) []byte {
	must(err)
	return v
}

type setter struct {
	o interface{}
}

func (s *setter) field(name string) *setter {
	return &setter{s.o.(map[string]interface{})[name]}
}

func (s *setter) item(i int) *setter {
	return &setter{s.o.([]interface{})[i]}
}

func (s *setter) setField(name string, value interface{}) *setter {
	s.o.(map[string]interface{})[name] = value
	return s
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("CRD filename required")
	}
	crdFilename := os.Args[1]

	yamlBytes := mustGet(os.ReadFile(crdFilename))

	o := make(map[string]interface{})
	must(yaml.Unmarshal(yamlBytes, &o))

	s := &setter{o}

	s.field("spec").
		field("versions").
		item(0).
		field("schema").
		field("openAPIV3Schema").
		field("properties").
		field("spec").
		field("properties").
		field("hosts").
		field("items").
		setField("pattern", hostPattern).
		setField("maxLength", hostMaxLength)

	s.field("spec").
		field("versions").
		item(0).
		field("schema").
		field("openAPIV3Schema").
		field("properties").
		field("spec").
		field("properties").
		field("tls").
		field("items").
		field("properties").
		field("hosts").
		field("items").
		setField("pattern", hostPattern).
		setField("maxLength", hostMaxLength)

	outYaml := mustGet(yaml.Marshal(o))
	must(os.WriteFile(crdFilename, outYaml, 0664))
}
