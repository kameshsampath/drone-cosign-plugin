# Examples

The folder has few Drone pipeline examples that shows how to use the drone `cosign` plugin.

For the demonstration we will use the local key pair.

Run the following command to generate the key pair,

```shell
cosign generate-key-pair
```

The command will generate the key pair in the `$PWD` with names `cosign.key` and `cosign.pub`.

All the examples requires you to have a secret file `secrets.txt` with the following contents,

```shell
key_password: <the password that you gave while creating the key pair>
```

For this example we will push the images to the local container registry. Start a local container registry.

```shell
docker run -d  -p 5001:5000 --name=registry --restart always registry:2 
```

Once you create this file run the drone pipeline using the command,

```shell
drone exec --trusted --secret-file <your secret file>
```

The pipeline will build the image from the Dockefile, sign it using the `cosign.key` and push it to the repo `localhost:5001/example/hello-world` and finally verify the signature using the public key `cosign.pub`.
