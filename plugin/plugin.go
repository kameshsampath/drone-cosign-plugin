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

	// Annotations the annotations to add to Signatures
	Annotations []string `envconfig:"PLUGIN_ANNOTATIONS"`
	// DryRun disables upload of the signatures
	DryRun bool `envconfig:"PLUGIN_DRY_RUN"`
	// Key is key that will be used to sign the image or public key to verify the image
	Key string `envconfig:"PLUGIN_KEY"`
	// KeyPassword is the password that will be used with the Key
	KeyPassword string `envconfig:"PLUGIN_KEY_PASSWORD"`
	// Images is the image repo/digest that will be signed
	Images []string `envconfig:"PLUGIN_IMAGES"`
	// Insecure marks the repo as insecure. Should be used only for testing
	Insecure bool `envconfig:"PLUGIN_INSECURE"`
	// Verify does the verification of the signature.
	Verify bool `envconfig:"PLUGIN_VERIFY"`
	// CheckClaims verifies the claims when doing verification. Defaults to true
	CheckClaims bool `envconfig:"PLUGIN_CHECK_CLAIMS"`
	// Signature
	Signature string `envconfig:"PLUGIN_VERIFY_SIGNATURE"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	if args.Verify {
		return cosignVerify(ctx, args)
	}
	return cosignSign(ctx, args)
}
