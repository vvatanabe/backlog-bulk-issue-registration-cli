package internal

import (
	"github.com/vvatanabe/shot/shot"
	"github.com/pkg/errors"
)

func NewModule(cfg *Config) (Module, error) {
	injector, err := shot.CreateInjector(func(binder shot.Binder) {
		binder.Bind(new(Config)).ToInstance(cfg)
		if cfg.Lang == En {
			binder.Bind(new(Messages)).ToConstructor(NewEnglish).AsEagerSingleton()
		} else {
			binder.Bind(new(Messages)).ToConstructor(NewJapanese).AsEagerSingleton()
		}
		binder.Bind(new(BacklogAPIClient)).ToConstructor(NewBacklogAPIClient).AsEagerSingleton()
		binder.Bind(new(IssueRepository)).ToConstructor(NewIssueHTTPClient).AsEagerSingleton()
		binder.Bind(new(ProjectRepository)).ToConstructor(NewProjectHTTPClient).AsEagerSingleton()
		binder.Bind(new(CommandConverter)).ToConstructor(NewCommandConverter).AsEagerSingleton()
		binder.Bind(new(CommandBuilder)).ToConstructor(NewCommandBuilder).AsEagerSingleton()
		binder.Bind(new(BulkCommandExecutor)).ToConstructor(NewBulkCommandExecutor).AsEagerSingleton()
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not resolve internal modules.")
	}

	return &module{injector}, nil
}

type Module interface {
	GetMessages() Messages
	GetCommandConverter() *CommandConverter
	GetBulkCommandExecutor() *BulkCommandExecutor
	GetProjectRepository() ProjectRepository
}

type module struct {
	shot.Injector
}

func (m *module) GetMessages() Messages {
	return m.Get(new(Messages)).(Messages)
}

func (m *module) GetCommandConverter() *CommandConverter {
	return m.Get(new(CommandConverter)).(*CommandConverter)
}

func (m *module) GetBulkCommandExecutor() *BulkCommandExecutor {
	return m.Get(new(BulkCommandExecutor)).(*BulkCommandExecutor)
}

func (m *module) GetProjectRepository() ProjectRepository {
	return m.Get(new(ProjectRepository)).(ProjectRepository)
}
