package server

import (
	"context"

	"cloud.google.com/go/datastore"
	"golang.org/x/crypto/acme/autocert"
)

// https://www.captaincodeman.com/2017/05/07/automatic-https-with-free-ssl-certificates-using-go-lets-encrypt

const (
	certCacheKind = "letsencrypt"
)

// LetsencryptCert -
type LetsencryptCert struct {
	Data []byte `datastore:"data,noindex"`
}

type datastoreCertCache struct {
	client *datastore.Client
}

func newDatastoreCertCache(client *datastore.Client) *datastoreCertCache {
	return &datastoreCertCache{client: client}
}

// Get -
func (d *datastoreCertCache) Get(ctx context.Context, key string) ([]byte, error) {
	var cert LetsencryptCert
	k := datastore.NameKey(certCacheKind, key, nil)
	if err := d.client.Get(ctx, k, &cert); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, autocert.ErrCacheMiss
		}
		return nil, err
	}
	return cert.Data, nil
}

// Put -
func (d *datastoreCertCache) Put(ctx context.Context, key string, data []byte) error {
	k := datastore.NameKey(certCacheKind, key, nil)
	cert := LetsencryptCert{
		Data: data,
	}
	if _, err := d.client.Put(ctx, k, &cert); err != nil {
		return err
	}
	return nil
}

// Delete -
func (d *datastoreCertCache) Delete(ctx context.Context, key string) error {
	k := datastore.NameKey(certCacheKind, key, nil)
	if err := d.client.Delete(ctx, k); err != nil {
		return err
	}
	return nil
}
