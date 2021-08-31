all: install build
.PHONY: all

LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH}"
include Makefile.Inc

install:
	go mod download
	GO111MODULE=off go get golang.org/x/lint/golint
	#go run cmd/fetch_artifacts/fetch_artifacts.go master "" build
.PHONY: install

build:
	go build ./...
.PHONY: build

SVM_DIR=svm-dist
build-rust-svm:
	if [ -d svm-dist ]; then cd $(SVM_DIR) && git pull; else git clone git@github.com:spacemeshos/svm $(SVM_DIR); fi
	mkdir -p $(BIN_DIR) svm/wasm
	cd $(SVM_DIR) && cargo +nightly build  --release --package svm-runtime-ffi --features=default-cranelift,default-memory --no-default-features
	cd $(SVM_DIR) && cargo +nightly build  --release --package svm-cli --features=default-cranelift,default-memory --no-default-features
	cp $(SVM_DIR)/target/release/svm.h svm/
	cd $(SVM_DIR)/crates/runtime-ffi/tests/wasm/counter && sh ./build.sh
	cp $(SVM_DIR)/crates/runtime-ffi/tests/wasm/counter.wasm svm/wasm/
	cd $(SVM_DIR)/crates/runtime-ffi/tests/wasm/failure && sh ./build.sh
	cp $(SVM_DIR)/crates/runtime-ffi/tests/wasm/failure.wasm svm/wasm/
ifeq ($(platform),linux)
	cp $(SVM_DIR)/target/release/libsvm_runtime_ffi.so $(BIN_DIR)libsvm.so
	cp $(SVM_DIR)/target/release/svm-cli $(BIN_DIR)svm-cli
endif
	cd $(SVM_DIR)/crates/cli/examples/craft-deploy \
		&& cargo +nightly build --features=ffi,static-alloc,meta --no-default-features --release --target wasm32-unknown-unknown \
        && rm -f craft_deploy_example.wasm \
        && cp ./target/wasm32-unknown-unknown/release/svm_cli_craft_deploy_example.wasm ./craft_deploy_example.wasm \
        && ./../../../../target/release/svm-cli craft-deploy --smwasm craft_deploy_example.wasm --meta Template-meta.json --output craft_deploy_example.bin
	cp $(SVM_DIR)/crates/cli/examples/craft-deploy/craft_deploy_example.bin svm/test_assets/

.PHONY: build-rust-svm

test-tidy:
	# Working directory must be clean, or this test would be destructive
	git diff --quiet || (echo "\033[0;31mWorking directory not clean!\033[0m" && git --no-pager diff && exit 1)
	# We expect `go mod tidy` not to change anything, the test should fail otherwise
	make tidy
	git diff --exit-code || (git --no-pager diff && git checkout . && exit 1)
.PHONY: test-tidy

test-fmt:
	git diff --quiet || (echo "\033[0;31mWorking directory not clean!\033[0m" && git --no-pager diff && exit 1)
	# We expect `go fmt` not to change anything, the test should fail otherwise
	go fmt ./...
	git diff --exit-code || (git --no-pager diff && git checkout . && exit 1)
.PHONY: test-fmt

tidy:
	go mod tidy
.PHONY: tidy

lint:
	golint --set_exit_status ./...
	go vet ./...
.PHONY: lint

ci-test: build-rust-svm install build test-all
