package accessor

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	schemaName = "common"
	tableName  = "casbinrule"
	prefix     = ""
)

func NewEnforcerFromGORM(db *gorm.DB, configPath string) (*casbin.Enforcer, error) {
	db.Exec(fmt.Sprintf("SET search_path TO %s", schemaName))
	gormadapter.TurnOffAutoMigrate(db)

	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, prefix, tableName)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize GORM adapter: %w", err)
	}

	casbinModel, err := model.NewModelFromFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize Casbin model: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(casbinModel, adapter)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize Casbin Enforcer: %w", err)
	}
	return enforcer, nil
}

type casbinAccessor struct {
	enforcer casbin.IEnforcer
}

func NewCasbinAccessor(enforcer casbin.IEnforcer) *casbinAccessor {
	return &casbinAccessor{enforcer: enforcer}
}

func (c *casbinAccessor) CheckAccess(role string, res, act string) (AccessStatus, error) {
	isAccessable, err := c.enforcer.Enforce(role, res, act)
	if err != nil {
		return AccessUnknown, fmt.Errorf("enforce error: %w", err)
	}
	if !isAccessable {
		return AccessDenied, nil
	}
	return AccessGranted, nil
}
