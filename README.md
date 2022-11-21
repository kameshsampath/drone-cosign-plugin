# Drone cosign Plugin

A [Drone](https://drone.io) plugin sign the container images using [sigstore](https://www.sigstore.dev/).

>**IMPORTANT:** The plugin is under development and currently does not support all the options provided by `cosign`.

## Usage

The following settings changes this plugin's behavior.

* `key`: The google cloud service account JSON.
* `verify`: The google cloud project to use.
* `key_password`: The google cloud region to set as default region. If provided the Artifact Registry for Docker will be enabled in these regions.
* `repo`: An array of Google Cloud Artifact registry locations to configure for docker authentication.
* `insecure`: An array of Google Cloud Artifact registry locations to configure for docker authentication.
* `dry_run`: An array of Google Cloud Artifact registry locations to configure for docker authentication.

### Sign

```yaml
kind: pipeline
type: docker
name: default

steps:
# build docker image 
- name: sign
  image: kameshsampath/drone-cosign
  pull: never
  settings:
      # path relative to sources or load it from secret
      key: cosign.key
      key_password: 
        from_secret: key_password
      repo: mycontainer-registry/my-image
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
      repo: mycontainer-registry/my-image
```

Please check the examples folder for `.drone.yml` with other settings.

## Building

Run the following command to build and push the image manually

```text
./scripts/build.sh
```

## Testing

```shell
docker run --rm \
  -e PLUGIN_KEY=$PLUGIN_COSIGN_KEY \
  -e PLUGIN_REPO="foo/example" \
  -e PLUGIN_KEY_PASSWORD=$PLUGIN_COSIGN_KEY_PASSWORD \
  kameshsampath/drone-cosign
```
