package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/application/usecases"
	commanddomain "gomander/internal/command/domain"
	commanddomainevent "gomander/internal/command/domain/event"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/testutils"
)

func TestDefaultDuplicateCommand_Execute(t *testing.T) {
	t.Run("Should duplicate the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)
		mockEventBus := new(MockEventBus)

		projectId := "project1"
		sut := usecases.NewDuplicateCommand(mockUserConfigRepository, mockCommandRepository, mockEventBus)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		originalCmdData := testutils.NewCommand().WithProjectId(projectId).Data()
		originalCmd := commandDataToDomain(originalCmdData)

		mockCommandRepository.On("Get", originalCmd.Id).Return(&originalCmd, nil)
		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{originalCmd}, nil)

		mockCommandRepository.On("Create", mock.MatchedBy(func(cmd *commanddomain.Command) bool {
			return cmd.Id != originalCmd.Id &&
				cmd.ProjectId == originalCmd.ProjectId &&
				cmd.Name == originalCmd.Name+" (copy)" &&
				cmd.Command == originalCmd.Command &&
				cmd.WorkingDirectory == originalCmd.WorkingDirectory &&
				cmd.Position == 1
		})).Return(nil)

		mockEventBus.On("PublishSync", mock.MatchedBy(func(event interface{}) bool {
			e, ok := event.(commanddomainevent.CommandDuplicatedEvent)
			return ok && e.InsideGroupId == ""
		})).Return(make([]error, 0))

		// Act
		err := sut.Execute(originalCmd.Id, "")

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockEventBus,
		)
	})

	t.Run("Should duplicate the command with a target group", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)
		mockEventBus := new(MockEventBus)

		projectId := "project1"
		targetGroupId := "group1"

		sut := usecases.NewDuplicateCommand(mockUserConfigRepository, mockCommandRepository, mockEventBus)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		originalCmdData := testutils.NewCommand().WithProjectId(projectId).Data()
		originalCmd := commandDataToDomain(originalCmdData)

		mockCommandRepository.On("Get", originalCmd.Id).Return(&originalCmd, nil)
		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{originalCmd}, nil)

		mockCommandRepository.On("Create", mock.MatchedBy(func(cmd *commanddomain.Command) bool {
			return cmd.Id != originalCmd.Id &&
				cmd.ProjectId == originalCmd.ProjectId &&
				cmd.Name == originalCmd.Name+" (copy)" &&
				cmd.Command == originalCmd.Command &&
				cmd.WorkingDirectory == originalCmd.WorkingDirectory &&
				cmd.Position == 1
		})).Return(nil)

		mockEventBus.On("PublishSync", mock.MatchedBy(func(event interface{}) bool {
			e, ok := event.(commanddomainevent.CommandDuplicatedEvent)
			return ok && e.InsideGroupId == targetGroupId
		})).Return(make([]error, 0))

		// Act
		err := sut.Execute(originalCmd.Id, targetGroupId)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockEventBus,
		)
	})

	t.Run("Should return an error if fails to get the user config", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)
		mockEventBus := new(MockEventBus)

		sut := usecases.NewDuplicateCommand(mockUserConfigRepository, mockCommandRepository, mockEventBus)
		commandId := "command1"
		expectedErr := errors.New("failed to get user config")
		mockUserConfigRepository.On("GetOrCreate").Return(nil, expectedErr)

		// Act
		err := sut.Execute(commandId, "")

		// Assert
		assert.ErrorIs(t, err, expectedErr)
		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if the command does not exist", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)
		mockEventBus := new(MockEventBus)

		sut := usecases.NewDuplicateCommand(mockUserConfigRepository, mockCommandRepository, mockEventBus)

		commandId := "non-existing-command"

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: "project1"}, nil)
		expectedErr := errors.New("command not found")
		mockCommandRepository.On("Get", commandId).Return(nil, expectedErr)

		// Act
		err := sut.Execute(commandId, "")

		// Assert
		assert.ErrorIs(t, err, expectedErr)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})

	t.Run("Should return an error if fails to get all commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)
		mockEventBus := new(MockEventBus)

		projectId := "project1"
		sut := usecases.NewDuplicateCommand(mockUserConfigRepository, mockCommandRepository, mockEventBus)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		originalCmd := commandDataToDomain(testutils.NewCommand().Data())
		mockCommandRepository.On("Get", originalCmd.Id).Return(&originalCmd, nil)

		expectedErr := errors.New("failed to get commands")
		mockCommandRepository.On("GetAll", projectId).Return(nil, expectedErr)

		// Act
		err := sut.Execute(originalCmd.Id, "")

		// Assert
		assert.ErrorIs(t, err, expectedErr)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})

	t.Run("Should return an error if fails to create the duplicated command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)
		mockEventBus := new(MockEventBus)

		projectId := "project1"
		sut := usecases.NewDuplicateCommand(mockUserConfigRepository, mockCommandRepository, mockEventBus)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		originalCmd := commandDataToDomain(testutils.NewCommand().Data())
		mockCommandRepository.On("Get", originalCmd.Id).Return(&originalCmd, nil)
		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{originalCmd}, nil)

		expectedErr := errors.New("failed to create command")
		mockCommandRepository.On("Create", mock.Anything).Return(expectedErr)

		// Act
		err := sut.Execute(originalCmd.Id, "")

		// Assert
		assert.ErrorIs(t, err, expectedErr)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})

	t.Run("Should return an error if side effects fail", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockConfigRepository)
		mockEventBus := new(MockEventBus)

		projectId := "project1"
		sut := usecases.NewDuplicateCommand(mockUserConfigRepository, mockCommandRepository, mockEventBus)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		originalCmdData := testutils.NewCommand().WithProjectId(projectId).Data()
		originalCmd := commandDataToDomain(originalCmdData)

		mockCommandRepository.On("Get", originalCmd.Id).Return(&originalCmd, nil)
		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{originalCmd}, nil)

		mockCommandRepository.On("Create", mock.Anything).Return(nil)

		sideEffectErrorMsg := "Something happened"
		mockEventBus.On("PublishSync", mock.MatchedBy(func(event interface{}) bool {
			e, ok := event.(commanddomainevent.CommandDuplicatedEvent)
			return ok && e.InsideGroupId == ""
		})).Return([]error{errors.New(sideEffectErrorMsg)})

		// Act
		err := sut.Execute(originalCmd.Id, "")

		// Assert
		assert.ErrorContains(t, err, sideEffectErrorMsg)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockEventBus,
		)
	})
}
