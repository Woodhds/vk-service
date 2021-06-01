package vkclient

type WallClient struct {
	baseclient *BaseClient
}

func NewWallClient(token string, version string) (*WallClient, error) {
	baseClient, e := New(token, version)

	if e != nil {
		return nil, e
	}

	return &WallClient{
		baseclient: baseClient,
	}, nil
}

func (wallClient *WallClient) Get() error {

}
