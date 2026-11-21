package users

type Owner struct {
	Email string `json:"email" bidirekt:"field"`
}

type Account struct {
	Owner
	Login string `json:"login" bidirekt:"field"`
}
