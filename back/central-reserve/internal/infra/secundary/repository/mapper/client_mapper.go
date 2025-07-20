package mapper

import (
	"central_reserve/internal/domain/entities"
)

func ClientToClient(client entities.Client) entities.Client {
	return client
}

func ClientSliceToClientSlice(clients []entities.Client) []entities.Client {
	return clients
}
