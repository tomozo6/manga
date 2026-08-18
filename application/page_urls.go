package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type pageResponse struct {
	Number   int    `json:"number"`
	ImageURL string `json:"imageUrl"`
}

func (a app) issuePageBatch(ctx context.Context, item manga, selectedVolume volume, start int) ([]pageResponse, error) {
	end := min(start+pageURLBatchSize-1, selectedVolume.PageCount)
	pages := make([]pageResponse, end-start+1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once
	for number := start; number <= end; number++ {
		index := number - start
		wg.Add(1)
		go func(number, index int) {
			defer wg.Done()
			key := fmt.Sprintf("manga/%s/%s/%03d.%s", item.ID, selectedVolume.ID, number, selectedVolume.PageExtension)
			imageURL, err := a.gcsMediaSigner.Issue(ctx, key, time.Now())
			if err != nil {
				errOnce.Do(func() {
					firstErr = err
					cancel()
				})
				return
			}
			pages[index] = pageResponse{Number: number, ImageURL: imageURL}
		}(number, index)
	}
	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return pages, nil
}
