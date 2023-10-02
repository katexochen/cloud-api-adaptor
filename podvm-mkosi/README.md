# Podvm images with mkosi

[mkosi](https://github.com/systemd/mkosi) builds a bootable OS image from scratch. This way, we have full control over every detail of the image, especially over the image format and boot process. On the long run, we will implement fully, bit-by-bit reproducible images with mkosi, and use measured boot and an immutable root FS to ensure the image integrity through remote attestation.

## Prerequisites

Currently, mksoi and other related tools are provided through a [Nix](https://nixos.org/) flake. Nix ensures all tools used in the build of the image are itself reproducible and pinned. mkosi requires a very recent systemd version, so using tools installed on the host is usually not possible.

The recommended way to install Nix on Linux systems is the [Determinate Systems nix-installer](https://github.com/DeterminateSystems/nix-installer).


## Building an image

### Building PodVM binaries and place them in the binaries-tree

```sh
docker buildx use default

docker buildx build \
    -t podvm-builder-fedora \
    --load \
    - < ../podvm/Dockerfile.podvm_builder.fedora

docker buildx build \
    --build-arg BUILDER_IMG=podvm-builder-fedora \
    -o type=local,dest="./resources/binaries-tree" \
    - < ../podvm/Dockerfile.podvm_binaries.fedora
```

### Activating the Nix development shell

This puts the dependencies from the Nix flake into your path.

```sh
nix develop ./..#podvm-mkosi
mkosi --version
```

You can exit the environment later using <kbd>Ctrl</kbd> + <kbd>D</kbd>

### Build the image

```
rm -rf ./build
mkosi
```

### Upload the image to the desired cloud provider

You can upload the image with the tool of your choice, but the recommended way is using [uplosi](https://github.com/edgelesssys/uplosi). Follow the uplosi readme to configure your upload for the desired cloud provider. Then run:

```sh
# Using -i to increment the image version after the upload.
uplosi -i build/system.raw
```

If you want to use the image with libvirt, run the following to convert to qcow2 format:

```sh
qemu-img convert -f raw -O qcow2 build/system.raw build/system.qcow2
```

## Custom image configuration

You can easily place additional files in `resources/binaries-tree` after it has been populated by the
binaries build step. Notice that systemd units need to be enabled in the presets and links in the tree
won't be copied into the image.

## Limitations

The following limitations apply to these images. Notice that the limitations are intentional to
reduce complexity of configuration and CI and shall not be seen as open to-dos.

- `DISABLE_CLOUD_CONFIG=false` is implied. The mkosi images are currently using
    cloud init/cloud config, and there is no possibility to disable it.
