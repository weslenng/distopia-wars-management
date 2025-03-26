package errors

type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

var PlayerNotFound = &ServiceError{
	Message: "Você não está registrado, digite `.help` para saber mais",
}

var OtherPlayerNotFound = &ServiceError{
	Message: "O jogador não foi encontrado",
}

var PlayerAlreadyRegistered = &ServiceError{
	Message: "Você já está registrado",
}

var NicknameAlreadyExists = &ServiceError{
	Message: "O nickname já está em uso",
}
