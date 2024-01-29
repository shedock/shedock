# shedock

Generate the most optimized & smallest `Dockerfile` for your shell scripts!

## Features ✨

- **Optimized** - Shedock will generate the most ethically minimal image for your shell script. This means
  - Automatically detecting and adding your script & shell dependencies to the final Docker image. Yes this includes external commands like `curl`, `wget`, `git`, etc.
  - Build the image with only the stuff necessary to run your script.
- **Easy** - Shedock is easy to use, just run `shedock <script.sh>` and it will generate the `Dockerfile` for you.
- **Heavily Documented Dockerfile** - Shedock will generate a heavily documented `Dockerfile` for you, so that you can understand what's going on.

## Who is `shedock` built for?

- **Shell script developers** who want to distribute their scripts as a Docker image.
- **Shell script users** who want to run shell scripts in a Docker container.

## When not to use `shedock`, or when not to write a `Dockerfile` for your script?

- If you are depending on the host machine's resources, like `notify-send`, `xdg-open`, or anything UI.

## Its `2024`, why are we still writing shell scripts?

- They work, deal with it.

To run on a Mac, you need to set the following environment variables:

```bash
export DOCKER_HOST="unix:///Users/$USER/.docker/run/docker.sock"
export DOCKER_API_VERSION=1.43
```