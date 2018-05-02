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

var _ = Describe("Formatter", func() {
	Context("valid content", func() {
		content := `apiVersion: extensions/v1beta1
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
		It("return no error", func() {
			_, err := formatYaml([]byte(content))
			Expect(err).To(BeNil())
		})
		It("output not empty no error", func() {
			output, _ := formatYaml([]byte(content))
			Expect(output).NotTo(HaveLen(0))
		})
		It("output match input", func() {
			output, _ := formatYaml([]byte(content))
			Expect(gbytes.BufferWithBytes(output)).To(gbytes.Say(content))
		})
	})
	Context("yaml long line", func() {
		content := `apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  creationTimestamp: null
  labels:
    app: long long long long long long long long long long long long long long long long long long long long long long long long long long long long long long long long long long long long long long long
spec: {}
status:
  loadBalancer: {}
`
		It("return no error", func() {
			_, err := formatYaml([]byte(content))
			Expect(err).To(BeNil())
		})
		It("output not empty no error", func() {
			output, _ := formatYaml([]byte(content))
			Expect(output).NotTo(HaveLen(0))
		})
		It("output match input", func() {
			output, _ := formatYaml([]byte(content))
			Expect(gbytes.BufferWithBytes(output)).To(gbytes.Say(content))
		})
	})
})
