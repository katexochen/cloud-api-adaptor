.DEFAULT_GOAL := all
PHONY: all
all: fedora-binaries-builder binaries image

PHONY: debug
debug: fedora-binaries-builder binaries image-debug

PHONY: fedora-binaries-builder
fedora-binaries-builder:
	@echo "Building fedora-binaries-builder image..."
	docker buildx build \
		-t fedora-binaries-builder \
		--load \
		- < ../podvm/Dockerfile.podvm_builder.fedora

PHONY: binaries
binaries:
	docker buildx use default
	@echo "Building binaries..."
	rm -rf ./resources/binaries-tree
	docker buildx build \
		--build-arg BUILDER_IMG=fedora-binaries-builder \
		-o type=local,dest="./resources/binaries-tree" \
		- < ../podvm/Dockerfile.podvm_binaries.fedora

define build-image
docker buildx use default
@echo "Building mkosi-builder..."
docker buildx build \
	-t mkosi-builder \
	--load \
	--no-cache \
	.
rm -rf ./build
mkdir -p build
@echo "Building $(1) image..."
docker run \
	--rm \
	--privileged \
	-v $(PWD)/build:/mkosi/build \
	mkosi-builder \
	mkosi \
	--environment=VARIANT_ID=$(1)
endef

PHONY: image
image:
	@echo "Enabling production preset..."
	rm -rf resources/buildDebugImage
	$(call build-image,production)


PHONY: image-debug
image-debug:
	@echo "Enabling debug preset..."
	touch resources/buildDebugImage
	$(call build-image,debug)
