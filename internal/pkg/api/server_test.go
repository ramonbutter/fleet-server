// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

//go:build !integration

package api

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	fbuild "github.com/elastic/fleet-server/v7/internal/pkg/build"
	"github.com/elastic/fleet-server/v7/internal/pkg/cache"
	"github.com/elastic/fleet-server/v7/internal/pkg/checkin"
	"github.com/elastic/fleet-server/v7/internal/pkg/config"
	"github.com/elastic/fleet-server/v7/internal/pkg/monitor/mock"
	"github.com/elastic/fleet-server/v7/internal/pkg/policy"
	ftesting "github.com/elastic/fleet-server/v7/internal/pkg/testing"
)

func Test_server_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port, err := ftesting.FreePort()
	require.NoError(t, err)
	cfg := &config.Server{}
	cfg.InitDefaults()
	cfg.Host = "localhost"
	cfg.Port = port
	addr := cfg.BindEndpoints()[0]

	verCon := mustBuildConstraints("8.0.0")
	c, err := cache.New(config.Cache{NumCounters: 100, MaxCost: 100000})
	require.NoError(t, err)
	bulker := ftesting.NewMockBulk()
	pim := mock.NewMockMonitor()
	pm := policy.NewMonitor(bulker, pim, 5*time.Millisecond)
	bc := checkin.NewBulk(nil)
	ct := NewCheckinT(verCon, cfg, c, bc, pm, nil, nil, nil, nil)
	et, err := NewEnrollerT(verCon, cfg, nil, c)
	require.NoError(t, err)

	srv := NewServer(addr, cfg, ct, et, nil, nil, nil, nil, fbuild.Info{}, nil, nil, nil)
	errCh := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = srv.Run(ctx)
		wg.Done()
	}()
	var errFromChan error
	select {
	case err := <-errCh:
		errFromChan = err
	case <-time.After(500 * time.Millisecond):
		break
	}
	cancel()
	wg.Wait()
	require.NoError(t, errFromChan)
	if !errors.Is(err, http.ErrServerClosed) {
		require.NoError(t, err)
	}
}
