package helpers

func RetrieveHostAvailableMemory() (int64, error) {
	// This function should return the available memory on the host.
	// For now, we will return a dummy value.
	// In a real implementation, you would use a library or system call to get the actual memory.
	return 1024, nil // Returning 1 GB as an example
}

func RetrieveHostAvailableVCpu() (float64, error) {
	// This function should return the available vCPU on the host.
	// For now, we will return a dummy value.
	// In a real implementation, you would use a library or system call to get the actual vCPU.
	return 2, nil // Returning 2 vCPUs as an example
}
