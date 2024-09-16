# Functions

!!! warning "Experimental functions"

    Experimental functions are denoted by the `x_` prefix.

    These are functions that are subject to signature change in a future version.

## built-in functions

| Name | Signature |
|---|---|
| abs | `abs(number)` |
| avg | `avg(array[number])` |
| ceil | `ceil(number)` |
| contains | `contains(array|string, any)` |
| ends_with | `ends_with(string, string)` |
| find_first | `find_first(string, string, number, number)` |
| find_last | `find_last(string, string, number, number)` |
| floor | `floor(number)` |
| from_items | `from_items(array[array])` |
| group_by | `group_by(array, expref)` |
| items | `items(object)` |
| join | `join(string, array[string])` |
| keys | `keys(object)` |
| length | `length(string|array|object)` |
| lower | `lower(string)` |
| map | `map(expref, array)` |
| max | `max(array[number]|array[string])` |
| max_by | `max_by(array, expref)` |
| merge | `merge(object)` |
| min | `min(array[number]|array[string])` |
| min_by | `min_by(array, expref)` |
| not_null | `not_null(any)` |
| pad_left | `pad_left(string, number, string)` |
| pad_right | `pad_right(string, number, string)` |
| replace | `replace(string, string, string, number)` |
| reverse | `reverse(array|string)` |
| sort | `sort(array[string]|array[number])` |
| sort_by | `sort_by(array, expref)` |
| split | `split(string, string, number)` |
| starts_with | `starts_with(string, string)` |
| sum | `sum(array[number])` |
| to_array | `to_array(any)` |
| to_number | `to_number(any)` |
| to_string | `to_string(any)` |
| trim | `trim(string, string)` |
| trim_left | `trim_left(string, string)` |
| trim_right | `trim_right(string, string)` |
| type | `type(any)` |
| upper | `upper(string)` |
| values | `values(object)` |
| zip | `zip(array, array)` |

## kyverno-json functions

| Name | Signature |
|---|---|
| at | `at(array, any)` |
| concat | `concat(string, string)` |
| json_parse | `json_parse(string)` |
| wildcard | `wildcard(string, string)` |

## kyverno functions

| Name | Signature |
|---|---|
| compare | `compare(string, string)` |
| equal_fold | `equal_fold(string, string)` |
| replace | `replace(string, string, string, number)` |
| replace_all | `replace_all(string, string, string)` |
| to_upper | `to_upper(string)` |
| to_lower | `to_lower(string)` |
| trim | `trim(string, string)` |
| trim_prefix | `trim_prefix(string, string)` |
| split | `split(string, string)` |
| regex_replace_all | `regex_replace_all(string, string|number, string|number)` |
| regex_replace_all_literal | `regex_replace_all_literal(string, string|number, string|number)` |
| regex_match | `regex_match(string, string|number)` |
| pattern_match | `pattern_match(string, string|number)` |
| label_match | `label_match(object, object)` |
| to_boolean | `to_boolean(string)` |
| add | `add(any, any)` |
| sum | `sum(array)` |
| subtract | `subtract(any, any)` |
| multiply | `multiply(any, any)` |
| divide | `divide(any, any)` |
| modulo | `modulo(any, any)` |
| round | `round(number, number)` |
| base64_decode | `base64_decode(string)` |
| base64_encode | `base64_encode(string)` |
| time_since | `time_since(string, string, string)` |
| time_now | `time_now()` |
| time_now_utc | `time_now_utc()` |
| path_canonicalize | `path_canonicalize(string)` |
| truncate | `truncate(string, number)` |
| semver_compare | `semver_compare(string, string)` |
| parse_json | `parse_json(string)` |
| parse_yaml | `parse_yaml(string)` |
| lookup | `lookup(object|array, string|number)` |
| items | `items(object|array, string, string)` |
| object_from_lists | `object_from_lists(array, array)` |
| random | `random(string)` |
| x509_decode | `x509_decode(string)` |
| time_to_cron | `time_to_cron(string)` |
| time_add | `time_add(string, string)` |
| time_parse | `time_parse(string, string)` |
| time_utc | `time_utc(string)` |
| time_diff | `time_diff(string, string)` |
| time_before | `time_before(string, string)` |
| time_after | `time_after(string, string)` |
| time_between | `time_between(string, string, string)` |
| time_truncate | `time_truncate(string, string)` |

## chainsaw functions

| Name | Signature |
|---|---|
| env | `env(string)` |
| x_k8s_get | `x_k8s_get(any, string, string, string, string)` |
| x_k8s_list | `x_k8s_list(any, string, string, string)` |
| x_k8s_exists | `x_k8s_exists(any, string, string, string, string)` |
| x_k8s_resource_exists | `x_k8s_resource_exists(any, string, string)` |
| x_k8s_server_version | `x_k8s_server_version(any)` |
| x_metrics_decode | `x_metrics_decode(string)` |
| trim_space | `trim_space(string)` |
| as_string | `as_string(any)` |

## examples

- [x_k8s_get](./examples/x_k8s_get.md)

### compare

compares two strings lexicographically

### equal_fold

allows comparing two strings for equivalency where the only differences are letter cases

### replace

replaces a specified number of instances of the source string with the replacement string in a parent 

### replace_all

replace all instances of one string with another in an overall parent string

### to_upper

takes in a string and outputs the same string with all upper-case letters

### to_lower

takes in a string and outputs the same string with all lower-case letters

### trim

trims both ends of the source string by characters appearing in the second string

### trim_prefix

trims the second string prefix from the first string if the first string starts with the prefix

### split

splits the first string when the second string is found and converts it into an array 

### regex_replace_all

converts all parameters to string

### regex_replace_all_literal

converts all parameters to string

### regex_match

first string is the regular exression which is compared with second input which can be a number or string

### pattern_match

'*' matches zero or more alphanumeric characters, '?' matches a single alphanumeric character

### label_match

object arguments must be enclosed in backticks; ex. `{{request.object.spec.template.metadata.labels}}`

### to_boolean

It returns true or false for any string, such as 'True', 'TruE', 'False', 'FAlse', 'faLSE', etc.

### add

does arithmetic addition of two specified values of numbers, quantities, and durations

### sum

does arithmetic addition of specified array of values of numbers, quantities, and durations

### subtract

does arithmetic subtraction of two specified values of numbers, quantities, and durations

### multiply

does arithmetic multiplication of two specified values of numbers, quantities, and durations

### divide

divisor must be non zero

### modulo

divisor must be non-zero, arguments must be integers

### round

does roundoff to upto the given decimal places

### base64_decode

decodes a base 64 string

### base64_encode

encodes a regular, plaintext and unencoded string to base64

### time_since

calculate the difference between a start and end period of time where the end may either be a static definition or the then-current time

### time_now

returns current time in RFC 3339 format

### time_now_utc

returns current UTC time in RFC 3339 format

### path_canonicalize

normalizes or canonicalizes a given path by removing excess slashes

### truncate

length argument must be enclosed in backticks; ex. "{{request.object.metadata.name | truncate(@, `9`)}}"

### semver_compare

compares two strings which comply with the semantic versioning schema and outputs a boolean response as to the position of the second relative to the first

### parse_json

decodes a valid JSON encoded string to the appropriate type. Opposite of `to_string` function

### parse_yaml

decodes a valid YAML encoded string to the appropriate type provided it can be represented as JSON

### lookup

returns the value corresponding to the given key/index in the given object/array

### items

converts a map or array to an array of objects where each key:value is an item in the array

### object_from_lists

converts a pair of lists containing keys and values to an object

### random

Generates a random sequence of characters

### x509_decode

decodes an x.509 certificate to an object. you may also use this in conjunction with `base64_decode` jmespath function to decode a base64-encoded certificate

### time_to_cron

converts a time (RFC 3339) to a cron expression (string).

### time_add

adds duration (second string) to a time value (first string)

### time_parse

changes a time value of a given layout to RFC 3339

### time_utc

calcutes time in UTC from a given time in RFC 3339 format

### time_diff

calculate the difference between a start and end date in RFC3339 format

### time_before

checks if a time is before another time, both in RFC3339 format

### time_after

checks if a time is after another time, both in RFC3339 format

### time_between

checks if a time is between a start and end time, all in RFC3339 format

### time_truncate

returns the result of rounding time down to a multiple of duration

