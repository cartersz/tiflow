// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package tp

import (
	"context"

	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/cdc/scheduler"
	"github.com/pingcap/tiflow/cdc/scheduler/tp/schedulepb"
)

type schedule interface {
	Name() string
	Schedule(
		currentTables []model.TableID,
		captures map[model.CaptureID]*model.CaptureInfo,
		captureTables map[model.CaptureID]captureTableSet,
	) []*scheduleTask
}

type coordinator struct {
	trans   transport
	manager *replicationManager
	// balancer and drainer
	scheduler []schedule
}

func NewCoordinator() scheduler.ScheduleDispatcher {
	return nil
}

func (c *coordinator) Tick(
	ctx context.Context,
	// Latest global checkpoint of the changefeed
	checkpointTs model.Ts,
	// All tables that SHOULD be replicated (or started) at the current checkpoint.
	currentTables []model.TableID,
	// All captures that are alive according to the latest Etcd states.
	aliveCaptures map[model.CaptureID]*model.CaptureInfo,
) (newCheckpointTs, newResolvedTs model.Ts, err error) {
	captureTables := c.manager.captureTableSets()
	allTasks := make([]*scheduleTask, 0)
	for _, sched := range c.scheduler {
		tasks := sched.Schedule(currentTables, aliveCaptures, captureTables)
		allTasks = append(allTasks, tasks...)
	}
	msgs := c.recvMessages()
	c.manager.poll(ctx, checkpointTs, currentTables, aliveCaptures, msgs, allTasks)
	c.sendMessages(nil)

	return scheduler.CheckpointCannotProceed, scheduler.CheckpointCannotProceed, nil
}

func (c *coordinator) recvMessages() []*schedulepb.Message {
	return nil
}

func (c *coordinator) sendMessages(msgs []*schedulepb.Message) error {
	return nil
}

func (s *replicationManager) recvTasks() []*scheduleTask {
	// TODO: do not receive tasks if len(s.runningTasks) > threshold.
	return nil
}
