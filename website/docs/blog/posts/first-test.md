---
date: 2023-10-16
slug: first-test
categories:
  - announcements
authors:
  - eddycharly
---

# First working test !

<p align="center">
  <img src="https://www.alertdriving.co.nz/uploads/1/4/8/0/14809288/112642145.jpg" />
</p>

First test run !

<!-- more -->

The [quick start](../../quick-start.md) test is passing :tada:

```bash
# define kind image
export KIND_IMAGE="kindest/node:v1.28.0"

# create cluster
kind create cluster --image $KIND_IMAGE

# create test folder
mkdir quick-start

# enter test folder
cd quick-start

# create a ConfigMap
cat > configmap.yaml << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
data:
  foo: bar
EOF

# create test file
cat > chainsaw-test.yaml << EOF
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: quick-start
spec:
  steps:
  # first step applies the config map
  - apply:
    # file is relative to the test folder
    - file: configmap.yaml
  # second step verifies the config map exists and contains the expected data
  - assert:
    # file is relative to the test folder
    - file: configmap.yaml
EOF

../chainsaw test --test-dir .

Running without configuration
Loading tests...
- quick-start (.)
Running tests...
=== RUN   chainsaw
=== RUN   chainsaw/quick-start
=== PAUSE chainsaw/quick-start
=== CONT  chainsaw/quick-start
    runner.go:35: step-1
    runner.go:35: apply chainsaw-polite-chamois/chainsaw-quick-start (v1/ConfigMap)
    runner.go:35: step-2
    runner.go:35: assert chainsaw-polite-chamois/chainsaw-quick-start (v1/ConfigMap)
    runner.go:68: cleanup namespace: chainsaw-polite-chamois
--- PASS: chainsaw (0.00s)
    --- PASS: chainsaw/quick-start (5.09s)
PASS
Done.
```

Impressive ! :joy:
