package types

type GetGameStateRequest struct {
	JoinCode JoinCode `json:"join_code"`
	UserIds  []Uuid   `json:"user_ids"`
}

type GetGameStateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Game    Game   `json:"game"`
	Tiles   []Tile `json:"tiles"`
	Turn    Turn   `json:"turn"`
}

type GetGameStateResponseMap = map[Uuid]GetGameStateResponse

func GetGameStateError(message string) GetGameStateResponse {
	return GetGameStateResponse{
		Success: false,
		Message: message,
	}
}

func GetGameStateErrorMap(playerIds []Uuid, message string) GetGameStateResponseMap {
	res := make(GetGameStateResponseMap)
	for _, id := range playerIds {
		res[id] = GetGameStateError(message)
	}
	return res
}
