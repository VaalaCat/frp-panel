package rbac

import (
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func InitializeCasbin(ctx *app.Context, db *gorm.DB) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		logger.Logger(ctx).Fatalf("error creating Casbin GORM adapter: %v", err)
		return nil, err
	}

	// Load the Casbin model from file
	m, err := model.NewModelFromString(RBAC_MODEL)
	if err != nil {
		logger.Logger(ctx).Fatalf("error loading Casbin model: %v", err)
		return nil, err
	}

	// Create the Casbin enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		logger.Logger(ctx).Fatalf("error creating Casbin enforcer: %v", err)
		return nil, err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		logger.Logger(ctx).WithError(err).Warnf("could not load Casbin policy from DB")
		// return nil, err
	}

	enforcer.EnableAutoSave(true)

	logger.Logger(ctx).Infof("Casbin Enforcer initialized successfully.")
	return enforcer, nil
}
