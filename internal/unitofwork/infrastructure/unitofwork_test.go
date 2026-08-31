package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	commanddomain "gomander/internal/command/domain"
	commandtest "gomander/internal/command/domain/test"
	commandinfrastructure "gomander/internal/command/infrastructure"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	commandgrouptest "gomander/internal/commandgroup/domain/test"
	commandgroupinfrastructure "gomander/internal/commandgroup/infrastructure"
	projectdomain "gomander/internal/project/domain"
	projecttest "gomander/internal/project/domain/test"
	projectinfrastructure "gomander/internal/project/infrastructure"
	"gomander/internal/testdb"
	"gomander/internal/unitofwork"
	"gomander/internal/unitofwork/infrastructure"
)

func TestGormUnitOfWork_Do(t *testing.T) {
	t.Run("Should keep the writes made inside it", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		db := testdb.New(t)
		sut := infrastructure.NewGormUnitOfWork(db, ctx)

		project := projecttest.NewProjectBuilder().Build()
		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		commandGroup := commandgrouptest.NewCommandGroupBuilder().WithProjectId(project.Id).Build()

		// Act
		err := sut.Do(func(repositories unitofwork.Repositories) error {
			if err := repositories.Projects.Create(project); err != nil {
				return err
			}
			if err := repositories.Commands.Create(&command); err != nil {
				return err
			}
			return repositories.CommandGroups.CreateWithCommandIds(&commandGroup)
		})

		// Assert
		assert.NoError(t, err)

		storedProjects, err := projectinfrastructure.NewGormProjectRepository(db, ctx).GetAll()
		assert.NoError(t, err)
		assert.Equal(t, []projectdomain.Project{project}, storedProjects)

		storedCommands, err := commandinfrastructure.NewGormCommandRepository(db, ctx).GetAll(project.Id)
		assert.NoError(t, err)
		assert.Equal(t, []commanddomain.Command{command}, storedCommands)

		storedCommandGroups, err := commandgroupinfrastructure.NewGormCommandGroupRepository(db, ctx).GetAllWithCommandIds(project.Id)
		assert.NoError(t, err)
		assert.Equal(t, []commandgroupdomain.CommandGroupWithCommandIds{commandGroup}, storedCommandGroups)
	})

	t.Run("Should undo every write it made when one of them fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		db := testdb.New(t)
		sut := infrastructure.NewGormUnitOfWork(db, ctx)

		project := projecttest.NewProjectBuilder().Build()
		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		commandGroup := commandgrouptest.NewCommandGroupBuilder().WithProjectId(project.Id).Build()

		// Act
		err := sut.Do(func(repositories unitofwork.Repositories) error {
			if err := repositories.Projects.Create(project); err != nil {
				return err
			}
			if err := repositories.Commands.Create(&command); err != nil {
				return err
			}
			if err := repositories.CommandGroups.CreateWithCommandIds(&commandGroup); err != nil {
				return err
			}
			return assert.AnError
		})

		// Assert
		assert.ErrorIs(t, err, assert.AnError)

		storedProjects, err := projectinfrastructure.NewGormProjectRepository(db, ctx).GetAll()
		assert.NoError(t, err)
		assert.Empty(t, storedProjects)

		storedCommands, err := commandinfrastructure.NewGormCommandRepository(db, ctx).GetAll(project.Id)
		assert.NoError(t, err)
		assert.Empty(t, storedCommands)

		storedCommandGroups, err := commandgroupinfrastructure.NewGormCommandGroupRepository(db, ctx).GetAllWithCommandIds(project.Id)
		assert.NoError(t, err)
		assert.Empty(t, storedCommandGroups)
	})
}
