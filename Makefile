
help:
	@echo Please choose a valid target: build, test, release, final-release, clean, delete-deployments, wipe-everything

build:
	{ \
	  cd src/; \
	  go build -mod=vendor -o a.out main/main.go; \
	  rm a.out; \
	}

test:
	{ \
	  cd src/; \
	  go version; \
	  go test -v ./...; \
	}

release:
	bosh create-release --force --tarball $(PWD)/cpi

final-release:
	last_tag=$(git describe --abbrev=0)
	if grep "version: ${last_tag}" releases/libvirt-bosh-cpi/index.yml > /dev/null
	then
		echo "Nothing to do."
	else
		bosh create-release --final --version=${last_tag}
	fi

clean:
	[ -f cpi ] && rm cpi

delete-deployments:
	bosh --json deployments | jq -r '.Tables[].Rows[].name' | \
		xargs --verbose --max-args=1 --replace={} bosh --non-interactive --deployment {} delete-deployment

wipe-everything:
	@echo TODO
	# dump all libvirt VMs
	# clean all secrets
	# clean state.json etc
