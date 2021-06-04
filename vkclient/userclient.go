package vkclient

type UserClient struct {
	baseClient *BaseClient
}

func NewUserClient(token string, v string) (*UserClient, error) {
	client, err := New(token, v)
	if err != nil {
		return nil, err
	}

	return &UserClient{
		baseClient: client,
	}, nil
}

func (userClient *UserClient) Search() {

}
