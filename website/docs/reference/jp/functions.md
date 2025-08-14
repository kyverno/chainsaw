# Functions

!!! warning "Experimental functions"

    Experimental functions are denoted by the `x_` prefix.

    These are functions that are subject to signature change in a future version.

## built-in functions

| Name | Description |
|---|---|
| [abs](./examples/abs.md) | Returns the absolute value of the provided argument. |
| [avg](./examples/avg.md) | Returns the average of the elements in the provided array. An empty array will produce a return value of null. |
| [ceil](./examples/ceil.md) | Returns the next highest integer value by rounding up if necessary. |
| [contains](./examples/contains.md) | Returns `true` if the given subject contains the provided search value. If the subject is an array, this function returns `true` if one of the elements in the array is equal to the provided search value. If the provided subject is a string, this function returns `true` if the string contains the provided search argument. |
| [ends_with](./examples/ends_with.md) | Reports whether the given string ends with the provided suffix argument. |
| [find_first](./examples/find_first.md) | Returns the zero-based index of the first occurence where the substring appears in a string or null if it does not appear. |
| [find_last](./examples/find_last.md) | Returns the zero-based index of the last occurence where the substring appears in a string or null if it does not appear. |
| [floor](./examples/floor.md) | Returns the next lowest integer value by rounding down if necessary. |
| [from_items](./examples/from_items.md) | Returns an object from the provided array of key value pairs. This function is the inversed of the `items()` function. |
| [group_by](./examples/group_by.md) | Groups an array of objects using an expression as the group key. |
| [items](./examples/items.md) | Converts a given object into an array of key-value pairs. |
| [join](./examples/join.md) | Returns all of the elements from the provided array joined together using the glue argument as a separator between each. |
| [keys](./examples/keys.md) | Returns an array containing the keys of the provided object. |
| [length](./examples/length.md) | Returns the length of the given argument. If the argument is a string this function returns the number of code points in the string. If the argument is an array this function returns the number of elements in the array. If the argument is an object this function returns the number of key-value pairs in the object. |
| [lower](./examples/lower.md) | Returns the given string with all Unicode letters mapped to their lower case. |
| [map](./examples/map.md) | Transforms elements in a given array and returns the result. |
| [max](./examples/max.md) | Returns the highest found element in the provided array argument. An empty array will produce a return value of null. |
| [max_by](./examples/max_by.md) | Returns the highest found element using a custom expression to compute the associated value for each element in the input array. |
| [merge](./examples/merge.md) | Meges a list of objects together and returns the result. |
| [min](./examples/min.md) | Returns the lowest found element in the provided array argument. |
| [min_by](./examples/min_by.md) | Returns the lowest found element using a custom expression to compute the associated value for each element in the input array. |
| [not_null](./examples/not_null.md) | Returns the first non null element in the input array. |
| [pad_left](./examples/pad_left.md) | Adds characters to the beginning of a string. |
| [pad_right](./examples/pad_right.md) | Adds characters to the end of a string. |
| [replace](./examples/replace.md) | Returns a copy of the input string with instances of old string argument replaced by new string argument. |
| [reverse](./examples/reverse.md) | Reverses the input string or array and returns the result. |
| [sort](./examples/sort.md) | This function accepts an array argument and returns the sorted elements as an array. |
| [sort_by](./examples/sort_by.md) | This function accepts an array argument and returns the sorted elements as an array using a custom expression to compute the associated value for each element. |
| [split](./examples/split.md) | Slices input string into substrings separated by a string argument and returns an array of the substrings between those separators. |
| [starts_with](./examples/starts_with.md) | Reports whether the input string begins with the provided string prefix argument. |
| [sum](./examples/sum.md) | Returns the sum of all numbers contained in the provided array. |
| [to_array](./examples/to_array.md) | Returns a one element array containing the passed in argument, or the passed in value if it's an array. |
| [to_number](./examples/to_number.md) | Returns the parsed number. |
| [to_string](./examples/to_string.md) | The JSON encoded value of the given argument. |
| [trim](./examples/trim.md) | Removes the leading and trailing characters found in the passed in string argument. |
| [trim_left](./examples/trim_left.md) | Removes the leading characters found in the passed in string argument. |
| [trim_right](./examples/trim_right.md) | Removes the trailing characters found in the passed in string argument. |
| [type](./examples/type.md) | Returns the JavaScript type of the given argument as a string value. |
| [upper](./examples/upper.md) | Returns the given string with all Unicode letters mapped to their upper case. |
| [values](./examples/values.md) | Returns the values of the provided object. |
| [zip](./examples/zip.md) | Accepts one or more arrays as arguments and returns an array of arrays in which the i-th array contains the i-th element from each of the argument arrays. The returned array is truncated to the length of the shortest argument array. |

## kyverno-json functions

| Name | Description |
|---|---|
| [at](./examples/at.md) | Returns the element in an array at the given index. |
| [concat](./examples/concat.md) | Concatenates two strings together and returns the result. |
| [json_parse](./examples/json_parse.md) | Parses a given JSON string into an object. |
| [wildcard](./examples/wildcard.md) | Compares a wildcard pattern with a given string and returns if they match or not. |

## kyverno functions

| Name | Description |
|---|---|
| [compare](./examples/compare.md) | Compares two strings lexicographically. |
| [equal_fold](./examples/equal_fold.md) | Allows comparing two strings for equivalency where the only differences are letter cases. |
| [replace](./examples/replace.md) | Replaces a specified number of instances of the source string with the replacement string in a parent. |
| [replace_all](./examples/replace_all.md) | Replace all instances of one string with another in an overall parent string. |
| [to_upper](./examples/to_upper.md) | Takes in a string and outputs the same string with all upper-case letters. |
| [to_lower](./examples/to_lower.md) | Takes in a string and outputs the same string with all lower-case letters. |
| [trim](./examples/trim.md) | Trims both ends of the source string by characters appearing in the second string. |
| [trim_prefix](./examples/trim_prefix.md) | Trims the second string prefix from the first string if the first string starts with the prefix. |
| [split](./examples/split.md) | Splits the first string when the second string is found and converts it into an array. |
| [regex_replace_all](./examples/regex_replace_all.md) | Converts all parameters to string. |
| [regex_replace_all_literal](./examples/regex_replace_all_literal.md) | Converts all parameters to string. |
| [regex_match](./examples/regex_match.md) | First string is the regular exression which is compared with second input which can be a number or string. |
| [pattern_match](./examples/pattern_match.md) | '*' matches zero or more alphanumeric characters, '?' matches a single alphanumeric character. |
| [label_match](./examples/label_match.md) | Object arguments must be enclosed in backticks; ex. `{{request.object.spec.template.metadata.labels}}`. |
| [to_boolean](./examples/to_boolean.md) | It returns true or false for any string, such as 'True', 'TruE', 'False', 'FAlse', 'faLSE', etc. |
| [add](./examples/add.md) | Does arithmetic addition of two specified values of numbers, quantities, and durations. |
| [sum](./examples/sum.md) | Does arithmetic addition of specified array of values of numbers, quantities, and durations. |
| [subtract](./examples/subtract.md) | Does arithmetic subtraction of two specified values of numbers, quantities, and durations. |
| [multiply](./examples/multiply.md) | Does arithmetic multiplication of two specified values of numbers, quantities, and durations. |
| [divide](./examples/divide.md) | Divisor must be non zero. |
| [modulo](./examples/modulo.md) | Divisor must be non-zero, arguments must be integers. |
| [round](./examples/round.md) | Does roundoff to upto the given decimal places. |
| [base64_decode](./examples/base64_decode.md) | Decodes a base 64 string. |
| [base64_encode](./examples/base64_encode.md) | Encodes a regular, plaintext and unencoded string to base64. |
| [time_since](./examples/time_since.md) | Calculate the difference between a start and end period of time where the end may either be a static definition or the then-current time. |
| [time_now](./examples/time_now.md) | Returns current time in RFC 3339 format. |
| [time_now_utc](./examples/time_now_utc.md) | Returns current UTC time in RFC 3339 format. |
| [path_canonicalize](./examples/path_canonicalize.md) | Normalizes or canonicalizes a given path by removing excess slashes. |
| [truncate](./examples/truncate.md) | Length argument must be enclosed in backticks; ex. "{{request.object.metadata.name | truncate(@, `9`)}}". |
| [semver_compare](./examples/semver_compare.md) | Compares two strings which comply with the semantic versioning schema and outputs a boolean response as to the position of the second relative to the first. |
| [parse_json](./examples/parse_json.md) | Decodes a valid JSON encoded string to the appropriate type. Opposite of `to_string` function. |
| [parse_yaml](./examples/parse_yaml.md) | Decodes a valid YAML encoded string to the appropriate type provided it can be represented as JSON. |
| [lookup](./examples/lookup.md) | Returns the value corresponding to the given key/index in the given object/array. |
| [items](./examples/items.md) | Converts a map or array to an array of objects where each key:value is an item in the array. |
| [object_from_lists](./examples/object_from_lists.md) | Converts a pair of lists containing keys and values to an object. |
| [random](./examples/random.md) | Generates a random sequence of characters. |
| [x509_decode](./examples/x509_decode.md) | Decodes an x.509 certificate to an object. you may also use this in conjunction with `base64_decode` jmespath function to decode a base64-encoded certificate. |
| [time_to_cron](./examples/time_to_cron.md) | Converts a time (RFC 3339) to a cron expression (string). |
| [time_add](./examples/time_add.md) | Adds duration (second string) to a time value (first string). |
| [time_parse](./examples/time_parse.md) | Changes a time value of a given layout to RFC 3339. |
| [time_utc](./examples/time_utc.md) | Calcutes time in UTC from a given time in RFC 3339 format. |
| [time_diff](./examples/time_diff.md) | Calculate the difference between a start and end date in RFC3339 format. |
| [time_before](./examples/time_before.md) | Checks if a time is before another time, both in RFC3339 format. |
| [time_after](./examples/time_after.md) | Checks if a time is after another time, both in RFC3339 format. |
| [time_between](./examples/time_between.md) | Checks if a time is between a start and end time, all in RFC3339 format. |
| [time_truncate](./examples/time_truncate.md) | Returns the result of rounding time down to a multiple of duration. |

## chainsaw functions

| Name | Description |
|---|---|
| [env](./examples/env.md) | Returns the value of the environment variable passed in argument. |
| [x_k8s_get](./examples/x_k8s_get.md) | Gets a resource from a Kubernetes cluster. |
| [x_k8s_list](./examples/x_k8s_list.md) | Lists resources from a Kubernetes cluster. |
| [x_k8s_exists](./examples/x_k8s_exists.md) | Checks if a given resource exists in a Kubernetes cluster. |
| [x_k8s_resource_exists](./examples/x_k8s_resource_exists.md) | Checks if a given resource type is available in a Kubernetes cluster. |
| [x_k8s_server_version](./examples/x_k8s_server_version.md) | Returns the version of a Kubernetes cluster. |
| [x_metrics_decode](./examples/x_metrics_decode.md) | Decodes metrics in the Prometheus text format. |
| [trim_space](./examples/trim_space.md) | Trims leading and trailing spaces from the string passed in argument. |
| [as_string](./examples/as_string.md) | Returns the passed in argument converted into a string. |

