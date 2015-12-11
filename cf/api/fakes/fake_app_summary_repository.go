// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/models"
)

type FakeAppSummaryRepository struct {
	GetSummariesInCurrentSpaceStub        func() (apps []models.Application, apiErr error)
	getSummariesInCurrentSpaceMutex       sync.RWMutex
	getSummariesInCurrentSpaceArgsForCall []struct{}
	getSummariesInCurrentSpaceReturns struct {
		result1 []models.Application
		result2 error
	}
	GetSummaryStub        func(appGuid string) (summary models.Application, apiErr error)
	getSummaryMutex       sync.RWMutex
	getSummaryArgsForCall []struct {
		appGuid string
	}
	getSummaryReturns struct {
		result1 models.Application
		result2 error
	}
}

func (fake *FakeAppSummaryRepository) GetSummariesInCurrentSpace() (apps []models.Application, apiErr error) {
	fake.getSummariesInCurrentSpaceMutex.Lock()
	fake.getSummariesInCurrentSpaceArgsForCall = append(fake.getSummariesInCurrentSpaceArgsForCall, struct{}{})
	fake.getSummariesInCurrentSpaceMutex.Unlock()
	if fake.GetSummariesInCurrentSpaceStub != nil {
		return fake.GetSummariesInCurrentSpaceStub()
	} else {
		return fake.getSummariesInCurrentSpaceReturns.result1, fake.getSummariesInCurrentSpaceReturns.result2
	}
}

func (fake *FakeAppSummaryRepository) GetSummariesInCurrentSpaceCallCount() int {
	fake.getSummariesInCurrentSpaceMutex.RLock()
	defer fake.getSummariesInCurrentSpaceMutex.RUnlock()
	return len(fake.getSummariesInCurrentSpaceArgsForCall)
}

func (fake *FakeAppSummaryRepository) GetSummariesInCurrentSpaceReturns(result1 []models.Application, result2 error) {
	fake.GetSummariesInCurrentSpaceStub = nil
	fake.getSummariesInCurrentSpaceReturns = struct {
		result1 []models.Application
		result2 error
	}{result1, result2}
}

func (fake *FakeAppSummaryRepository) GetSummary(appGuid string) (summary models.Application, apiErr error) {
	fake.getSummaryMutex.Lock()
	fake.getSummaryArgsForCall = append(fake.getSummaryArgsForCall, struct {
		appGuid string
	}{appGuid})
	fake.getSummaryMutex.Unlock()
	if fake.GetSummaryStub != nil {
		return fake.GetSummaryStub(appGuid)
	} else {
		return fake.getSummaryReturns.result1, fake.getSummaryReturns.result2
	}
}

func (fake *FakeAppSummaryRepository) GetSummaryCallCount() int {
	fake.getSummaryMutex.RLock()
	defer fake.getSummaryMutex.RUnlock()
	return len(fake.getSummaryArgsForCall)
}

func (fake *FakeAppSummaryRepository) GetSummaryArgsForCall(i int) string {
	fake.getSummaryMutex.RLock()
	defer fake.getSummaryMutex.RUnlock()
	return fake.getSummaryArgsForCall[i].appGuid
}

func (fake *FakeAppSummaryRepository) GetSummaryReturns(result1 models.Application, result2 error) {
	fake.GetSummaryStub = nil
	fake.getSummaryReturns = struct {
		result1 models.Application
		result2 error
	}{result1, result2}
}

var _ api.AppSummaryRepository = new(FakeAppSummaryRepository)
