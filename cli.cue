package main

import (
    "dagger.io/dagger"

    "universe.dagger.io/alpine"
    "universe.dagger.io/docker"
    "universe.dagger.io/bash"
)

dagger.#Plan & {
    client: filesystem: ".": read: exclude: [
        ".github",
    ]
    actions: {
        _source: client.filesystem["."].read.contents

        lint: {
            "cue": #Lint & {
                source: _source
            }
        }
    }
}

#Lint: {
    source: dagger.#FS

    docker.#Build & {
        steps: [
            alpine.#Build & {
                packages: bash: _
                packages: curl: _
                packages: git:  _
            },

            docker.#Copy & {
                contents: source
            },

            // Install CUE
            bash.#Run & {
                script: contents: #"""
                            export CUE_VERSION="v0.4.3"
                            export CUE_TARBALL="cue_${CUE_VERSION}_linux_amd64.tar.gz"
                            echo "Installing cue version $CUE_VERSION"
                            curl -L "https://github.com/cue-lang/cue/releases/download/${CUE_VERSION}/${CUE_TARBALL}" | tar zxf - -C /usr/local/bin
                            cue version
                    """#
            },
            // LINT
            bash.#Run & {
                workdir: "/cue"
                script: contents: #"""
                    cue vet sample.yaml check.cue
                    """#
            },
        ]
    }
}
