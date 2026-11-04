package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"
	"runtime"

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

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"

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
	// ShouldQuit has to be handed to application.New, before the app it asks exists
	var app *application.App
	var gomander *internalapp.App
	var useCases usecases.Registry
	zoom := zoomOperationsOf(&useCases)

	app = application.New(application.Options{
		Name:   "Gomander",
		Logger: productionLogger(),
		ShouldQuit: func() bool {
			return !gomander.OnBeforeClose(app.Context())
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	ctx := app.Context()
	runtimeFacade := facade.NewRuntimeFacade(app)
	dialogs := dialog.NewWailsDialogs(app.Dialog)

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Gomander",
		Width:            1024,
		Height:           768,
		BackgroundColour: application.NewRGB(27, 38, 54),
		KeyBindings:      macOSZoomKeyBindings(runtime.GOOS, zoom, runtimeFacade.LogError),
	})
	closeWindowByQuitting(window, app)
	if runtime.GOOS == "darwin" {
		app.Menu.Set(macOSApplicationMenu(zoom, runtimeFacade.LogError))
	}

	gormDb := configDB(ctx)
	gomander, useCases = buildDeps(gormDb, ctx, dialogs, runtimeFacade, facade.NewWindowFacade(window))

	window.OnWindowEvent(events.Common.WindowRuntimeReady, func(*application.WindowEvent) {
		if err := useCases.RestoreZoom.Execute(); err != nil {
			runtimeFacade.LogError(err.Error())
		}
	})

	gomander.RegisterHandlers()

	if err := gomander.Startup(ctx); err != nil {
		panic(err)
	}

	app.RegisterService(application.NewService(path.NewUiPathHelper()))
	app.RegisterService(application.NewService(NewWailsControllers(useCases)))
	app.RegisterService(application.NewService(fs.NewUIFsHelper(dialogs, runtimeFacade)))
	app.RegisterService(application.NewService(os_internal.NewUIOsHelper()))

	server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
	go func() {
		if err := server.RegisterHandlers(); err != nil {
			panic(err)
		}
		server.Start()
	}()

	if err := app.Run(); err != nil {
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

func buildDeps(gormDb *gorm.DB, ctx context.Context, dialogs dialog.Dialogs, runtimeFacade facade.RuntimeFacade, window configusecases.Window) (*internalapp.App, usecases.Registry) {
	// Initialize deps
	l := logger.NewDefaultLogger(runtimeFacade)
	ee := event.NewDefaultEventEmitter(runtimeFacade)
	r := runner.NewDefaultRunner(l, ee)
	releaseFeed := releaseinfrastructure.NewGithubReleaseFeed(ctx, facade.DefaultIOFacade{}, releaseinfrastructure.DefaultLatestReleaseUrl)
	releaseDownloader := releaseinfrastructure.NewGithubReleaseDownloader(ctx, facade.DefaultOSFacade{}, facade.DefaultIOFacade{}, releaseinfrastructure.DefaultBinaryDownloadBaseUrl)
	releaseInstaller := releaseinfrastructure.NewOSReleaseInstaller(facade.DefaultOSFacade{}, facade.DefaultOpenFacade{})

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
	zoom := configusecases.NewZoom(configRepo, window)
	zoomIn := configusecases.NewZoomIn(zoom)
	zoomOut := configusecases.NewZoomOut(zoom)
	resetZoom := configusecases.NewResetZoom(zoom)
	restoreZoom := configusecases.NewRestoreZoom(configRepo, window)
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
	projectFileV1 := projectinfrastructure.NewProjectFileV1(facade.DefaultFsFacade{})
	exportProject := projectusecases.NewExportProject(projectRepo, commandRepo, commandGroupRepo, dialogs, projectFileV1)
	importProject := projectusecases.NewImportProject(unitOfWork)
	getProjectToImport := projectusecases.NewGetProjectToImport(dialogs, map[projectusecases.FileType]projectusecases.BlueprintReader{
		projectusecases.FileTypeGomander:    projectFileV1,
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
	installReleaseAndQuit := releaseusecases.NewInstallReleaseAndQuit(releaseInstaller, runtimeFacade)

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
		ZoomIn:         zoomIn,
		ZoomOut:        zoomOut,
		ResetZoom:      resetZoom,
		RestoreZoom:    restoreZoom,
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
