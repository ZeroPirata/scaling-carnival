package helper

func CopiarLabels(original map[string]string) map[string]string {
	if original == nil {
		return nil
	}
	copia := make(map[string]string, len(original))
	for k, v := range original {
		copia[k] = v
	}
	return copia
}
