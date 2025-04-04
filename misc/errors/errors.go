package errors

type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

var NicknameAlreadyExists = &ServiceError{
	Message: "O nickname já está em uso",
}

var PlayerAlreadyJoined = &ServiceError{
	Message: "Você já está em um time",
}

var PlayerAlreadyRegistered = &ServiceError{
	Message: "Você já está registrado",
}

var PlayerNotFound = &ServiceError{
	Message: "Você não está registrado, digite `.help` para saber mais",
}

var PlayerCannotChangeTeam = &ServiceError{
	Message: "Você não tem permissão para mudar de time",
}
