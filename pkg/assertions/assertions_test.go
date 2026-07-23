package assertions

import "testing"

// TestCountResourcesOfType exercises the resource-counting logic that backs
// AssertResourceCount, including the cases the previous hardcoded regex
// (`module\.example\.<type>\.`) silently mis-counted as 0: count/for_each
// indexed module instances and resources nested in submodules.
func TestCountResourcesOfType(t *testing.T) {
	cases := []struct {
		name         string
		stateList    string
		resourceType string
		want         int
	}{
		{
			name:         "direct child of module.example",
			stateList:    "module.example.local_file_bundle.this",
			resourceType: "local_file_bundle",
			want:         1,
		},
		{
			name:         "for_each indexed resource instances",
			stateList:    "module.example.local_file_bundle.this[\"http\"]\nmodule.example.local_file_bundle.this[\"https\"]",
			resourceType: "local_file_bundle",
			want:         2,
		},
		{
			name:         "count/for_each indexed module instance",
			stateList:    "module.example[\"a\"].local_file_bundle.this\nmodule.example[\"b\"].local_file_bundle.this",
			resourceType: "local_file_bundle",
			want:         2,
		},
		{
			name:         "nested submodule",
			stateList:    "module.example.module.inner.local_file_bundle.this",
			resourceType: "local_file_bundle",
			want:         1,
		},
		{
			name:         "type is matched as a full segment (local_file must not match local_file_bundle)",
			stateList:    "module.example.local_file_bundle.this\nmodule.example.local_file.this",
			resourceType: "local_file",
			want:         1,
		},
		{
			name:         "root-level fixtures outside module.example are not counted",
			stateList:    "local_file.fixture\nlocal_file_bundle.fixture\nmodule.example.local_file_bundle.this",
			resourceType: "local_file_bundle",
			want:         1,
		},
		{
			name:         "same-named data source under module.example is not counted as a resource",
			stateList:    "module.example.data.random_pet.selected\nmodule.example.random_pet.this",
			resourceType: "random_pet",
			want:         1,
		},
		{
			name:         "no matching resources",
			stateList:    "module.example.null_resource.this",
			resourceType: "local_file_bundle",
			want:         0,
		},
		{
			name:         "empty state list",
			stateList:    "",
			resourceType: "local_file_bundle",
			want:         0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := countResourcesOfType(tc.stateList, tc.resourceType)
			if got != tc.want {
				t.Errorf("countResourcesOfType(resourceType=%q) = %d, want %d", tc.resourceType, got, tc.want)
			}
		})
	}
}
