package check

import (
	"sort"

	"github.com/hashicorp/go-multierror"
)

type Check struct {
	Options *CheckOptions
}

type CheckOptions struct {
	LegacyDataSourceFile *LegacyDataSourceFileOptions
	LegacyGuideFile      *LegacyGuideFileOptions
	LegacyIndexFile      *LegacyIndexFileOptions
	LegacyResourceFile   *LegacyResourceFileOptions

	RegistryDataSourceFile *RegistryDataSourceFileOptions
	RegistryGuideFile      *RegistryGuideFileOptions
	RegistryIndexFile      *RegistryIndexFileOptions
	RegistryResourceFile   *RegistryResourceFileOptions
}

func NewCheck(opts *CheckOptions) *Check {
	check := &Check{
		Options: opts,
	}

	if check.Options == nil {
		check.Options = &CheckOptions{}
	}

	return check
}

func (check *Check) Run(directories map[string][]string) error {
	if err := InvalidDirectoriesCheck(directories); err != nil {
		return err
	}

	if err := MixedDirectoriesCheck(directories); err != nil {
		return err
	}

	var result *multierror.Error

	if files, ok := directories[RegistryDataSourcesDirectory]; ok {
		if err := NewRegistryDataSourceFileCheck(check.Options.RegistryDataSourceFile).RunAll(files); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if files, ok := directories[RegistryGuidesDirectory]; ok {
		if err := NewRegistryGuideFileCheck(check.Options.RegistryGuideFile).RunAll(files); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if files, ok := directories[RegistryIndexDirectory]; ok {
		if err := NewRegistryIndexFileCheck(check.Options.RegistryIndexFile).RunAll(files); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if files, ok := directories[RegistryResourcesDirectory]; ok {
		if err := NewRegistryResourceFileCheck(check.Options.RegistryResourceFile).RunAll(files); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if files, ok := directories[LegacyDataSourcesDirectory]; ok {
		if err := NewLegacyDataSourceFileCheck(check.Options.LegacyDataSourceFile).RunAll(files); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if files, ok := directories[LegacyGuidesDirectory]; ok {
		if err := NewLegacyGuideFileCheck(check.Options.LegacyGuideFile).RunAll(files); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if files, ok := directories[LegacyIndexDirectory]; ok {
		if err := NewLegacyIndexFileCheck(check.Options.LegacyIndexFile).RunAll(files); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if files, ok := directories[LegacyResourcesDirectory]; ok {
		if err := NewLegacyResourceFileCheck(check.Options.LegacyResourceFile).RunAll(files); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if result != nil {
		sort.Sort(result)
	}

	return result.ErrorOrNil()
}
