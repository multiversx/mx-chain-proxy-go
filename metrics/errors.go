package metrics

import "errors"

var errNegativeCacheSize = errors.New("negative cache size")

var errNilMetricSaveFunction = errors.New("nil metrics save function")
