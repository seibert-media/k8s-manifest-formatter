package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

func TestFormatter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8s Manifest Formatter Suite")
}

var invalidContent = "foo bar"

var formattedContent = `apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: traefik
    traefik.frontend.priority: "10000"
  creationTimestamp: null
  labels:
    app: webdav
  name: webdav
  namespace: webdav
spec:
  rules:
  - host: webdav.benjamin-borbe.de
    http:
      paths:
      - backend:
          serviceName: webdav
          servicePort: web
        path: /
status:
  loadBalancer: {}
`

var unformattedContent = `apiVersion: extensions/v1beta1
kind: Ingress
spec:
  rules:
  - host: webdav.benjamin-borbe.de
    http:
      paths:
      - backend:
          serviceName: webdav
          servicePort: web
        path: /
metadata:
  labels:
    app: webdav
  name: webdav
  namespace: webdav
  annotations:
    traefik.frontend.priority: "10000"
    kubernetes.io/ingress.class: traefik
`
var _ = Describe("Formatter", func() {
	Context("invalid content", func() {
		It("return error", func() {
			_, err := formatYaml([]byte(invalidContent))
			Expect(err).NotTo(BeNil())
		})
	})
	Context("valid but not formatted content", func() {
		It("return no error", func() {
			_, err := formatYaml([]byte(unformattedContent))
			Expect(err).To(BeNil())
		})
		It("output is not empty", func() {
			output, _ := formatYaml([]byte(unformattedContent))
			Expect(output).NotTo(HaveLen(0))
		})
		It("output match input", func() {
			output, _ := formatYaml([]byte(unformattedContent))
			Expect(gbytes.BufferWithBytes(output)).To(gbytes.Say(formattedContent))
		})
	})
	Context("valid and formatted content", func() {
		It("return no error", func() {
			_, err := formatYaml([]byte(formattedContent))
			Expect(err).To(BeNil())
		})
		It("output is not empty", func() {
			output, _ := formatYaml([]byte(formattedContent))
			Expect(output).NotTo(HaveLen(0))
		})
		It("output match input", func() {
			output, _ := formatYaml([]byte(formattedContent))
			Expect(gbytes.BufferWithBytes(output)).To(gbytes.Say(formattedContent))
		})
	})
})
