package errorx

import "bytes"

type (
	// A BatchError is an error that can hold multiple errors.
	BatchError struct {
		errs arrayErrors
	}

	arrayErrors []error
)

// Add adds err to be.
func (be *BatchError) Add(err error) {
	if err != nil {
		be.errs = append(be.errs, err)
	}
}

// Err returns an error that represents all errors.
func (be *BatchError) Err() error {
	switch len(be.errs) {
	case 0:
		return nil
	case 1:
		return be.errs[0]
	default:
		return be.errs
	}
}

// NotNil checks if any error inside.
func (be *BatchError) NotNil() bool {
	return len(be.errs) > 0
}

// Error returns a string that represents inside errors.
func (ea arrayErrors) Error() string {
	var buf bytes.Buffer

	for i := range ea {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(ea[i].Error())
	}

	return buf.String()
}
