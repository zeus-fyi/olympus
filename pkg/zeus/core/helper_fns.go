package zeus_core

func IsVersionNew(oldLabels, newLabels map[string]string) bool {
	currentVersion := GetVersionLabel(oldLabels)
	proposedVersion := GetVersionLabel(newLabels)

	if proposedVersion != currentVersion {
		return true
	}
	return false
}

func GetVersionLabel(labels map[string]string) string {
	if v, ok := labels["version"]; ok {
		return v
	}
	return "undefined"
}
