package stage

import (
	"github.com/flant/dapp/pkg/build/builder"
	"github.com/flant/dapp/pkg/config"
	"github.com/flant/dapp/pkg/image"
	"github.com/flant/dapp/pkg/util"
)

func GenerateInstallStage(dimgConfig config.DimgInterface, extra *builder.Extra, baseStageOptions *NewBaseStageOptions) *InstallStage {
	b := getBuilder(dimgConfig, extra)
	if b != nil && !b.IsInstallEmpty() {
		return newInstallStage(b, baseStageOptions)
	}

	return nil
}

func newInstallStage(builder builder.Builder, baseStageOptions *NewBaseStageOptions) *InstallStage {
	s := &InstallStage{}
	s.UserWithGAPatchStage = newUserWithGAPatchStage(builder, baseStageOptions)
	return s
}

type InstallStage struct {
	*UserWithGAPatchStage
}

func (s *InstallStage) Name() StageName {
	return Install
}

func (s *InstallStage) GetDependencies(_ Conveyor, _ image.Image) (string, error) {
	stageDependenciesChecksum, err := s.GetStageDependenciesChecksum(Install)
	if err != nil {
		return "", err
	}

	return util.Sha256Hash(s.builder.InstallChecksum(), stageDependenciesChecksum), nil
}

func (s *InstallStage) PrepareImage(c Conveyor, prevBuiltImage, image image.Image) error {
	if err := s.UserWithGAPatchStage.PrepareImage(c, prevBuiltImage, image); err != nil {
		return nil
	}

	if err := s.builder.Install(image.BuilderContainer()); err != nil {
		return err
	}

	return nil
}
