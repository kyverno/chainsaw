# `x_k8s_get`

```yaml
    - description: >-
        1. Look up some details about a Deployment object with 'x_k8s_get(...)', and store them in a binding.
        2. Make use of the binding in a `script` block through environment variables.
        3. Refer to the same bound Deployment object in an `assert` block.
      bindings:
        - name: xpNS
          value: crossplane-system
        - name: xpDeploy
          # Arguments to 'x_k8s_get(any, string, string, string, string)':
          #
          # 1. any: Supply the '$client' built-in binding here, which is needed to connect to the cluster:
          #         https://kyverno.github.io/chainsaw/latest/reference/builtins/
          #
          # 2. string: 'apiVersion' field on the object.
          #
          # 3. string: 'Kind' field on the object.
          #
          # 4. string: The namespace of the object. If the object type is not namespaced, this field can be an empty
          #            string, or any string; it doesn't seem to matter.
          #
          # 5. string: The name of the object.
          value: (x_k8s_get($client, 'apps/v1', 'Deployment', $xpNS, 'crossplane'))
      try:
        - script:
            bindings:
              - # Re-bind the version label from the Deployment to give us a more succinct name to refer to.
                name: deployVersion
                # If the label key has any periods in it, double-quote the whole key.
                value: ($xpDeploy.metadata.labels."app.kubernetes.io/version")
            env:
              - # Refer to the Deployment version label the long way.
                name: DEP_VER_LONG
                value: ($xpDeploy.metadata.labels."app.kubernetes.io/version")
              - # Refer to the Deployment version label the short way, through the additional binding scoped to this
                # 'script' block.
                name: DEP_VER_SHORT
                value: ($deployVersion)
            # The version values printed by the script here will be the same, even though they were derived through
            # slightly different routes.
            content: |-
              echo "DEP_VER_LONG:  '$DEP_VER_LONG'"
              echo "DEP_VER_SHORT: '$DEP_VER_SHORT'"
        - assert:
            bindings:
              - name: depVer
                value: ($xpDeploy.metadata.labels."app.kubernetes.io/version")
            resource:
              apiVersion: v1
              kind: Pod
              metadata:
                # Match the pod(s) with namespace and label selectors.
                namespace: ($xpNS)
                labels:
                  app: crossplane
              # Make a binding that holds this Pod's version from its label.
              (metadata.labels."app.kubernetes.io/version")->podVer:
                # Assert that the version labels for the Deployment and this Pod equal each other, using the
                # 'semver_compare()' function described here:
                # https://kyverno.io/docs/writing-policies/jmespath/#semver_compare
                (semver_compare($depVer, $podVer)): true
```
