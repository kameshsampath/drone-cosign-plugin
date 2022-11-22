package plugin

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/sigstore/cosign/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/cmd/cosign/cli/sign"
	"github.com/sigstore/cosign/cmd/cosign/cli/verify"
)

func cosignSign(_ context.Context, args Args) error {
	annotations := make(map[string]interface{}, len(args.Annotations))
	for _, a := range args.Annotations {
		kv := strings.Split(a, "=")
		if len(kv) > 0 {
			annotations[kv[0]] = kv[1]
		}
	}
	ro := &options.RootOptions{
		Timeout: 10 * time.Second,
	}
	ko := options.KeyOpts{
		KeyRef: args.Key,
		PassFunc: func(b bool) ([]byte, error) {
			pw, ok := os.LookupEnv("COSIGN_PASSWORD")
			if ok {
				return []byte(pw), nil
			}
			return []byte(args.KeyPassword), nil
		},
	}
	regOpts := options.RegistryOptions{
		AllowInsecure: args.Insecure,
	}

	return sign.SignCmd(ro, ko, regOpts, annotations, args.Images, "", "", !args.DryRun, "", "", "", false, false, "", false)
}

func cosignVerify(ctx context.Context, args Args) error {
	regOpts := options.RegistryOptions{
		AllowInsecure: args.Insecure,
	}

	if !args.CheckClaims {
		args.CheckClaims = true
	}

	//TODO(kamesh) add support for Signatures
	vc := &verify.VerifyCommand{
		RegistryOptions: regOpts,
		KeyRef:          args.Key,
		CheckClaims:     args.CheckClaims,
	}

	return vc.Exec(ctx, args.Images)
}
