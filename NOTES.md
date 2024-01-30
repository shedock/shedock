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
- https://elv.sh/ (not on alpine main repo, only on edge)

## Gotchas to tackle in the future

- Brew on alpine doesn't work because of glibc.
  - https://github.com/Homebrew/brew/issues/8130
  - https://github.com/sgerrand/alpine-pkg-glibc/issues/178
  - https://github.com/prantlf/docker-alpine-glibc
  - https://stackoverflow.com/questions/37818831/is-there-a-best-practice-on-setting-up-glibc-on-docker-alpine-linux-base-image
- Understand the difference b/w `glibc` and `musl`
  - https://stackoverflow.com/a/66974607
- `notify-send` wont work on containers
  - https://github.com/mikaelbr/node-notifier/issues/200
- Automatically detect color codes in the shell script and add `TERM=xterm-256color` to the Dockerfile.
  - https://stackoverflow.com/a/20983251
  - https://chadaustin.me/2024/01/truecolor-terminal-emacs/

## Projects to use for testing

- https://github.com/KevCui/animepahe-dl
- https://github.com/jarun/ddgr
- https://github.com/HarshitJoshi9152/libgen
- https://github.com/kamranahmedse/git-standup
- Find more on https://github.com/search?q=language%3AShell+&type=repositories


## Blog Ideas

- The need for dockerizing your shell scripts.
- The awesomeness of Alpine images.
- Static v/s Dynamic Libraries, Linking etc.
- Cool things you can do via shell scripts.
- Best practices for shell scripts.
