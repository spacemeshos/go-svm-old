# Change the dynamic link to a relative one,
# for breaking the dependency with the previous file path.
# The relative path will be adjusted by cgo.
fix-path:
	#!/usr/bin/env bash
	set -euo pipefail

	case "{{os()}}" in
			"macos")
				shared_library=libsvm_runtime_c_api.dylib
				install_name_tool -id "@rpath/${shared_library}" svm/${shared_library}
				echo "{{os()}}: lib path fixed to @rpath/${shared_library}"
				;;
			"windows")
				echo "{{os()}}: no fix is required"
				exit 1
				;;
			*)
				echo "{{os()}}: no fix is required"
				exit 1
		esac

# Fetch pre-compiled libs for all platforms from github.com/spacemeshos/svm CI.
fetch-artifacts branch token:
	#!/usr/bin/env bash
	set -euo pipefail

	dest=$(pwd)/svm
	pushd cmd/fetch_artifacts
	go build && ./fetch_artifacts -branch={{branch}} -token={{token}} -dest=$dest
	popd

# Re-build SVM on your platform.
build-svm:
	#!/usr/bin/env bash
	set -euo pipefail

	pushd svm-dep
	cargo +nightly build --release
	popd

	rm -f svm/svm.h
	cp svm-dep/target/release/svm.h svm/svm.h

	case "{{os()}}" in
		"macos")
			shared_library_path=$( ls -t svm-dep/target/release/deps/libsvm_runtime_c_api*.dylib | head -n 1 )
			shared_library=libsvm_runtime_c_api.dylib

			rm -f svm/${shared_library}
			cp ${shared_library_path} svm/${shared_library}
			;;
		"windows")
			echo "{{os()}}: local build not supported yet"
			exit 1
			;;
		*)
			echo "{{os()}}: local build not supported yet"
			exit 1
	esac

# Run all the tests.
test:
    just example
    GODEBUG=cgocheck=2 go test ./... -v

# Run the example.
example:
	#!/usr/bin/env bash
	set -euo pipefail

	pushd examples/counter/wasm
	./build.sh
	popd

	pushd examples/counter
	go build && ./counter
	popd

# Generate cgo debug objects.
debug-cgo:
	cd svm && go tool cgo bridge.go && cd _obj && ls -d "$PWD/"*
