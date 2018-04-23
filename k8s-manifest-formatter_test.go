package main

import "testing"

func TestFormatYamlIlllega(t *testing.T) {
	_, err := formatYaml([]byte(`illegal content`))
	if err == nil {
		t.Fatal("err expected")
	}
}

func TestFormatYaml(t *testing.T) {
	content:=`apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: traefik
    traefik.frontend.priority: '10000'
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
`
	_, err := formatYaml([]byte(content))
	if err != nil {
		t.Fatal("err not expected")
	}
}
