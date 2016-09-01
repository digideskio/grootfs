// This file was generated by counterfeiter
package storefakes

import (
	"sync"

	"code.cloudfoundry.org/grootfs/store"
	"code.cloudfoundry.org/lager"
)

type FakeSnapshotDriver struct {
	SnapshotStub        func(logger lager.Logger, fromPath, toPath string) error
	snapshotMutex       sync.RWMutex
	snapshotArgsForCall []struct {
		logger   lager.Logger
		fromPath string
		toPath   string
	}
	snapshotReturns struct {
		result1 error
	}
	ApplyDiskLimitStub        func(logger lager.Logger, path string, diskLimit int64) error
	applyDiskLimitMutex       sync.RWMutex
	applyDiskLimitArgsForCall []struct {
		logger    lager.Logger
		path      string
		diskLimit int64
	}
	applyDiskLimitReturns struct {
		result1 error
	}
	DestroyStub        func(logger lager.Logger, path string) error
	destroyMutex       sync.RWMutex
	destroyArgsForCall []struct {
		logger lager.Logger
		path   string
	}
	destroyReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSnapshotDriver) Snapshot(logger lager.Logger, fromPath string, toPath string) error {
	fake.snapshotMutex.Lock()
	fake.snapshotArgsForCall = append(fake.snapshotArgsForCall, struct {
		logger   lager.Logger
		fromPath string
		toPath   string
	}{logger, fromPath, toPath})
	fake.recordInvocation("Snapshot", []interface{}{logger, fromPath, toPath})
	fake.snapshotMutex.Unlock()
	if fake.SnapshotStub != nil {
		return fake.SnapshotStub(logger, fromPath, toPath)
	} else {
		return fake.snapshotReturns.result1
	}
}

func (fake *FakeSnapshotDriver) SnapshotCallCount() int {
	fake.snapshotMutex.RLock()
	defer fake.snapshotMutex.RUnlock()
	return len(fake.snapshotArgsForCall)
}

func (fake *FakeSnapshotDriver) SnapshotArgsForCall(i int) (lager.Logger, string, string) {
	fake.snapshotMutex.RLock()
	defer fake.snapshotMutex.RUnlock()
	return fake.snapshotArgsForCall[i].logger, fake.snapshotArgsForCall[i].fromPath, fake.snapshotArgsForCall[i].toPath
}

func (fake *FakeSnapshotDriver) SnapshotReturns(result1 error) {
	fake.SnapshotStub = nil
	fake.snapshotReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSnapshotDriver) ApplyDiskLimit(logger lager.Logger, path string, diskLimit int64) error {
	fake.applyDiskLimitMutex.Lock()
	fake.applyDiskLimitArgsForCall = append(fake.applyDiskLimitArgsForCall, struct {
		logger    lager.Logger
		path      string
		diskLimit int64
	}{logger, path, diskLimit})
	fake.recordInvocation("ApplyDiskLimit", []interface{}{logger, path, diskLimit})
	fake.applyDiskLimitMutex.Unlock()
	if fake.ApplyDiskLimitStub != nil {
		return fake.ApplyDiskLimitStub(logger, path, diskLimit)
	} else {
		return fake.applyDiskLimitReturns.result1
	}
}

func (fake *FakeSnapshotDriver) ApplyDiskLimitCallCount() int {
	fake.applyDiskLimitMutex.RLock()
	defer fake.applyDiskLimitMutex.RUnlock()
	return len(fake.applyDiskLimitArgsForCall)
}

func (fake *FakeSnapshotDriver) ApplyDiskLimitArgsForCall(i int) (lager.Logger, string, int64) {
	fake.applyDiskLimitMutex.RLock()
	defer fake.applyDiskLimitMutex.RUnlock()
	return fake.applyDiskLimitArgsForCall[i].logger, fake.applyDiskLimitArgsForCall[i].path, fake.applyDiskLimitArgsForCall[i].diskLimit
}

func (fake *FakeSnapshotDriver) ApplyDiskLimitReturns(result1 error) {
	fake.ApplyDiskLimitStub = nil
	fake.applyDiskLimitReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSnapshotDriver) Destroy(logger lager.Logger, path string) error {
	fake.destroyMutex.Lock()
	fake.destroyArgsForCall = append(fake.destroyArgsForCall, struct {
		logger lager.Logger
		path   string
	}{logger, path})
	fake.recordInvocation("Destroy", []interface{}{logger, path})
	fake.destroyMutex.Unlock()
	if fake.DestroyStub != nil {
		return fake.DestroyStub(logger, path)
	} else {
		return fake.destroyReturns.result1
	}
}

func (fake *FakeSnapshotDriver) DestroyCallCount() int {
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return len(fake.destroyArgsForCall)
}

func (fake *FakeSnapshotDriver) DestroyArgsForCall(i int) (lager.Logger, string) {
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return fake.destroyArgsForCall[i].logger, fake.destroyArgsForCall[i].path
}

func (fake *FakeSnapshotDriver) DestroyReturns(result1 error) {
	fake.DestroyStub = nil
	fake.destroyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSnapshotDriver) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.snapshotMutex.RLock()
	defer fake.snapshotMutex.RUnlock()
	fake.applyDiskLimitMutex.RLock()
	defer fake.applyDiskLimitMutex.RUnlock()
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeSnapshotDriver) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ store.SnapshotDriver = new(FakeSnapshotDriver)