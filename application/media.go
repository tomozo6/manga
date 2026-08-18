package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iamcredentials/v1"
)

const (
	pageURLBatchSize               = 8
	mangaImageBucket               = "tomozo-manga-images"
	mangaImageSignerServiceAccount = "manga-media-signer@tomozo6.iam.gserviceaccount.com"
	mediaURLTTL                    = time.Hour
)

type gcsMediaSigner struct {
	signBytes func(context.Context, []byte) ([]byte, error)
}

func newGCSMediaSigner(ctx context.Context) (gcsMediaSigner, error) {
	service, err := iamcredentials.NewService(ctx)
	if err != nil {
		return gcsMediaSigner{}, fmt.Errorf("create IAM Credentials client: %w", err)
	}
	return gcsMediaSigner{signBytes: func(ctx context.Context, payload []byte) ([]byte, error) {
		response, err := service.Projects.ServiceAccounts.SignBlob(
			"projects/-/serviceAccounts/"+mangaImageSignerServiceAccount,
			&iamcredentials.SignBlobRequest{Payload: base64.StdEncoding.EncodeToString(payload)},
		).Context(ctx).Do()
		if err != nil {
			return nil, err
		}
		return base64.StdEncoding.DecodeString(response.SignedBlob)
	}}, nil
}

func (s gcsMediaSigner) Issue(ctx context.Context, key string, now time.Time) (string, error) {
	url, err := storage.SignedURL(mangaImageBucket, key, &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         http.MethodGet,
		Expires:        now.Add(mediaURLTTL),
		GoogleAccessID: mangaImageSignerServiceAccount,
		SignBytes:      func(payload []byte) ([]byte, error) { return s.signBytes(ctx, payload) },
	})
	if err != nil {
		return "", fmt.Errorf("sign GCS URL for %q: %w", key, err)
	}
	return url, nil
}
