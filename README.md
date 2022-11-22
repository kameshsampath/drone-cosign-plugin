# Drone cosign Plugin

A [Drone](https://drone.io) plugin sign the container images using [sigstore](https://www.sigstore.dev/).

>**IMPORTANT:**

- The plugin is under development and currently does not support all the options provided by `cosign`.
- Current version only supports keyful signing and verification using keys

## Usage

The following settings changes this plugin's behavior.

- `key`: The private key file that will be used to sign the image or public key file that will be used to verify the image
- `key_password`: The password to use with the key. If not provided will lookup an environment variable named `COSIGN_PASSWORD`.
- `verify`: Flag to indicate to **sign** or **verify**.
- `key_password`: The google cloud region to set as default region. If provided the Artifact Registry for Docker will be enabled in these regions.
- `images`: An array images that needs to be signed.
- `insecure`: Whether an insecure registry is used. **Should only be used for testing/development**
- `dry_run`: Whether to upload the signature to the repository.
- `check_claims`: Flag to indicate to check the claims. Defaults `true`.

### Sign

```yaml
kind: pipeline
type: docker
name: default

steps:
- name: sign
  image: kameshsampath/drone-cosign
  pull: never
  settings:
      # path relative to sources or load it from secret
      key: cosign.key
      key_password: 
        from_secret: key_password
      images: 
        - mycontainer-registry/my-image:latest@sha256:fc48fd8b997337537ddd4cf954931bddbcf28467c326645d87f83b0e0d46f4f9
```

### Verify

```yaml
kind: pipeline
type: docker
name: default

steps:
# build docker image 
- name: verify
  image: kameshsampath/drone-cosign
  pull: never
  settings:
      verify: true
      # path relative to sources or load it from secret
      key: cosign.pub
      images: 
        - mycontainer-registry/my-image:latest@sha256:fc48fd8b997337537ddd4cf954931bddbcf28467c326645d87f83b0e0d46f4f9
```

Please check the examples folder for `.drone.yml` with other settings.

## Building

Run the following command to build and push the image manually

```text
drone exec
```

## Testing

```shell
docker run --rm \
  -e PLUGIN_KEY=$PLUGIN_COSIGN_KEY \
  -e PLUGIN_REPO="foo/example" \
  -e PLUGIN_KEY_PASSWORD=$PLUGIN_COSIGN_KEY_PASSWORD \
  kameshsampath/drone-cosign
```
