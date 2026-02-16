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

func TestIssue112(t *testing.T) {
	// https://github.com/josephburnett/jd/issues/112
	// Arrays with more than 10 elements triggered the Myers diff
	// algorithm which had an off-by-one in its backtracking logic.
	ctx := newTestContext(t)
	checkDiff(ctx,
		`{"key1":["v01","v02","v03","v04","v05","v06","v07","v08","v09","v10","v11"]}`,
		`{"key1":["v01","v02","v03","v04","v05","v06","v07","v08","v09","v10","v11 "]}`,
		`@ ["key1",10]`,
		`  "v10"`,
		`- "v11"`,
		`+ "v11 "`,
		`]`,
	)
}

func TestDebug(t *testing.T) {
	fuzz(t, `0`, ``, 0)
}
