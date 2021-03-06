// This file was generated by counterfeiter
package image_pullerfakes

import (
	"io"
	"net/url"
	"sync"

	"code.cloudfoundry.org/grootfs/image_puller"
	"code.cloudfoundry.org/lager"
)

type FakeFetcher struct {
	ImageInfoStub        func(logger lager.Logger, imageURL *url.URL) (image_puller.ImageInfo, error)
	imageInfoMutex       sync.RWMutex
	imageInfoArgsForCall []struct {
		logger   lager.Logger
		imageURL *url.URL
	}
	imageInfoReturns struct {
		result1 image_puller.ImageInfo
		result2 error
	}
	StreamBlobStub        func(logger lager.Logger, imageURL *url.URL, source string) (io.ReadCloser, int64, error)
	streamBlobMutex       sync.RWMutex
	streamBlobArgsForCall []struct {
		logger   lager.Logger
		imageURL *url.URL
		source   string
	}
	streamBlobReturns struct {
		result1 io.ReadCloser
		result2 int64
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeFetcher) ImageInfo(logger lager.Logger, imageURL *url.URL) (image_puller.ImageInfo, error) {
	fake.imageInfoMutex.Lock()
	fake.imageInfoArgsForCall = append(fake.imageInfoArgsForCall, struct {
		logger   lager.Logger
		imageURL *url.URL
	}{logger, imageURL})
	fake.recordInvocation("ImageInfo", []interface{}{logger, imageURL})
	fake.imageInfoMutex.Unlock()
	if fake.ImageInfoStub != nil {
		return fake.ImageInfoStub(logger, imageURL)
	} else {
		return fake.imageInfoReturns.result1, fake.imageInfoReturns.result2
	}
}

func (fake *FakeFetcher) ImageInfoCallCount() int {
	fake.imageInfoMutex.RLock()
	defer fake.imageInfoMutex.RUnlock()
	return len(fake.imageInfoArgsForCall)
}

func (fake *FakeFetcher) ImageInfoArgsForCall(i int) (lager.Logger, *url.URL) {
	fake.imageInfoMutex.RLock()
	defer fake.imageInfoMutex.RUnlock()
	return fake.imageInfoArgsForCall[i].logger, fake.imageInfoArgsForCall[i].imageURL
}

func (fake *FakeFetcher) ImageInfoReturns(result1 image_puller.ImageInfo, result2 error) {
	fake.ImageInfoStub = nil
	fake.imageInfoReturns = struct {
		result1 image_puller.ImageInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeFetcher) StreamBlob(logger lager.Logger, imageURL *url.URL, source string) (io.ReadCloser, int64, error) {
	fake.streamBlobMutex.Lock()
	fake.streamBlobArgsForCall = append(fake.streamBlobArgsForCall, struct {
		logger   lager.Logger
		imageURL *url.URL
		source   string
	}{logger, imageURL, source})
	fake.recordInvocation("StreamBlob", []interface{}{logger, imageURL, source})
	fake.streamBlobMutex.Unlock()
	if fake.StreamBlobStub != nil {
		return fake.StreamBlobStub(logger, imageURL, source)
	} else {
		return fake.streamBlobReturns.result1, fake.streamBlobReturns.result2, fake.streamBlobReturns.result3
	}
}

func (fake *FakeFetcher) StreamBlobCallCount() int {
	fake.streamBlobMutex.RLock()
	defer fake.streamBlobMutex.RUnlock()
	return len(fake.streamBlobArgsForCall)
}

func (fake *FakeFetcher) StreamBlobArgsForCall(i int) (lager.Logger, *url.URL, string) {
	fake.streamBlobMutex.RLock()
	defer fake.streamBlobMutex.RUnlock()
	return fake.streamBlobArgsForCall[i].logger, fake.streamBlobArgsForCall[i].imageURL, fake.streamBlobArgsForCall[i].source
}

func (fake *FakeFetcher) StreamBlobReturns(result1 io.ReadCloser, result2 int64, result3 error) {
	fake.StreamBlobStub = nil
	fake.streamBlobReturns = struct {
		result1 io.ReadCloser
		result2 int64
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeFetcher) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.imageInfoMutex.RLock()
	defer fake.imageInfoMutex.RUnlock()
	fake.streamBlobMutex.RLock()
	defer fake.streamBlobMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeFetcher) recordInvocation(key string, args []interface{}) {
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

var _ image_puller.Fetcher = new(FakeFetcher)
