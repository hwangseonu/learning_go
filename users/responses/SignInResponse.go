package responses

type SignInResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}