package state

// FilterByProvider returns only the resources whose Provider field matches
// the given provider name (e.g. "aws", "google").
func FilterByProvider(resources []Resource, provider string) []Resource {
	if provider == "" {
		return resources
	}
	out := make([]Resource, 0, len(resources))
	for _, r := range resources {
		if r.Provider == provider {
			out = append(out, r)
		}
	}
	return out
}

// FilterByType returns only the resources whose Type field matches the given
// resource type string (e.g. "aws_instance").
func FilterByType(resources []Resource, resourceType string) []Resource {
	if resourceType == "" {
		return resources
	}
	out := make([]Resource, 0, len(resources))
	for _, r := range resources {
		if r.Type == resourceType {
			out = append(out, r)
		}
	}
	return out
}

// IndexByName returns a map keyed by resource Name for quick lookups.
func IndexByName(resources []Resource) map[string]Resource {
	idx := make(map[string]Resource, len(resources))
	for _, r := range resources {
		idx[r.Name] = r
	}
	return idx
}

// FilterDrifted returns only the resources that are marked as drifted.
// A resource is considered drifted when its Drifted field is true.
func FilterDrifted(resources []Resource) []Resource {
	out := make([]Resource, 0, len(resources))
	for _, r := range resources {
		if r.Drifted {
			out = append(out, r)
		}
	}
	return out
}
