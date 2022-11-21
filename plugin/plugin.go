// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import "context"

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	// DryRun disables upload of the signatures
	DryRun bool `envconfig:"PLUGIN_DRY_RUN"`
	// Key is key that will be used to sign the image
	Key string `envconfig:"PLUGIN_KEY"`
	// KeyPassword is the password that will be used with the Key
	KeyPassword string `envconfig:"PLUGIN_KEY_PASSWORD"`
	// Repo is the image repo/digest that will be signed
	Repo string `envconfig:"PLUGIN_REPO"`
	// Insecure marks the repo as insecure. Should be used only for testing
	Insecure bool `envconfig:"PLUGIN_INSECURE"`
	// Verify does the verification of the signature.
	Verify bool `envconfig:"PLUGIN_VERIFY"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	// if !args.Verify {
	// 	return sign(args)
	// }
	return nil
}
