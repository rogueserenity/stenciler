package config

// Merge merges the local template with the repository template. It uses the contents of the repository template
// and fills in the values of parameters from the local template. It also sets the repository URL to the value from
// the local template.
func Merge(repoTemplate, localTemplate *Template) *Template {
	merged := Template{}
	merged.Repository = localTemplate.Repository
	merged.Directory = repoTemplate.Directory
	merged.Params = mergeParams(repoTemplate.Params, localTemplate.Params)
	merged.InitOnlyPaths = repoTemplate.InitOnlyPaths
	merged.RawCopyPaths = repoTemplate.RawCopyPaths
	merged.PreInitHookPaths = repoTemplate.PreInitHookPaths
	merged.PostInitHookPaths = repoTemplate.PostInitHookPaths
	merged.PreUpdateHookPaths = repoTemplate.PreUpdateHookPaths
	merged.PostUpdateHookPaths = repoTemplate.PostUpdateHookPaths

	return &merged
}

func mergeParams(repoParams, localParams []*Param) []*Param {
	params := make([]*Param, 0, len(repoParams))
	localValues := localParamValues(localParams)
	for _, p := range repoParams {
		param := &Param{
			Name:           p.Name,
			Prompt:         p.Prompt,
			Default:        p.Default,
			ValidationHook: p.ValidationHook,
			Value:          p.Value,
		}
		if value, ok := localValues[p.Name]; ok {
			param.Value = value
		}
		params = append(params, param)
	}
	if len(params) == 0 {
		return nil
	}
	return params
}

func localParamValues(params []*Param) map[string]string {
	values := make(map[string]string, len(params))
	for _, p := range params {
		if len(p.Prompt) > 0 {
			// only add if it was prompted for
			values[p.Name] = p.Value
		}
	}
	return values
}
