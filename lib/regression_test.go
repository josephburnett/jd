package jd

import "testing"

func TestIssue25(t *testing.T) {
	// https://github.com/josephburnett/jd/issues/25
	aNode, _ := ReadYamlString(`
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx:1.14.2
        name: nginx
        ports:
        - containerPort: 8080
`)
	patch, _ := ReadDiffString(`
@ ["spec","template","spec","containers",{"name":"nginx"},"ports",0]
- 8080
+ 8081
`)
	_, err := aNode.Patch(patch)
	if err != nil {
		t.Errorf("wanted no err. got %v", err)
	}
}
