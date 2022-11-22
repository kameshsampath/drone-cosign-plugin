// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/google/ko/pkg/commands"
	kooptions "github.com/google/ko/pkg/commands/options"
	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	cwd, _ := os.Getwd()
	ctx := context.Background()

	namespace := "base"
	s, err := registryServerWithImage(namespace)
	if err != nil {
		t.Fatalf("could not create test registry server: %v", err)
	}
	defer s.Close()
	repo := "localhost:5001"
	baseImage := fmt.Sprintf("%s/%s", repo, namespace)

	bo := &kooptions.BuildOptions{
		BaseImage:        baseImage,
		InsecureRegistry: true,
		UserAgent:        "ko",
		ConcurrentBuilds: 1,
		Platforms:        []string{"all"},
	}

	bi, err := commands.NewBuilder(ctx, bo)
	if err != nil {
		t.Fatal(err)
	}

	importpath, err := bi.QualifyImport(path.Join(cwd, "testdata"))
	if err != nil {
		t.Fatal(err)
	}

	br, err := bi.Build(ctx, importpath)
	if err != nil {
		t.Fatal(err)
	}

	po := &kooptions.PublishOptions{
		DockerRepo:       repo,
		InsecureRegistry: true,
		Tags:             []string{"latest"},
		ImageNamer:       md5Namer,
		Push:             true,
	}

	p, err := commands.NewPublisher(po)
	if err != nil {
		t.Fatal(err)
	}

	imageRef, err := p.Publish(ctx, br, importpath)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Published: %s", imageRef.Name())
	assert.NoError(t, err)
	privateKey := path.Join(cwd, "testdata", "test.key")
	err = cosignSign(ctx, Args{
		Key:         privateKey,
		KeyPassword: "password",
		Images:      []string{imageRef.Name()},
		Insecure:    true,
	})
	assert.NoError(t, err)

	//Verify
	publicKey := path.Join(cwd, "testdata", "test.pub")
	err = cosignVerify(ctx, Args{
		Key:      publicKey,
		Images:   []string{imageRef.Name()},
		Insecure: true,
	})
	assert.NoError(t, err)
}

func md5Namer(base, importpath string) string {
	hasher := md5.New() // nolint: gosec // No strong cryptography needed.
	hasher.Write([]byte(importpath))
	return path.Join(base, path.Base(importpath)+"-"+hex.EncodeToString(hasher.Sum(nil)))
}

func registryServerWithImage(namespace string) (*httptest.Server, error) {
	nopLog := log.New(ioutil.Discard, "", 0)
	r := registry.New(registry.Logger(nopLog))
	s := httptest.NewServer(r)
	imageName := fmt.Sprintf("%s/%s", s.Listener.Addr().String(), namespace)
	image, err := random.Image(1024, 1)
	if err != nil {
		return nil, fmt.Errorf("random.Image(): %w", err)
	}
	crane.Push(image, imageName)
	return s, nil
}
