.PHONY: dist

# Compile the binaries for all targets.
build: \
	build-x86_64-unknown-linux-musl \
	build-aarch64-unknown-linux-musl \
	build-armv7-unknown-linux-musleabihf

build-x86_64-unknown-linux-musl:
	cross build --target x86_64-unknown-linux-musl --release

build-aarch64-unknown-linux-musl:
	cross build --target aarch64-unknown-linux-musl --release

build-armv7-unknown-linux-musleabihf:
	cross build --target armv7-unknown-linux-musleabihf --release

# Build distributable binaries for all targets.
dist: \
	dist-x86_64-unknown-linux-musl \
	dist-aarch64-unknown-linux-musl \
	dist-armv7-unknown-linux-musleabihf

dist-x86_64-unknown-linux-musl: build-x86_64-unknown-linux-musl package-x86_64-unknown-linux-musl

dist-aarch64-unknown-linux-musl: build-aarch64-unknown-linux-musl package-aarch64-unknown-linux-musl

dist-armv7-unknown-linux-musleabihf: build-armv7-unknown-linux-musleabihf package-armv7-unknown-linux-musleabihf

# Package the compiled binaries
package-x86_64-unknown-linux-musl:
	$(eval PKG_VERSION := $(shell cargo metadata --no-deps --format-version 1 | jq -r '.packages[0].version'))
	mkdir -p dist

	# .tar.gz
	tar -czvf dist/chirpstack-packet-multiplexer_$(PKG_VERSION)_amd64.tar.gz -C target/x86_64-unknown-linux-musl/release chirpstack-packet-multiplexer

	# .deb
	cargo deb --target x86_64-unknown-linux-musl --no-build --no-strip
	cp target/x86_64-unknown-linux-musl/debian/*.deb ./dist

package-aarch64-unknown-linux-musl:
	$(eval PKG_VERSION := $(shell cargo metadata --no-deps --format-version 1 | jq -r '.packages[0].version'))
	mkdir -p dist

	# .tar.gz
	tar -czvf dist/chirpstack-packet-multiplexer_$(PKG_VERSION)_arm64.tar.gz -C target/aarch64-unknown-linux-musl/release chirpstack-packet-multiplexer

	# .deb
	cargo deb --target aarch64-unknown-linux-musl --no-build --no-strip
	cp target/aarch64-unknown-linux-musl/debian/*.deb ./dist

package-armv7-unknown-linux-musleabihf:
	$(eval PKG_VERSION := $(shell cargo metadata --no-deps --format-version 1 | jq -r '.packages[0].version'))
	mkdir -p dist

	# .tar.gz
	tar -czvf dist/chirpstack-packet-multiplexer_$(PKG_VERSION)_armv7hf.tar.gz -C target/armv7-unknown-linux-musleabihf/release chirpstack-packet-multiplexer

	# .deb
	cargo deb --target armv7-unknown-linux-musleabihf --no-build --no-strip
	cp target/armv7-unknown-linux-musleabihf/debian/*.deb ./dist

# Update the version.
version:
	test -n "$(VERSION)"
	sed -i 's/^  version.*/  version = "$(VERSION)"/g' ./Cargo.toml
	make test
	git add .
	git commit -v -m "Bump version to $(VERSION)"
	git tag -a v$(VERSION) -m "v$(VERSION)"

# Cleanup dist.
clean:
	cargo clean
	rm -rf dist

# Run tests
test:
	cargo clippy --no-deps
	cargo test

# Enter the devshell.
devshell:
	nix-shell
