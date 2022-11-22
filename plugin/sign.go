package plugin

import (
	"os"
	"time"

	"github.com/sigstore/cosign/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/cmd/cosign/cli/sign"
)

func cosignImage(keyRef string, keyPassword string, imgs []string, insecure bool, annotations map[string]interface{}) error {
	ro := &options.RootOptions{
		Timeout: 10 * time.Second,
	}
	ko := options.KeyOpts{
		KeyRef: keyRef,
		PassFunc: func(b bool) ([]byte, error) {
			pw, ok := os.LookupEnv("COSIGN_PASSWORD")
			if ok {
				return []byte(pw), nil
			}
			return []byte(keyPassword), nil
		},
	}
	regOpts := options.RegistryOptions{
		AllowInsecure: insecure,
	}

	//sign
	err := sign.SignCmd(ro, ko, regOpts, annotations, imgs, "", "", true, "", "", "", false, false, "", false)

	return err
}
