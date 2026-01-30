package data

func handleValid() (int, error) {
	return 1, nil
}

// FAIL: not start with handle
func notValid() (int, error) { // want `private function 'notValid' have more than one value return must start with 'handle'`
	return 1, nil
}

func setup() {}

// FAIL: not start with is
func set() bool { // want `private function 'set' have bool value must start with 'is'`
	return false
}

func foo() int {
	return 1
}
