package mapper

import (
	"testing"

	"github.com/compose-spec/compose-go/v2/types"
)

func TestUnit_DependsOn(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", DependsOn: types.DependsOnConfig{"db": {Condition: "service_started", Required: true}, "redis": {Condition: "service_started", Required: false}}}
	dirs := Unit(svc)

	assertDirective(t, dirs, "Requires", "db.container")
	assertDirective(t, dirs, "Wants", "redis.container")
	assertDirective(t, dirs, "After", "db.container")
	assertDirective(t, dirs, "After", "redis.container")
}

func TestUnit_DependsOn_Restart(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", DependsOn: types.DependsOnConfig{"db": {Condition: "service_started", Required: true, Restart: true}}}
	dirs := Unit(svc)

	assertDirective(t, dirs, "Requires", "db.container")
	assertDirective(t, dirs, "BindsTo", "db.container")
	assertDirective(t, dirs, "After", "db.container")
}

func TestUnit_DependsOn_Empty(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", DependsOn: types.DependsOnConfig{}}
	dirs := Unit(svc)
	if len(dirs) != 0 {
		t.Fatalf("expected no directives, got %v", dirs)
	}
}
