package hooks

func ExecuteValidationHook(hook, name, value string) (string, error) {
	return value, nil
}

func ExecuteHooks(hooks []string) error {
	for _, hook := range hooks {
		if err := executeHook(hook); err != nil {
			return err
		}
	}
	return nil
}

func executeHook(hook string) error {
	return nil
}
