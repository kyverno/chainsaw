# Functions

!!! warning "Experimental functions"

    Experimental functions are denoted by the `x_` prefix.

    These are functions that are subject to signature change in a future version.

## built-in functions

| Name | Signature | Description |
|---|---|---|
| [abs](./examples/abs.md) | `abs(number)` | Returns the absolute value of the provided argument. |
| [avg](./examples/avg.md) | `avg(array[number])` | Returns the average of the elements in the provided array. An empty array will produce a return value of null. |
| [ceil](./examples/ceil.md) | `ceil(number)` | Returns the next highest integer value by rounding up if necessary. |
| [contains](./examples/contains.md) | `contains(array|string, any)` | Returns `true` if the given subject contains the provided search value. If the subject is an array, this function returns `true` if one of the elements in the array is equal to the provided search value. If the provided subject is a string, this function returns `true` if the string contains the provided search argument. |
| [ends_with](./examples/ends_with.md) | `ends_with(string, string)` | Reports whether the given string ends with the provided suffix argument. |
| [find_first](./examples/find_first.md) | `find_first(string, string, number, number)` |  |
| [find_last](./examples/find_last.md) | `find_last(string, string, number, number)` |  |
| [floor](./examples/floor.md) | `floor(number)` | Returns the next lowest integer value by rounding down if necessary. |
| [from_items](./examples/from_items.md) | `from_items(array[array])` |  |
| [group_by](./examples/group_by.md) | `group_by(array, expref)` |  |
| [items](./examples/items.md) | `items(object)` | Converts a given object into an array of key-value pairs. |
| [join](./examples/join.md) | `join(string, array[string])` | Returns all of the elements from the provided array joined together using the glue argument as a separator between each. |
| [keys](./examples/keys.md) | `keys(object)` | Returns an array containing the keys of the provided object. |
| [length](./examples/length.md) | `length(string|array|object)` | Returns the length of the given argument. If the argument is a string this function returns the number of code points in the string. If the argument is an array this function returns the number of elements in the array. If the argument is an object this function returns the number of key-value pairs in the object. |
| [lower](./examples/lower.md) | `lower(string)` | Returns the given string with all Unicode letters mapped to their lower case. |
| [map](./examples/map.md) | `map(expref, array)` | Transforms elements in a given array and returns the result. |
| [max](./examples/max.md) | `max(array[number]|array[string])` | Returns the highest found element in the provided array argument. An empty array will produce a return value of null. |
| [max_by](./examples/max_by.md) | `max_by(array, expref)` | Returns the highest found element using a custom expression to compute the associated value for each element in the input array. |
| [merge](./examples/merge.md) | `merge(object)` | Meges a list of objects together and returns the result. |
| [min](./examples/min.md) | `min(array[number]|array[string])` | Returns the lowest found element in the provided array argument. |
| [min_by](./examples/min_by.md) | `min_by(array, expref)` | Returns the lowest found element using a custom expression to compute the associated value for each element in the input array. |
| [not_null](./examples/not_null.md) | `not_null(any)` | Returns the first non null element in the input array. |
| [pad_left](./examples/pad_left.md) | `pad_left(string, number, string)` |  |
| [pad_right](./examples/pad_right.md) | `pad_right(string, number, string)` |  |
| [replace](./examples/replace.md) | `replace(string, string, string, number)` | Returns a copy of the input string with instances of old string argument replaced by new string argument. |
| [reverse](./examples/reverse.md) | `reverse(array|string)` | Reverses the input string or array and returns the result. |
| [sort](./examples/sort.md) | `sort(array[string]|array[number])` | This function accepts an array argument and returns the sorted elements as an array. |
| [sort_by](./examples/sort_by.md) | `sort_by(array, expref)` | This function accepts an array argument and returns the sorted elements as an array using a custom expression to compute the associated value for each element. |
| [split](./examples/split.md) | `split(string, string, number)` | Slices input string into substrings separated by a string argument and returns an array of the substrings between those separators. |
| [starts_with](./examples/starts_with.md) | `starts_with(string, string)` | Reports whether the input string begins with the provided string prefix argument. |
| [sum](./examples/sum.md) | `sum(array[number])` | Returns the sum of all numbers contained in the provided array. |
| [to_array](./examples/to_array.md) | `to_array(any)` |  |
| [to_number](./examples/to_number.md) | `to_number(any)` | Returns the parsed number. |
| [to_string](./examples/to_string.md) | `to_string(any)` | The JSON encoded value of the given argument. |
| [trim](./examples/trim.md) | `trim(string, string)` |  |
| [trim_left](./examples/trim_left.md) | `trim_left(string, string)` |  |
| [trim_right](./examples/trim_right.md) | `trim_right(string, string)` |  |
| [type](./examples/type.md) | `type(any)` | Returns the JavaScript type of the given argument as a string value. |
| [upper](./examples/upper.md) | `upper(string)` | Returns the given string with all Unicode letters mapped to their upper case. |
| [values](./examples/values.md) | `values(object)` | Returns the values of the provided object. |
| [zip](./examples/zip.md) | `zip(array, array)` |  |

## kyverno-json functions

| Name | Signature | Description |
|---|---|---|
| [at](./examples/at.md) | `at(array, number)` | Returns the element in an array at the given index. |
| [concat](./examples/concat.md) | `concat(string, string)` | Concatenates two strings together and returns the result. |
| [json_parse](./examples/json_parse.md) | `json_parse(string)` | Parses a given JSON string into an object. |
| [wildcard](./examples/wildcard.md) | `wildcard(string, string)` | Compares a wildcard pattern with a given string and returns if they match or not. |

## kyverno functions

| Name | Signature | Description |
|---|---|---|
| [compare](./examples/compare.md) | `compare(string, string)` | compares two strings lexicographically |
| [equal_fold](./examples/equal_fold.md) | `equal_fold(string, string)` | allows comparing two strings for equivalency where the only differences are letter cases |
| [replace](./examples/replace.md) | `replace(string, string, string, number)` | replaces a specified number of instances of the source string with the replacement string in a parent  |
| [replace_all](./examples/replace_all.md) | `replace_all(string, string, string)` | replace all instances of one string with another in an overall parent string |
| [to_upper](./examples/to_upper.md) | `to_upper(string)` | takes in a string and outputs the same string with all upper-case letters |
| [to_lower](./examples/to_lower.md) | `to_lower(string)` | takes in a string and outputs the same string with all lower-case letters |
| [trim](./examples/trim.md) | `trim(string, string)` | trims both ends of the source string by characters appearing in the second string |
| [trim_prefix](./examples/trim_prefix.md) | `trim_prefix(string, string)` | trims the second string prefix from the first string if the first string starts with the prefix |
| [split](./examples/split.md) | `split(string, string)` | splits the first string when the second string is found and converts it into an array  |
| [regex_replace_all](./examples/regex_replace_all.md) | `regex_replace_all(string, string|number, string|number)` | converts all parameters to string |
| [regex_replace_all_literal](./examples/regex_replace_all_literal.md) | `regex_replace_all_literal(string, string|number, string|number)` | converts all parameters to string |
| [regex_match](./examples/regex_match.md) | `regex_match(string, string|number)` | first string is the regular exression which is compared with second input which can be a number or string |
| [pattern_match](./examples/pattern_match.md) | `pattern_match(string, string|number)` | '*' matches zero or more alphanumeric characters, '?' matches a single alphanumeric character |
| [label_match](./examples/label_match.md) | `label_match(object, object)` | object arguments must be enclosed in backticks; ex. `{{request.object.spec.template.metadata.labels}}` |
| [to_boolean](./examples/to_boolean.md) | `to_boolean(string)` | It returns true or false for any string, such as 'True', 'TruE', 'False', 'FAlse', 'faLSE', etc. |
| [add](./examples/add.md) | `add(any, any)` | does arithmetic addition of two specified values of numbers, quantities, and durations |
| [sum](./examples/sum.md) | `sum(array)` | does arithmetic addition of specified array of values of numbers, quantities, and durations |
| [subtract](./examples/subtract.md) | `subtract(any, any)` | does arithmetic subtraction of two specified values of numbers, quantities, and durations |
| [multiply](./examples/multiply.md) | `multiply(any, any)` | does arithmetic multiplication of two specified values of numbers, quantities, and durations |
| [divide](./examples/divide.md) | `divide(any, any)` | divisor must be non zero |
| [modulo](./examples/modulo.md) | `modulo(any, any)` | divisor must be non-zero, arguments must be integers |
| [round](./examples/round.md) | `round(number, number)` | does roundoff to upto the given decimal places |
| [base64_decode](./examples/base64_decode.md) | `base64_decode(string)` | decodes a base 64 string |
| [base64_encode](./examples/base64_encode.md) | `base64_encode(string)` | encodes a regular, plaintext and unencoded string to base64 |
| [time_since](./examples/time_since.md) | `time_since(string, string, string)` | calculate the difference between a start and end period of time where the end may either be a static definition or the then-current time |
| [time_now](./examples/time_now.md) | `time_now()` | returns current time in RFC 3339 format |
| [time_now_utc](./examples/time_now_utc.md) | `time_now_utc()` | returns current UTC time in RFC 3339 format |
| [path_canonicalize](./examples/path_canonicalize.md) | `path_canonicalize(string)` | normalizes or canonicalizes a given path by removing excess slashes |
| [truncate](./examples/truncate.md) | `truncate(string, number)` | length argument must be enclosed in backticks; ex. "{{request.object.metadata.name | truncate(@, `9`)}}" |
| [semver_compare](./examples/semver_compare.md) | `semver_compare(string, string)` | compares two strings which comply with the semantic versioning schema and outputs a boolean response as to the position of the second relative to the first |
| [parse_json](./examples/parse_json.md) | `parse_json(string)` | decodes a valid JSON encoded string to the appropriate type. Opposite of `to_string` function |
| [parse_yaml](./examples/parse_yaml.md) | `parse_yaml(string)` | decodes a valid YAML encoded string to the appropriate type provided it can be represented as JSON |
| [lookup](./examples/lookup.md) | `lookup(object|array, string|number)` | returns the value corresponding to the given key/index in the given object/array |
| [items](./examples/items.md) | `items(object|array, string, string)` | converts a map or array to an array of objects where each key:value is an item in the array |
| [object_from_lists](./examples/object_from_lists.md) | `object_from_lists(array, array)` | converts a pair of lists containing keys and values to an object |
| [random](./examples/random.md) | `random(string)` | Generates a random sequence of characters |
| [x509_decode](./examples/x509_decode.md) | `x509_decode(string)` | decodes an x.509 certificate to an object. you may also use this in conjunction with `base64_decode` jmespath function to decode a base64-encoded certificate |
| [time_to_cron](./examples/time_to_cron.md) | `time_to_cron(string)` | converts a time (RFC 3339) to a cron expression (string). |
| [time_add](./examples/time_add.md) | `time_add(string, string)` | adds duration (second string) to a time value (first string) |
| [time_parse](./examples/time_parse.md) | `time_parse(string, string)` | changes a time value of a given layout to RFC 3339 |
| [time_utc](./examples/time_utc.md) | `time_utc(string)` | calcutes time in UTC from a given time in RFC 3339 format |
| [time_diff](./examples/time_diff.md) | `time_diff(string, string)` | calculate the difference between a start and end date in RFC3339 format |
| [time_before](./examples/time_before.md) | `time_before(string, string)` | checks if a time is before another time, both in RFC3339 format |
| [time_after](./examples/time_after.md) | `time_after(string, string)` | checks if a time is after another time, both in RFC3339 format |
| [time_between](./examples/time_between.md) | `time_between(string, string, string)` | checks if a time is between a start and end time, all in RFC3339 format |
| [time_truncate](./examples/time_truncate.md) | `time_truncate(string, string)` | returns the result of rounding time down to a multiple of duration |

## chainsaw functions

| Name | Signature | Description |
|---|---|---|
| [env](./examples/env.md) | `env(string)` | Returns the value of the environment variable passed in argument. |
| [x_k8s_get](./examples/x_k8s_get.md) | `x_k8s_get(any, string, string, string, string)` | Gets a resource from a Kubernetes cluster. |
| [x_k8s_list](./examples/x_k8s_list.md) | `x_k8s_list(any, string, string, string)` | Lists resources from a Kubernetes cluster. |
| [x_k8s_exists](./examples/x_k8s_exists.md) | `x_k8s_exists(any, string, string, string, string)` | Checks if a given resource exists in a Kubernetes cluster. |
| [x_k8s_resource_exists](./examples/x_k8s_resource_exists.md) | `x_k8s_resource_exists(any, string, string)` | Checks if a given resource type is available in a Kubernetes cluster. |
| [x_k8s_server_version](./examples/x_k8s_server_version.md) | `x_k8s_server_version(any)` | Returns the version of a Kubernetes cluster. |
| [x_metrics_decode](./examples/x_metrics_decode.md) | `x_metrics_decode(string)` | Decodes metrics in the Prometheus text format. |
| [trim_space](./examples/trim_space.md) | `trim_space(string)` | Trims leading and trailing spaces from the string passed in argument. |
| [as_string](./examples/as_string.md) | `as_string(any)` | Returns the passed in argument converted into a string. |

