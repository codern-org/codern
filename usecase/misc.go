package usecase

import "github.com/codern-org/codern/domain"

type miscUsecase struct {
	miscRepository domain.MiscRepository
}

func NewMiscUsecase(miscRepository domain.MiscRepository) domain.MiscUsecase {
	return &miscUsecase{miscRepository: miscRepository}
}

func (u *miscUsecase) GetFeatureFlag(feature string) (bool, error) {
	return u.miscRepository.GetFeatureFlag(feature)
}
