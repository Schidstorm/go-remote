package action

type RegistrationHandler interface {
	HandleRegistration(registrator Registrator) error
}
