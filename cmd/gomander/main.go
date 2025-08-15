package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	commmandinfrastructure "gomander/internal/command/infrastructure"
	commandgroupinfrastructure "gomander/internal/commandgroup/infrastructure"
	configinfrastructure "gomander/internal/config/infrastructure"
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
	gormDb, err := gorm.Open(sqlite.Open(getDbFile()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db, err := gormDb.DB()
	if err != nil {
		panic(err)
	}

	// Execute migrations
	err = goose.SetDialect("sqlite3")
	if err != nil {
		panic(err)
	}

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
	l := logger.NewDefaultLogger(ctx)
	ee := event.NewDefaultEventEmitter(ctx)
	r := runner.NewDefaulRunner(l, ee)

	// Initialize repos
	comandRepo := commmandinfrastructure.NewGormCommandRepository(gormDb, ctx)
	comandGroupRepo := commandgroupinfrastructure.NewGormCommandGroupRepository(gormDb, ctx)
	projectRepo := projectinfrastructure.NewGormProjectRepository(gormDb, ctx)
	configRepo := configinfrastructure.NewGormConfigRepository(gormDb, ctx)

	// Initialize event emitter
	app.LoadDependencies(l, ee, r, comandRepo, comandGroupRepo, projectRepo, configRepo)
}
