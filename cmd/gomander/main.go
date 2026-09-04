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
	"gomander/internal/dialog"
	"gomander/internal/eventbus"
	"gomander/internal/facade"
	localizationusecases "gomander/internal/localization/application/usecases"
	localizationinfrastructure "gomander/internal/localization/infrastructure"
	"gomander/internal/logger"
	"gomander/internal/openedproject"
	projectusecases "gomander/internal/project/application/usecases"
	projectinfrastructure "gomander/internal/project/infrastructure"
	releaseusecases "gomander/internal/releases/application/usecases"
	releaseinfrastructure "gomander/internal/releases/infrastructure"
	"gomander/internal/runner"
	"gomander/internal/uihelpers/fs"
	"gomander/internal/uihelpers/os_internal"
	"gomander/internal/uihelpers/path"
	unitofworkinfrastructure "gomander/internal/unitofwork/infrastructure"
	"gomander/internal/usecases"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	internalapp "gomander/internal/app"
	"gomander/internal/event"
	_ "gomander/migrations"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed locales
var localeFs embed.FS

const ConfigFolderPathName = "gomander"

func main() {
	// The app can only be built once Wails hands us its context, so OnStartup fills this in
	var app *internalapp.App

	// Create instance of helpers
	uiPathHelper := path.NewUiPathHelper()
	dialogs := dialog.NewWailsDialogs()
	uiFsHelper := fs.NewUIFsHelper(dialogs, facade.DefaultRuntimeFacade{})
	uiOsHelper := os_internal.NewUIOsHelper()

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

			// Load context into the desktop adapters
			dialog.SetWailsDialogsContext(dialogs, ctx)

			// Build deps
			var useCases usecases.Registry
			app, useCases = buildDeps(gormDb, ctx, dialogs)

			// Register event handlers
			app.RegisterHandlers()

			// Start app
			if err := app.Startup(ctx); err != nil {
				panic(err)
			}

			// Initialize controllers
			controllers.loadUseCases(useCases)

			// Start http server for 3rd party integrations
			server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)

			go func() {
				err := server.RegisterHandlers()
				if err != nil {
					panic(err)
				}
				server.Start()
			}()
		},
		Bind: []interface{}{
			uiPathHelper,
			controllers,
			uiFsHelper,
			uiOsHelper,
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			return app.OnBeforeClose(ctx)
		},
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

func buildDeps(gormDb *gorm.DB, ctx context.Context, dialogs dialog.Dialogs) (*internalapp.App, usecases.Registry) {
	// Initialize deps
	l := logger.NewDefaultLogger(ctx, facade.DefaultRuntimeFacade{})
	ee := event.NewDefaultEventEmitter(ctx, facade.DefaultRuntimeFacade{})
	r := runner.NewDefaultRunner(l, ee)
	releaseFeed := releaseinfrastructure.NewGithubReleaseFeed(ctx, facade.DefaultIOFacade{}, releaseinfrastructure.DefaultLatestReleaseUrl)
	releaseDownloader := releaseinfrastructure.NewGithubReleaseDownloader(ctx, facade.DefaultOSFacade{}, facade.DefaultIOFacade{}, releaseinfrastructure.DefaultBinaryDownloadBaseUrl)
	releaseInstaller := releaseinfrastructure.NewOSReleaseInstaller(facade.DefaultOSFacade{}, facade.DefaultOpenFacade{})
	appControl := releaseinfrastructure.NewShellAppControl(ctx, facade.DefaultRuntimeFacade{})

	// Initialize repos
	commandRepo := commmandinfrastructure.NewGormCommandRepository(gormDb, ctx)
	commandGroupRepo := commandgroupinfrastructure.NewGormCommandGroupRepository(gormDb, ctx)
	projectRepo := projectinfrastructure.NewGormProjectRepository(gormDb, ctx)
	configRepo := configinfrastructure.NewGormConfigRepository(gormDb, ctx)

	// Initialize the unit of work
	unitOfWork := unitofworkinfrastructure.NewGormUnitOfWork(gormDb, ctx)

	// Initialize event handlers
	cleanCommandGroupsOnCommandDeleted := commandgrouphandlers.NewCleanCommandGroupsOnCommandDeleted(unitOfWork, commandGroupRepo, ee)
	cleanCommandGroupsOnProjectDeleted := commandgrouphandlers.NewCleanCommandGroupsOnProjectDeleted(unitOfWork, ee)
	cleanCommandsOnProjectDeleted := handlers.NewCleanCommandOnProjectDeleted(commandRepo)
	addCommandToGroupOnCommandDuplicated := commandgrouphandlers.NewAddCommandToGroupOnCommandDuplicated(commandRepo, commandGroupRepo)

	// Initialize event bus
	eventBus := eventbus.NewInMemoryEventBus()

	// Initialize the opened project
	openedProject := openedproject.NewOpenedProject(configRepo, projectRepo)

	// Initialize use cases

	// Configuration
	getUserConfig := configusecases.NewGetUserConfig(configRepo)
	saveUserConfig := configusecases.NewSaveUserConfig(configRepo)
	// Localization
	localeFiles := localizationinfrastructure.NewLocaleFiles(localeFs)
	getTranslation := localizationusecases.NewGetTranslation(localeFiles)
	getSupportedLanguages := localizationusecases.NewGetSupportedLanguages(localeFiles)
	// Projects
	getCurrentProject := projectusecases.NewGetCurrentProject(openedProject)
	getAvailableProjects := projectusecases.NewGetAvailableProjects(projectRepo)
	openProject := projectusecases.NewOpenProject(openedProject)
	createProject := projectusecases.NewCreateProject(projectRepo)
	editProject := projectusecases.NewEditProject(projectRepo)
	closeProject := projectusecases.NewCloseProject(openedProject)
	deleteProject := projectusecases.NewDeleteProject(projectRepo, eventBus)
	projectFile := projectinfrastructure.NewProjectFile(facade.DefaultFsFacade{})
	exportProject := projectusecases.NewExportProject(projectRepo, commandRepo, commandGroupRepo, dialogs, projectFile)
	importProject := projectusecases.NewImportProject(unitOfWork)
	getProjectToImport := projectusecases.NewGetProjectToImport(dialogs, map[projectusecases.FileType]projectusecases.BlueprintReader{
		projectusecases.FileTypeGomander:    projectFile,
		projectusecases.FileTypePackageJSON: projectinfrastructure.NewPackageJSONFile(facade.DefaultFsFacade{}),
	})
	// Command Groups
	getCommandGroups := commandgroupusecases.NewGetCommandGroups(openedProject, commandGroupRepo)
	createCommandGroup := commandgroupusecases.NewCreateCommandGroup(openedProject, commandGroupRepo)
	updateCommandGroup := commandgroupusecases.NewUpdateCommandGroup(commandGroupRepo)
	deleteCommandGroup := commandgroupusecases.NewDeleteCommandGroup(commandGroupRepo, ee)
	removeCommandFromCommandGroup := commandgroupusecases.NewRemoveCommandFromCommandGroup(commandGroupRepo)
	reorderCommandGroups := commandgroupusecases.NewReorderCommandGroups(openedProject, commandGroupRepo)
	runCommandGroup := commandgroupusecases.NewRunCommandGroup(openedProject, commandGroupRepo, commandRepo, r)
	stopCommandGroup := commandgroupusecases.NewStopCommandGroup(commandGroupRepo, r)
	// Commands
	getCommands := commandusecases.NewGetCommands(openedProject, commandRepo)
	addCommand := commandusecases.NewAddCommand(openedProject, commandRepo)
	duplicateCommand := commandusecases.NewDuplicateCommand(openedProject, commandRepo, eventBus)
	removeCommand := commandusecases.NewRemoveCommand(commandRepo, eventBus)
	editCommand := commandusecases.NewEditCommand(commandRepo)
	reorderCommands := commandusecases.NewReorderCommands(openedProject, commandRepo)
	runCommand := commandusecases.NewRunCommand(openedProject, commandRepo, r)
	stopCommand := commandusecases.NewStopCommand(commandRepo, r)
	getRunningCommandIds := commandusecases.NewGetRunningCommandIds(r)
	// Releases
	getCurrentRelease := releaseusecases.NewGetCurrentRelease()
	checkForNewRelease := releaseusecases.NewCheckForNewRelease(releaseFeed)
	downloadRelease := releaseusecases.NewDownloadRelease(releaseDownloader)
	installReleaseAndQuit := releaseusecases.NewInstallReleaseAndQuit(releaseInstaller, appControl)

	app := internalapp.NewApp(l, r, configRepo, eventBus, internalapp.EventHandlers{
		CleanCommandGroupsOnCommandDeleted:   cleanCommandGroupsOnCommandDeleted,
		CleanCommandGroupsOnProjectDeleted:   cleanCommandGroupsOnProjectDeleted,
		CleanCommandsOnProjectDeleted:        cleanCommandsOnProjectDeleted,
		AddCommandToGroupOnCommandDuplicated: addCommandToGroupOnCommandDuplicated,
	})

	return app, usecases.Registry{
		// Configuration
		GetUserConfig:  getUserConfig,
		SaveUserConfig: saveUserConfig,
		// Localization
		GetTranslation:        getTranslation,
		GetSupportedLanguages: getSupportedLanguages,
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
		// Releases
		GetCurrentRelease:     getCurrentRelease,
		CheckForNewRelease:    checkForNewRelease,
		DownloadRelease:       downloadRelease,
		InstallReleaseAndQuit: installReleaseAndQuit,
	}
}
