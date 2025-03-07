# Skipping Steps

Chainsaw allows you to conditionally skip steps in your tests. This is useful when you want to run different steps based on environment variables or other conditions.

## Basic Usage

To skip a step, add the `skip` field to your step definition:

```yaml
steps:
  - name: Check if resource is deployed
    skip: true
    try:
      - assert:
          file: ../assert-resources.yaml
```

When the `skip` field is set to `true`, Chainsaw will skip the step and continue with the next one.

## Dynamic Skipping with Templates

You can also use template expressions to dynamically determine whether to skip a step:

```yaml
steps:
  - name: Check if resource is deployed
    skip: "{{ .SKIP_STEP }}"
    try:
      - assert:
          file: ../assert-resources.yaml
```

In this example, the step will be skipped if the `SKIP_STEP` environment variable is set to `true`.

## Conditional Skipping

You can use more complex expressions to conditionally skip steps:

```yaml
steps:
  - name: Skip in development environment
    skip: "{{ eq .ENV \"development\" }}"
    try:
      - apply:
          file: ../production-resources.yaml
```

This step will be skipped if the `ENV` environment variable is set to `development`.

## Notes

- The `skip` field accepts a string value that will be evaluated as a boolean.
- Valid values are `true`, `false`, `"true"`, `"false"`, or any template expression that evaluates to a boolean.
- If the template expression cannot be evaluated or does not result in a valid boolean value, an error will be reported.
- Skipped steps will be logged with a "SKIP" status in the test output. 