package domain

type MiscRepository interface {
	GetFeatureFlag(feature string) (bool, error)
}

type MiscUsecase interface {
	GetFeatureFlag(feature string) (bool, error)
}
