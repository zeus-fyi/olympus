package dev_hacks

// Use lets you avoid the annoying var not used message that prevents compiling
func Use(vals ...interface{}) error {
	for _, val := range vals {
		_ = val
	}
	return nil
}
