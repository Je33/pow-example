package pow

func (s *Service) Validate(challenge, nonce string) error {
	return s.vld.Validate(challenge, nonce)
}
