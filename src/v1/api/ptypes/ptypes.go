package ptypes

func NewEmpty(err error) (*Empty, error) {
	if err != nil {
		return nil, err
	}
	return &Empty{}, nil
}
