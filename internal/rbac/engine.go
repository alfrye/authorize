package rbac

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	sqladapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/dgraph-io/ristretto"
)

type Engine struct {
	enforcer *casbin.Enforcer
	cache    *ristretto.Cache
}

func NewEngine(policyPath string) (*Engine, error) {
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
	if err != nil {
		return nil, err
	}
	adapter := sqladapter.NewFilteredAdapter(policyPath)
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     32 << 20,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}
	return &Engine{enforcer: enforcer, cache: cache}, nil
}

func (e *Engine) CheckPermission(subject, resource, action string) (bool, error) {
	cacheKey := subject + ":" + resource + ":" + action
	if val, found := e.cache.Get(cacheKey); found {
		return val.(bool), nil
	}
	result, err := e.enforcer.Enforce(subject, resource, action)
	if err != nil {
		return false, err
	}
	e.cache.Set(cacheKey, result, 1)
	return result, nil
}

func (e *Engine) AddRoleForUser(user, role string) (bool, error) {
	e.cache.Clear()
	return e.enforcer.AddRoleForUser(user, role)
}

func (e *Engine) AddPolicy(role, resource, action string) (bool, error) {
	e.cache.Clear()
	return e.enforcer.AddPolicy(role, resource, action)
}
