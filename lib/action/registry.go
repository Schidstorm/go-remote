package action

type Registry struct {
	Children []RegistrationHandler
}

func (reg *Registry) HandleRegistration(registrator Registrator) error {
	for _, child := range reg.Children {
		err := child.HandleRegistration(registrator)

		if err != nil {
			return err
		}
	}

	return nil
}
