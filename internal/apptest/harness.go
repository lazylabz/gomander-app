// Package apptest drives the Gomander backend the way the desktop app drives
// it: real use cases, real repositories, a real event bus and a real database,
// with fakes only where the backend leaves the process - spawning processes,
// the desktop runtime, and the filesystem.
//
// A test written against it arranges data, runs a backend operation through the
// use case registry, and asserts on what the user or the operating system would
// see. No repository fake takes part.
package apptest

import (
	"context"
	"testing"

	"github.com/google/uuid"

	commandhandlers "gomander/internal/command/application/handlers"
	commandusecases "gomander/internal/command/application/usecases"
	commanddomain "gomander/internal/command/domain"
	commandinfrastructure "gomander/internal/command/infrastructure"
	commandgrouphandlers "gomander/internal/commandgroup/application/handlers"
	commandgroupusecases "gomander/internal/commandgroup/application/usecases"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	commandgroupinfrastructure "gomander/internal/commandgroup/infrastructure"
	configusecases "gomander/internal/config/application/usecases"
	configdomain "gomander/internal/config/domain"
	configinfrastructure "gomander/internal/config/infrastructure"
	"gomander/internal/event"
	"gomander/internal/eventbus"
	"gomander/internal/logger"
	"gomander/internal/openedproject"
	projectusecases "gomander/internal/project/application/usecases"
	projectdomain "gomander/internal/project/domain"
	projectinfrastructure "gomander/internal/project/infrastructure"
	"gomander/internal/testdb"
	"gomander/internal/usecases"

	internalapp "gomander/internal/app"
)

// Harness is the seam every backend behaviour is verified through. UseCases is
// the same registry the Wails controllers and the third-party HTTP server hold.
type Harness struct {
	t        *testing.T
	UseCases usecases.Registry

	// The repositories are here to arrange data, never to assert on it: what a
	// test asserts is what the app itself can see.
	commandRepository      commanddomain.Repository
	commandGroupRepository commandgroupdomain.Repository
	projectRepository      projectdomain.Repository
	configRepository       configdomain.Repository
	openedProject          openedproject.OpenedProject

	app           *internalapp.App
	processRunner *processRunnerFake
	events        *eventSinkFake
	dialogs       *dialogsFake
	fs            *fsFacadeFake
}

// New starts a backend with an empty database and no project open.
//
// The graph below mirrors buildDeps in cmd/gomander/main.go, which is where the
// app wires the same pieces together. The localization use cases are the one
// omission: they read the locale files embedded in the desktop binary, and no
// backend behaviour depends on them.
func New(t *testing.T) *Harness {
	t.Helper()

	ctx := context.Background()
	db := testdb.New(t)

	events := &eventSinkFake{}
	dialogs := &dialogsFake{}
	fsFacade := &fsFacadeFake{files: make(map[string][]byte)}
	processRunner := &processRunnerFake{
		startFailures: make(map[string]error),
		stopFailures:  make(map[string]error),
	}

	l := logger.NewDefaultLogger(ctx, &logSinkFake{})
	ee := event.NewDefaultEventEmitter(ctx, events)

	commandRepo := commandinfrastructure.NewGormCommandRepository(db, ctx)
	commandGroupRepo := commandgroupinfrastructure.NewGormCommandGroupRepository(db, ctx)
	projectRepo := projectinfrastructure.NewGormProjectRepository(db, ctx)
	configRepo := configinfrastructure.NewGormConfigRepository(db, ctx)

	eventBus := eventbus.NewInMemoryEventBus()

	openedProject := openedproject.NewOpenedProject(configRepo, projectRepo)

	app := internalapp.NewApp(l, processRunner, configRepo, eventBus, internalapp.EventHandlers{
		CleanCommandGroupsOnCommandDeleted:   commandgrouphandlers.NewCleanCommandGroupsOnCommandDeleted(commandGroupRepo, ee),
		CleanCommandGroupsOnProjectDeleted:   commandgrouphandlers.NewCleanCommandGroupsOnProjectDeleted(commandGroupRepo, ee),
		CleanCommandsOnProjectDeleted:        commandhandlers.NewCleanCommandOnProjectDeleted(commandRepo),
		AddCommandToGroupOnCommandDuplicated: commandgrouphandlers.NewAddCommandToGroupOnCommandDuplicated(commandRepo, commandGroupRepo),
	})
	app.RegisterHandlers()
	if err := app.Startup(ctx); err != nil {
		t.Fatalf("failed to start the app: %v", err)
	}

	return &Harness{
		t: t,
		UseCases: usecases.Registry{
			GetUserConfig:  configusecases.NewGetUserConfig(configRepo),
			SaveUserConfig: configusecases.NewSaveUserConfig(configRepo),

			GetCurrentProject:    projectusecases.NewGetCurrentProject(openedProject),
			GetAvailableProjects: projectusecases.NewGetAvailableProjects(projectRepo),
			OpenProject:          projectusecases.NewOpenProject(openedProject),
			CreateProject:        projectusecases.NewCreateProject(projectRepo),
			EditProject:          projectusecases.NewEditProject(projectRepo),
			CloseProject:         projectusecases.NewCloseProject(openedProject),
			DeleteProject:        projectusecases.NewDeleteProject(projectRepo, eventBus, l),
			ExportProject:        projectusecases.NewExportProject(projectRepo, commandRepo, commandGroupRepo, dialogs, fsFacade),
			ImportProject:        projectusecases.NewImportProject(projectRepo, commandRepo, commandGroupRepo),
			GetProjectToImport:   projectusecases.NewGetProjectToImport(dialogs, fsFacade),

			GetCommandGroups:              commandgroupusecases.NewGetCommandGroups(openedProject, commandGroupRepo),
			CreateCommandGroup:            commandgroupusecases.NewCreateCommandGroup(openedProject, commandGroupRepo),
			UpdateCommandGroup:            commandgroupusecases.NewUpdateCommandGroup(commandGroupRepo),
			DeleteCommandGroup:            commandgroupusecases.NewDeleteCommandGroup(commandGroupRepo, ee),
			RemoveCommandFromCommandGroup: commandgroupusecases.NewRemoveCommandFromCommandGroup(commandGroupRepo),
			ReorderCommandGroups:          commandgroupusecases.NewReorderCommandGroups(openedProject, commandGroupRepo),
			RunCommandGroup:               commandgroupusecases.NewRunCommandGroup(openedProject, commandGroupRepo, processRunner),
			StopCommandGroup:              commandgroupusecases.NewStopCommandGroup(commandGroupRepo, processRunner),

			GetCommands:          commandusecases.NewGetCommands(openedProject, commandRepo),
			AddCommand:           commandusecases.NewAddCommand(openedProject, commandRepo),
			DuplicateCommand:     commandusecases.NewDuplicateCommand(openedProject, commandRepo, eventBus),
			RemoveCommand:        commandusecases.NewRemoveCommand(commandRepo, eventBus),
			EditCommand:          commandusecases.NewEditCommand(commandRepo),
			ReorderCommands:      commandusecases.NewReorderCommands(openedProject, commandRepo),
			RunCommand:           commandusecases.NewRunCommand(openedProject, commandRepo, processRunner),
			StopCommand:          commandusecases.NewStopCommand(commandRepo, processRunner),
			GetRunningCommandIds: commandusecases.NewGetRunningCommandIds(processRunner),
		},
		commandRepository:      commandRepo,
		commandGroupRepository: commandGroupRepo,
		projectRepository:      projectRepo,
		configRepository:       configRepo,
		openedProject:          openedProject,
		app:                    app,
		processRunner:          processRunner,
		events:                 events,
		dialogs:                dialogs,
		fs:                     fsFacade,
	}
}

func (h *Harness) GivenProjects(projects ...projectdomain.Project) {
	h.t.Helper()

	for _, project := range projects {
		if err := h.projectRepository.Create(project); err != nil {
			h.t.Fatalf("failed to arrange project %s: %v", project.Id, err)
		}
	}
}

func (h *Harness) GivenOpenedProject(projectId string) {
	h.t.Helper()

	if err := h.openedProject.Open(projectId); err != nil {
		h.t.Fatalf("failed to arrange the opened project: %v", err)
	}
}

func (h *Harness) GivenCommands(commands ...commanddomain.Command) {
	h.t.Helper()

	for i := range commands {
		if err := h.commandRepository.Create(&commands[i]); err != nil {
			h.t.Fatalf("failed to arrange command %s: %v", commands[i].Id, err)
		}
	}
}

func (h *Harness) GivenCommandGroups(commandGroups ...commandgroupdomain.CommandGroup) {
	h.t.Helper()

	for i := range commandGroups {
		if err := h.commandGroupRepository.Create(&commandGroups[i]); err != nil {
			h.t.Fatalf("failed to arrange command group %s: %v", commandGroups[i].Id, err)
		}
	}
}

func (h *Harness) GivenEnvironmentPaths(paths ...string) {
	h.t.Helper()

	config := h.currentConfig()
	config.EnvironmentPaths = make([]configdomain.EnvironmentPath, 0, len(paths))
	for _, path := range paths {
		config.EnvironmentPaths = append(config.EnvironmentPaths, configdomain.EnvironmentPath{
			Id:   uuid.New().String(),
			Path: path,
		})
	}

	if err := h.configRepository.Update(config); err != nil {
		h.t.Fatalf("failed to arrange the environment paths: %v", err)
	}
}

// GivenProcessThatCannotStart makes spawning the Command fail, the way the
// operating system refuses one whose working directory is gone.
func (h *Harness) GivenProcessThatCannotStart(commandId string, err error) {
	h.processRunner.failStart(commandId, err)
}

// GivenProcessThatCannotStop makes signalling the Command's process fail.
func (h *Harness) GivenProcessThatCannotStop(commandId string, err error) {
	h.processRunner.failStop(commandId, err)
}

// GivenFileToImport puts a file where the open-file dialog will point the user.
func (h *Harness) GivenFileToImport(path string, contents []byte) {
	h.t.Helper()

	h.fs.put(path, contents)
	h.dialogs.openFilePath = path
}

// GivenMissingFileToImport points the open-file dialog at a path the
// filesystem has no file at.
func (h *Harness) GivenMissingFileToImport(path string) {
	h.t.Helper()

	h.dialogs.openFilePath = path
}

// GivenExportDestination answers the save-file dialog with path.
func (h *Harness) GivenExportDestination(path string) {
	h.t.Helper()

	h.dialogs.saveFilePath = path
}

// GivenDialogsThatFail makes the desktop toolkit refuse to put any dialog on
// screen.
func (h *Harness) GivenDialogsThatFail(err error) {
	h.t.Helper()

	h.dialogs.failure = err
}

// GivenADestinationThatCannotBeWritten makes the operating system refuse every
// write, the way a full disk or a read-only folder does.
func (h *Harness) GivenADestinationThatCannotBeWritten(err error) {
	h.t.Helper()

	h.fs.writeFailure = err
}

// ClosingTheApp runs the shutdown the desktop window triggers, and answers
// whether the app refused to close.
func (h *Harness) ClosingTheApp() (prevent bool) {
	h.t.Helper()

	return h.app.OnBeforeClose(context.Background())
}

func (h *Harness) StartedProcesses() []StartedProcess {
	return h.processRunner.started()
}

func (h *Harness) StoppedProcessIds() []string {
	return h.processRunner.stopped()
}

func (h *Harness) EmittedEvents() []EmittedEvent {
	return h.events.emitted()
}

func (h *Harness) ExportedFile(path string) ([]byte, bool) {
	return h.fs.file(path)
}

func (h *Harness) currentConfig() *configdomain.Config {
	h.t.Helper()

	config, err := h.configRepository.GetOrCreate()
	if err != nil {
		h.t.Fatalf("failed to read the user config: %v", err)
	}

	return config
}
