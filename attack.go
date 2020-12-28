package stormbringer

type Attack interface {
	Setup(c Config) error
	Do()
	Teardown() error
	Clone()
}
