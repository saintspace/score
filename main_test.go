package main

import (
	"embed"
	"errors"
	"flag"
	"os"
	"score/app/runners/confirmer"
	"score/app/runners/preauth"
	"score/app/runners/server"
	"score/app/runners/worker"
	"testing"

	"github.com/stretchr/testify/require"
)

// Mock the app functions
var mockRunServer = func(embed.FS) error {
	return nil
}

var mockRunWorker = func() error {
	return nil
}

var mockRunPostAccountConfirmationHandler = func() error {
	return nil
}

var mockRunPreAuthenticationHandler = func() error {
	return nil
}

func TestMain(m *testing.M) {
	server.Run = mockRunServer
	worker.Run = mockRunWorker
	confirmer.Run = mockRunPostAccountConfirmationHandler
	preauth.Run = mockRunPreAuthenticationHandler

	os.Exit(m.Run())
}

func TestGetExecutionMode(t *testing.T) {
	// Test the environment variable
	os.Setenv("SCORE_MODE", "worker")
	executionMode := getExecutionMode()
	require.Equal(t, ExecutionMode("worker"), executionMode, "Execution mode should be worker")
	os.Unsetenv("SCORE_MODE")

	// Reset the flags before the next test
	os.Args = []string{"cmd"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test the command line flag
	os.Args = []string{"cmd", "-mode=confirmer"}
	executionMode = getExecutionMode()
	require.Equal(t, ExecutionMode("confirmer"), executionMode, "Execution mode should be confirmer")

	// Reset the flags before the next test
	os.Args = []string{"cmd"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test the default value
	executionMode = getExecutionMode()
	require.Equal(t, ExecutionMode("server"), executionMode, "Execution mode should be server by default")
}

func TestRunApp(t *testing.T) {
	t.Run("server", func(t *testing.T) {
		originalRun := server.Run
		server.Run = func(fs embed.FS) error {
			require.NotNil(t, fs, "embed.FS should not be nil")
			return nil
		}
		defer func() { server.Run = originalRun }()
		err := RunApp(ServerMode)
		require.NoError(t, err, "Server mode should not return an error")
	})

	t.Run("worker", func(t *testing.T) {
		originalRun := worker.Run
		worker.Run = func() error {
			return nil
		}
		err := RunApp(WorkerMode)
		defer func() { worker.Run = originalRun }()
		require.NoError(t, err, "Worker mode should not return an error")
	})

	t.Run("confirmer", func(t *testing.T) {
		originalRun := confirmer.Run
		confirmer.Run = func() error {
			return nil
		}
		defer func() { confirmer.Run = originalRun }()
		err := RunApp(ConfirmerMode)
		require.NoError(t, err, "Confirmer mode should not return an error")
	})

	t.Run("preauth", func(t *testing.T) {
		originalRun := preauth.Run
		preauth.Run = func() error {
			return nil
		}
		defer func() { preauth.Run = originalRun }()
		err := RunApp(PreauthMode)
		require.NoError(t, err, "Preauth mode should not return an error")
	})

	t.Run("default", func(t *testing.T) {
		originalRun := server.Run
		server.Run = func(fs embed.FS) error {
			require.NotNil(t, fs, "embed.FS should not be nil")
			return nil
		}
		defer func() { server.Run = originalRun }()
		err := RunApp(ExecutionMode("invalid"))
		require.NoError(t, err, "Default mode should not return an error")
	})

	t.Run("error", func(t *testing.T) {
		originalRun := server.Run
		server.Run = func(fs embed.FS) error {
			return errors.New("some error")
		}
		defer func() { server.Run = originalRun }()
		err := RunApp(ServerMode)
		require.Error(t, err, "Should return an error if RunServer fails")
	})
}
