// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"path"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/publish"
	"github.com/sigstore/cosign/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/cmd/cosign/cli/sign"
	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	opts, err := buildOpts()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bi, err := build.NewGo(ctx, ".", opts...)
	if err != nil {
		t.Fatal(err)
	}

	importpath, err := bi.QualifyImport("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	img, err := bi.Build(ctx, importpath)
	if err != nil {
		t.Fatal(err)
	}

	digest, err := img.Digest()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Digest:%s", digest)

	p, err := publish.NewDefault("localhost:5001/test",
		publish.Insecure(true), publish.WithNamer(func(base, importpath string) string {
			hasher := md5.New() // nolint: gosec // No strong cryptography needed.
			hasher.Write([]byte(importpath))
			return path.Join(base, path.Base(importpath)+"-"+hex.EncodeToString(hasher.Sum(nil)))
		}))
	if err != nil {
		t.Fatal(err)
	}
	pres, err := p.Publish(ctx, img, importpath)
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)

	cwd, _ := os.Getwd()
	keyRef := path.Join(cwd, "testdata", "cosign.key")

	ro := &options.RootOptions{}
	ko := options.KeyOpts{
		KeyRef: keyRef,
		PassFunc: func(b bool) ([]byte, error) {
			pw, ok := os.LookupEnv("COSIGN_PASSWORD")
			if ok {
				return []byte(pw), nil
			}
			return []byte("password"), nil
		},
	}
	regOpts := options.RegistryOptions{
		AllowInsecure: true,
	}
	imgs := []string{pres.Name()}
	anMap := make(map[string]interface{})
	//sign
	err = sign.SignCmd(ro, ko, regOpts, anMap, imgs, "", "", true, "", "", "", false, false, "", false)
	assert.NoError(t, err)
}

func buildOpts() ([]build.Option, error) {
	var opts []build.Option
	base, err := random.Image(1024, 1)
	if err != nil {
		return nil, err
	}

	opts = append(opts,
		build.WithBaseImages(func(ctx context.Context, s string) (name.Reference, build.Result, error) {
			return name.MustParseReference("gcr.io/distroless/static:nonroot"), base, nil
		}))

	opts = append(opts,
		build.WithPlatforms("all"))

	return opts, nil
}
