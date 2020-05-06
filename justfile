# Build the runtime shared library for this specific system.
build-dep:
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

			# Change the dynamic link to a relative one, for breaking the
			# dependency with the previous file path. The relative link
			# will be adjusted by cgo.
			install_name_tool -id "@rpath/${shared_library}" svm/${shared_library}
			;;
		"windows")
			echo "{{os()}}: platform not supported yet due to upstream dependency (single-pass compiler)"
			exit 1
			;;
		*)
			echo "{{os()}}: platform not supported yet"
			exit 1
	esac

# Generate cgo debug objects.
debug-cgo:
	cd svm && go tool cgo bridge.go && cd _obj && ls -d "$PWD/"*

# Run all the tests.
test:
	GODEBUG=cgocheck=2 go test ./... -v
