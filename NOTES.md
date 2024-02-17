# Notes

## shedock image build flow

1. ~~Find all the commands in the shell script.~~
2. ~~Find the base shell type (bash, zsh, sh, etc.)~~.
3. ~~Find builtins for the base shell type~~.
4. ~~Filter out the shell builtins from the commands~~.
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


## Load most commonly used commands

- Load all utils from GNU https://www.gnu.org/manual/manual.html.
- GNU says that findutils contain `find`, `locate`, `updatedb`, `xargs`. But alpine doesn't have `locate` and `updatedb` in findutils. So, we need to install `mlocate` to get `locate` and `updatedb`.
- Find some utils from https://github.com/agarrharr/awesome-cli-apps
- https://github.com/tj/git-extras/tree/main/bin
- https://github.com/sharkdp/bat
- https://github.com/BurntSushi/ripgrep
- https://github.com/ogham/dog


## Shells to support in the future

- https://pkgs.alpinelinux.org/package/edge/main/armhf/mksh
- https://pkgs.alpinelinux.org/package/edge/main/armhf/dash
- https://github.com/ibara/oksh
- https://pkgs.alpinelinux.org/package/v3.18/community/x86_64/powershell
- https://github.com/nushell/nushell
- https://elv.sh/ (not on alpine main repo, only on edge)

## Gotchas to tackle in the future

- Brew on alpine doesn't work because of glibc.
  - https://github.com/Homebrew/brew/issues/8130
  - https://github.com/sgerrand/alpine-pkg-glibc/issues/178
  - https://github.com/prantlf/docker-alpine-glibc
  - https://stackoverflow.com/questions/37818831/is-there-a-best-practice-on-setting-up-glibc-on-docker-alpine-linux-base-image
- Understand the difference b/w `glibc` and `musl`
  - https://stackoverflow.com/a/66974607
  - https://wiki.musl-libc.org/
- `notify-send` wont work on containers
  - https://github.com/ku1ik/git-dude/blob/master/git-dude#L18-L30
  - https://github.com/mikaelbr/node-notifier/issues/200
- Automatically detect color codes in the shell script and add `TERM=xterm-256color` to the Dockerfile.
  - https://stackoverflow.com/a/20983251
  - https://chadaustin.me/2024/01/truecolor-terminal-emacs/
- Usage of `/dev/urandom`, `/dev/null` in shell scripts.

## Projects to use for testing

- https://github.com/KevCui/animepahe-dl
- https://github.com/HarshitJoshi9152/libgen
- https://github.com/kamranahmedse/git-standup
- https://github.com/bigH/git-fuzzy
- https://github.com/wfxr/forgit
- https://github.com/v1s1t0r1sh3r3/airgeddon
- https://github.com/fcambus/ansiweather (its already on alpine)
- Find more on https://github.com/search?q=language%3AShell+fork%3Afalse+stars%3A%3E1000&type=repositories&s=updated&o=desc


## Blog Ideas

- The need for dockerizing your shell scripts.
- The awesomeness of Alpine images.
- Static v/s Dynamic Libraries, Linking etc.
- Best practices for shell scripts.

## Reads

- Union File Systems
  - https://martinheinz.dev/blog/44
  - https://docs.docker.com/storage/storagedriver/overlayfs-driver/
  - https://docs.docker.com/storage/storagedriver/aufs-driver/
  - https://leftasexercise.com/2018/04/12/docker-internals-process-isolation-with-namespaces-and-cgroups/
- Busybox
  - https://www.busybox.net/
- Packaging
  - https://ramcq.net/2024/02/06/flathub-pros-and-cons-of-direct-uploads/
- CLI Guidelines
  - https://clig.dev/

## Checklist for use-cases while Dockerizing a shell script

- [ ] File operations
  - [ ] Read files and directories
  - [ ] Write files and directories
  - [ ] Watching files and directories
- [ ] Network operations
  - [ ] HTTP requests
  - [ ] DNS lookups
- [ ] Displaying TUI
- [ ] Scheduling work (cron jobs)
- [ ] Interacting with the system (userspace)
  - [ ] Running commands
  - [ ] Running commands in the background
  - [ ] Running commands in parallel
  - [ ] Setting environment variables
  - [ ] Reading environment variables


## Insights to provide

- Inconsistent usage of `echo` & `printf`.
- Unnecessary usage of HERE documents. Use multi-line `printf` instead.
- Incompatible usage of commands:
  - xdg-open, open, notify-send, etc.
- Unncessary usage of basename:
  - `$(basename "${BASH_SOURCE[0]}")` -> `<script_name>`
  - Reduce size by hardcoding the script name.

## Package Maintainers

- https://debian-handbook.info/browse/stable/sect.becoming-package-maintainer.html
- https://www.reddit.com/r/linux/comments/kmat5j/what_exactly_is_expected_of_a_package_maintainer/
- https://wiki.archlinux.org/title/Package_Maintainers
- https://unixsheikh.com/articles/the-heavy-responsibility-of-the-package-maintainer.html
- https://github.com/jubalh/awesome-package-maintainer

## Container Tooling

- https://buildah.io/