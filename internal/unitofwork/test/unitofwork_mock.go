package test

import (
	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/unitofwork"
)

type MockUnitOfWork struct {
	repositories unitofwork.Repositories
}

func NewMockUnitOfWork(
	projectRepository projectdomain.Repository,
	commandRepository commanddomain.Repository,
	commandGroupRepository commandgroupdomain.Repository,
) *MockUnitOfWork {
	return &MockUnitOfWork{
		repositories: unitofwork.Repositories{
			Projects:      projectRepository,
			Commands:      commandRepository,
			CommandGroups: commandGroupRepository,
		},
	}
}

// Do runs change against the repositories it was given: the transaction is the
// real Unit of Work's business, and a test that stubs one is testing GORM.
func (m *MockUnitOfWork) Do(change func(unitofwork.Repositories) error) error {
	return change(m.repositories)
}
