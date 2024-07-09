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

