package commands

var errorBucket []error

func appendError(err error) {
	if err != nil {
		errorBucket = append(errorBucket, err)
	}
}

func HasErrors() bool {
	return len(errorBucket) > 0
}
