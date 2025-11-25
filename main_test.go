package main

import (
	"context"
	"flag"
	"os"
	"testing"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-legacy-cache-api/features/steps"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var componentFlag = flag.Bool("component", false, "perform component tests")

type ComponentTest struct {
	MongoFeature *componenttest.MongoFeature
}

func (f *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	mongoURI, err := f.MongoFeature.GetConnectionString()
	if err != nil {
		panic(err)
	}
	mongoDatabaseName := f.MongoFeature.Database.Name()

	component, err := steps.NewComponent(mongoURI, mongoDatabaseName)
	if err != nil {
		panic(err)
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		component.Reset()

		if err := f.MongoFeature.Reset(); err != nil {
			panic(err)
		}

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if closeErr := component.Close(); closeErr != nil {
			panic(closeErr)
		}

		return ctx, nil
	})

	component.RegisterSteps(ctx)
	f.MongoFeature.RegisterSteps(ctx)
}

func (f *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	const MongoVersion = "4.4.8"
	const DatabaseName = "testing"

	ctx.BeforeSuite(func() {
		f.MongoFeature = componenttest.NewMongoFeature(componenttest.MongoOptions{MongoVersion: MongoVersion, DatabaseName: DatabaseName})
	})
	ctx.AfterSuite(func() {
		err := f.MongoFeature.Close()
		if err != nil {
			panic(err)
		}
	})
}

func TestComponent(t *testing.T) {
	if *componentFlag {
		status := 0

		var opts = godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  flag.Args(),
			Strict: true,
		}

		f := &ComponentTest{}

		status = godog.TestSuite{
			Name:                 "feature_tests",
			ScenarioInitializer:  f.InitializeScenario,
			TestSuiteInitializer: f.InitializeTestSuite,
			Options:              &opts,
		}.Run()

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
