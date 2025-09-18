package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"

	"gomander/cmd/gomander/thirdpartyserver"
	"gomander/internal/command/application/handlers"
	commandusecases "gomander/internal/command/application/usecases"
	commmandinfrastructure "gomander/internal/command/infrastructure"
	commandgrouphandlers "gomander/internal/commandgroup/application/handlers"
	commandgroupusecases "gomander/internal/commandgroup/application/usecases"
	commandgroupinfrastructure "gomander/internal/commandgroup/infrastructure"
	configusecases "gomander/internal/config/application/usecases"
	configinfrastructure "gomander/internal/config/infrastructure"
	"gomander/internal/eventbus"
	"gomander/internal/facade"
	"gomander/internal/logger"
	projectusecases "gomander/internal/project/application/usecases"
	projectinfrastructure "gomander/internal/project/infrastructure"
	"gomander/internal/releases"
	"gomander/internal/runner"
	"gomander/internal/uihelpers/fs"
	"gomander/internal/uihelpers/path"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	internalapp "gomander/internal/app"
	"gomander/internal/event"
	_ "gomander/migrations"
)

//go:embed all:frontend/dist
var assets embed.FS

const ConfigFolderPathName = "gomander"

func main() {
	// Create an instance of the app structure
	app := internalapp.NewApp()

	// Create instance of helpers
	uiPathHelper := path.NewUiPathHelper()
	uiFsHelper := fs.NewUIFsHelper(facade.DefaultRuntimeFacade{})
	releaseHelper := releases.NewReleaseHelper()

	// Create instance of controllers
	controllers := NewWailsControllers()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Gomander",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			// Initialize the database
			gormDb := configDB(ctx)

			// Register deps
			registerDeps(gormDb, ctx, app)

			// Register event handlers
			app.RegisterHandlers()

			// Start app
			app.Startup(ctx)

			// Initialize controllers
			controllers.loadUseCases(app.UseCases)

			// Load context into helpers
			uiFsHelper.SetContext(ctx)

			// Start http server for 3rd party integrations
			server := thirdpartyserver.NewThirdPartyIntegrationsServer(app.UseCases)

			go func() {
				err := server.RegisterHandlers()
				if err != nil {
					panic(err)
				}
				server.Start()
			}()
		},
		Bind: []interface{}{
			app,
			uiPathHelper,
			controllers,
			releaseHelper,
			uiFsHelper,
		},
		OnBeforeClose: app.OnBeforeClose,
		EnumBind: []interface{}{
			event.Events,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func configDB(ctx context.Context) *gorm.DB {
	gormDb, err := gorm.Open(sqlite.Open(getDbFile()+"?cache=shared"), &gorm.Config{
		// Uncomment when debugging
		// Logger: gormlogger.Default.LogMode(gormlogger.Info),
		Logger: gormlogger.Default.LogMode(gormlogger.Error),
	})
	if err != nil {
		panic(err)
	}

	db, err := gormDb.DB()

	if err != nil {
		panic(err)
	}

	if db == nil {
		panic("db is nil")
	}

	db.SetMaxOpenConns(1)

	// Execute migrations
	err = goose.SetDialect("sqlite3")
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embed.FS{})

	err = goose.UpContext(ctx, db, ".")
	if err != nil {
		panic(err)
	}
	return gormDb
}

func getDbFile() string {
	userConfig, err := os.UserConfigDir()

	if err != nil {
		panic(err)
	}

	configFolderPath := filepath.Join(userConfig, ConfigFolderPathName)
	err = os.MkdirAll(configFolderPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	dbLocation := filepath.Join(configFolderPath, "data.db")

	return dbLocation
}

func registerDeps(gormDb *gorm.DB, ctx context.Context, app *internalapp.App) {
	// Initialize deps
	l := logger.NewDefaultLogger(ctx, facade.DefaultRuntimeFacade{})
	ee := event.NewDefaultEventEmitter(ctx, facade.DefaultRuntimeFacade{})
	r := runner.NewDefaultRunner(l, ee)

	// Initialize repos
	commandRepo := commmandinfrastructure.NewGormCommandRepository(gormDb, ctx)
	commandGroupRepo := commandgroupinfrastructure.NewGormCommandGroupRepository(gormDb, ctx)
	projectRepo := projectinfrastructure.NewGormProjectRepository(gormDb, ctx)
	configRepo := configinfrastructure.NewGormConfigRepository(gormDb, ctx)

	// Initialize event handlers
	cleanCommandGroupsOnCommandDeleted := commandgrouphandlers.NewCleanCommandGroupsOnCommandDeleted(commandGroupRepo, ee)
	cleanCommandGroupsOnProjectDeleted := commandgrouphandlers.NewCleanCommandGroupsOnProjectDeleted(commandGroupRepo, ee)
	cleanCommandsOnProjectDeleted := handlers.NewCleanCommandOnProjectDeleted(commandRepo)
	addCommandToGroupOnCommandDuplicated := commandgrouphandlers.NewAddCommandToGroupOnCommandDuplicated(commandRepo, commandGroupRepo)

	// Initialize event bus
	eventBus := eventbus.NewInMemoryEventBus()

	// Initialize use cases

	// Configuration
	getUserConfig := configusecases.NewGetUserConfig(configRepo)
	saveUserConfig := configusecases.NewSaveUserConfig(configRepo)
	// Projects
	getCurrentProject := projectusecases.NewGetCurrentProject(configRepo, projectRepo)
	getAvailableProjects := projectusecases.NewGetAvailableProjects(projectRepo)
	openProject := projectusecases.NewOpenProject(configRepo, projectRepo)
	createProject := projectusecases.NewCreateProject(projectRepo)
	editProject := projectusecases.NewEditProject(projectRepo)
	closeProject := projectusecases.NewCloseProject(configRepo)
	deleteProject := projectusecases.NewDeleteProject(projectRepo, eventBus, l)
	exportProject := projectusecases.NewExportProject(ctx, projectRepo, commandRepo, commandGroupRepo, facade.DefaultRuntimeFacade{}, facade.DefaultFsFacade{})
	importProject := projectusecases.NewImportProject(projectRepo, commandRepo, commandGroupRepo)
	getProjectToImport := projectusecases.NewGetProjectToImport(ctx, facade.DefaultRuntimeFacade{}, facade.DefaultFsFacade{})
	// Command Groups
	getCommandGroups := commandgroupusecases.NewGetCommandGroups(configRepo, commandGroupRepo)
	createCommandGroup := commandgroupusecases.NewCreateCommandGroup(configRepo, commandGroupRepo)
	updateCommandGroup := commandgroupusecases.NewUpdateCommandGroup(commandGroupRepo)
	deleteCommandGroup := commandgroupusecases.NewDeleteCommandGroup(commandGroupRepo, ee)
	removeCommandFromCommandGroup := commandgroupusecases.NewRemoveCommandFromCommandGroup(commandGroupRepo)
	reorderCommandGroups := commandgroupusecases.NewReorderCommandGroups(configRepo, commandGroupRepo)
	runCommandGroup := commandgroupusecases.NewRunCommandGroup(configRepo, commandRepo, commandGroupRepo, projectRepo, r)
	stopCommandGroup := commandgroupusecases.NewStopCommandGroup(commandGroupRepo, r)
	// Commands
	getCommands := commandusecases.NewGetCommands(configRepo, commandRepo)
	addCommand := commandusecases.NewAddCommand(configRepo, commandRepo)
	duplicateCommand := commandusecases.NewDuplicateCommand(configRepo, commandRepo, eventBus)
	removeCommand := commandusecases.NewRemoveCommand(commandRepo, eventBus)
	editCommand := commandusecases.NewEditCommand(commandRepo)
	reorderCommands := commandusecases.NewReorderCommands(configRepo, commandRepo)
	runCommand := commandusecases.NewRunCommand(configRepo, commandRepo, projectRepo, r)
	stopCommand := commandusecases.NewStopCommand(commandRepo, r)
	getRunningCommandIds := commandusecases.NewGetRunningCommandIds(r)

	app.LoadDependencies(internalapp.Dependencies{
		Logger:       l,
		EventEmitter: ee,
		Runner:       r,

		CommandRepository:      commandRepo,
		CommandGroupRepository: commandGroupRepo,
		ProjectRepository:      projectRepo,
		ConfigRepository:       configRepo,

		FsFacade:      facade.DefaultFsFacade{},
		RuntimeFacade: facade.DefaultRuntimeFacade{},

		EventBus: eventBus,
		EventHandlers: internalapp.EventHandlers{
			CleanCommandGroupsOnCommandDeleted:   cleanCommandGroupsOnCommandDeleted,
			CleanCommandGroupsOnProjectDeleted:   cleanCommandGroupsOnProjectDeleted,
			CleanCommandsOnProjectDeleted:        cleanCommandsOnProjectDeleted,
			AddCommandToGroupOnCommandDuplicated: addCommandToGroupOnCommandDuplicated,
		},

		UseCases: internalapp.UseCases{
			// Configuration
			GetUserConfig:  getUserConfig,
			SaveUserConfig: saveUserConfig,
			// Projects
			GetCurrentProject:    getCurrentProject,
			GetAvailableProjects: getAvailableProjects,
			OpenProject:          openProject,
			CreateProject:        createProject,
			EditProject:          editProject,
			CloseProject:         closeProject,
			DeleteProject:        deleteProject,
			ExportProject:        exportProject,
			ImportProject:        importProject,
			GetProjectToImport:   getProjectToImport,
			// Command Groups
			GetCommandGroups:              getCommandGroups,
			CreateCommandGroup:            createCommandGroup,
			UpdateCommandGroup:            updateCommandGroup,
			DeleteCommandGroup:            deleteCommandGroup,
			RemoveCommandFromCommandGroup: removeCommandFromCommandGroup,
			ReorderCommandGroups:          reorderCommandGroups,
			RunCommandGroup:               runCommandGroup,
			StopCommandGroup:              stopCommandGroup,
			// Commands
			GetCommands:          getCommands,
			AddCommand:           addCommand,
			DuplicateCommand:     duplicateCommand,
			RemoveCommand:        removeCommand,
			EditCommand:          editCommand,
			ReorderCommands:      reorderCommands,
			RunCommand:           runCommand,
			StopCommand:          stopCommand,
			GetRunningCommandIds: getRunningCommandIds,
		},
	})
}
