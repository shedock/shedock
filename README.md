# shedock

Generate the most optimized & smallest `Dockerfile` for your shell scripts!

To run on a Mac, you need to set the following environment variables:

```bash
export DOCKER_HOST="unix:///Users/$USER/.docker/run/docker.sock"
export DOCKER_API_VERSION=1.43
```

## Flow

1. ~~Find all the commands in the shell script.~~
2. ~~Find the base shell type (bash, zsh, sh, etc.)~~.
3. ~~Find builtins for the base shell type~~.
4. Filter out the builtins from the commands.
   1. These builtins will be installed with the base shell type.
   2. The remaining commands will be searched, and or installed with the Alpine package manager.
5. Start an alpine container
   1. Update package manager cache `apk update`.
   2. Find the dependencies for base shell type (shared objects and location to builtins) `apk info -a <shell>`.
   3. Install the base shell type `apk add <shell>`.
   4. Search the remaining commands in the Alpine package manager.
      1. If found, install the package.
      2. Find the dependencies for the package.
6. Label commands not found on apk as "user defined".
7. Generate the Dockerfile.

```
docker run --rm -it alpine:3.18 /bin/sh -c "apk update && apk info -R wget"
fetch https://dl-cdn.alpinelinux.org/alpine/v3.18/main/aarch64/APKINDEX.tar.gz
fetch https://dl-cdn.alpinelinux.org/alpine/v3.18/community/aarch64/APKINDEX.tar.gz
v3.18.5-166-g7f9be8f79f0 [https://dl-cdn.alpinelinux.org/alpine/v3.18/main]
v3.18.5-166-g7f9be8f79f0 [https://dl-cdn.alpinelinux.org/alpine/v3.18/community]
OK: 19947 distinct packages available
wget-1.21.4-r0 depends on:
so:libc.musl-aarch64.so.1
so:libcrypto.so.3
so:libidn2.so.0
so:libpcre2-8.so.0
so:libssl.so.3
so:libz.so.1

```

## Ideas

- Dockerize Shebang Movement
- Use `bash -c "compgen -c" > output.txt` to find all shell builtins.
- Load all utils from GNU https://www.gnu.org/manual/manual.html

## Shells

- https://pkgs.alpinelinux.org/package/edge/main/armhf/mksh
- https://pkgs.alpinelinux.org/package/edge/main/armhf/dash
- https://github.com/ibara/oksh
- https://pkgs.alpinelinux.org/package/v3.18/community/x86_64/powershell
- https://elv.sh/ (not on alpine main repo, only on edge)

## Important Reads

- https://stackoverflow.com/a/66974607
- GNU says that findutils contain `find`, `locate`, `updatedb`, `xargs`. But alpine doesn't have `locate` and `updatedb` in findutils. So, we need to install `mlocate` to get `locate` and `updatedb`.
- Brew on alpine doesn't work because of glibc.
  - https://github.com/Homebrew/brew/issues/8130
  - https://github.com/sgerrand/alpine-pkg-glibc/issues/178
  - https://github.com/prantlf/docker-alpine-glibc
  - https://stackoverflow.com/questions/37818831/is-there-a-best-practice-on-setting-up-glibc-on-docker-alpine-linux-base-image

```
coreutils-9.3-r2 provides:
cmd:[=9.3-r2
cmd:b2sum=9.3-r2
cmd:base32=9.3-r2
cmd:base64=9.3-r2
cmd:basename=9.3-r2
cmd:basenc=9.3-r2
cmd:cat=9.3-r2
cmd:chcon=9.3-r2
cmd:chgrp=9.3-r2
cmd:chmod=9.3-r2
cmd:chown=9.3-r2
cmd:chroot=9.3-r2
cmd:cksum=9.3-r2
cmd:comm=9.3-r2
cmd:coreutils=9.3-r2
cmd:cp=9.3-r2
cmd:csplit=9.3-r2
cmd:cut=9.3-r2
cmd:date=9.3-r2
cmd:dd=9.3-r2
cmd:df=9.3-r2
cmd:dir=9.3-r2
cmd:dircolors=9.3-r2
cmd:dirname=9.3-r2
cmd:du=9.3-r2
cmd:echo=9.3-r2
cmd:env=9.3-r2
cmd:expand=9.3-r2
cmd:expr=9.3-r2
cmd:factor=9.3-r2
cmd:false=9.3-r2
cmd:fmt=9.3-r2
cmd:fold=9.3-r2
cmd:head=9.3-r2
cmd:hostid=9.3-r2
cmd:id=9.3-r2
cmd:install=9.3-r2
cmd:join=9.3-r2
cmd:link=9.3-r2
cmd:ln=9.3-r2
cmd:logname=9.3-r2
cmd:ls=9.3-r2
cmd:md5sum=9.3-r2
cmd:mkdir=9.3-r2
cmd:mkfifo=9.3-r2
cmd:mknod=9.3-r2
cmd:mktemp=9.3-r2
cmd:mv=9.3-r2
cmd:nice=9.3-r2
cmd:nl=9.3-r2
cmd:nohup=9.3-r2
cmd:nproc=9.3-r2
cmd:numfmt=9.3-r2
cmd:od=9.3-r2
cmd:paste=9.3-r2
cmd:pathchk=9.3-r2
cmd:pinky=9.3-r2
cmd:pr=9.3-r2
cmd:printenv=9.3-r2
cmd:printf=9.3-r2
cmd:ptx=9.3-r2
cmd:pwd=9.3-r2
cmd:readlink=9.3-r2
cmd:realpath=9.3-r2
cmd:rm=9.3-r2
cmd:rmdir=9.3-r2
cmd:runcon=9.3-r2
cmd:seq=9.3-r2
cmd:sha1sum=9.3-r2
cmd:sha224sum=9.3-r2
cmd:sha256sum=9.3-r2
cmd:sha384sum=9.3-r2
cmd:sha512sum=9.3-r2
cmd:shred=9.3-r2
cmd:shuf=9.3-r2
cmd:sleep=9.3-r2
cmd:sort=9.3-r2
cmd:split=9.3-r2
cmd:stat=9.3-r2
cmd:stdbuf=9.3-r2
cmd:stty=9.3-r2
cmd:sum=9.3-r2
cmd:sync=9.3-r2
cmd:tac=9.3-r2
cmd:tail=9.3-r2
cmd:tee=9.3-r2
cmd:test=9.3-r2
cmd:timeout=9.3-r2
cmd:touch=9.3-r2
cmd:tr=9.3-r2
cmd:true=9.3-r2
cmd:truncate=9.3-r2
cmd:tsort=9.3-r2
cmd:tty=9.3-r2
cmd:uname=9.3-r2
cmd:unexpand=9.3-r2
cmd:uniq=9.3-r2
cmd:unlink=9.3-r2
cmd:users=9.3-r2
cmd:vdir=9.3-r2
cmd:wc=9.3-r2
cmd:who=9.3-r2
cmd:whoami=9.3-r2
cmd:yes=9.3-r2
```


A few questions:

- Is Docker running?
- Are you running Docker locally or connecting to a remote Daemon?
- If locally: do you have the /var/run/docker.sock file?
- Do you have the DOCKER_HOST environment variable set?


## Projects to use for testing

- https://github.com/KevCui/animepahe-dl

https://github.com/search?q=language%3AShell+&type=repositories