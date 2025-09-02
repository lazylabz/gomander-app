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

	"gomander/internal/command/application/handlers"
	commmandinfrastructure "gomander/internal/command/infrastructure"
	commandgrouphandlers "gomander/internal/commandgroup/application/handlers"
	commandgroupinfrastructure "gomander/internal/commandgroup/infrastructure"
	configinfrastructure "gomander/internal/config/infrastructure"
	"gomander/internal/eventbus"
	"gomander/internal/facade"
	logger "gomander/internal/logger"
	projectinfrastructure "gomander/internal/project/infrastructure"
	"gomander/internal/runner"
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
	uiPathHelper := path.NewUiPathHelper()

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
		},
		Bind: []interface{}{
			app,
			uiPathHelper,
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
	cleanCommandGroupsOnCommandDeleted := commandgrouphandlers.NewDefaultCleanCommandGroupsOnCommandDeleted(commandGroupRepo)
	cleanCommandGroupsOnProjectDeleted := commandgrouphandlers.NewDefaultCleanCommandGroupsOnProjectDeleted(commandGroupRepo)
	cleanCommandsOnProjectDeleted := handlers.NewDefaultCleanCommandOnProjectDeleted(commandRepo)
	addCommandToGroupOnCommandDuplicated := commandgrouphandlers.NewDefaultAddCommandToGroupOnCommandDuplicated(commandRepo, commandGroupRepo)

	eventBus := eventbus.NewInMemoryEventBus()

	// Initialize event emitter
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
	})
}
