# shedock

Generate the most optimized & smallest `Dockerfile` for your shell scripts!

> [!IMPORTANT]
> shedock is under active development and is changing everyday. Please consider subscribing for updates by clicking on the `Watch` the repo for _releases_.


## Features (Goals) ✨

### Optimized

Shedock will generate the most ethically minimal `Dockerfile` for your shell script. This means.

- Automatically detecting and adding your script & shell dependencies to the final Docker image. Yes this includes external commands like `curl`, `wget`, `git`, etc.
- Build the image with only the stuff necessary to run your script.
- This ultimately results in a smaller image size.

### You own it

- Shedock will generate a heavily documented `Dockerfile` for you, so that you can understand what's going on.
- This way shedock encourages users to learn and maintain the `Dockerfile` on their own (for now).

### Easy to use

- Shedock doesn't have any unnecessary flags, just install & run `shedock /path/to/script.sh` and it will generate the Dockerfile for you.

## Who is `shedock` built for?

- Authors, folks who want to distribute their shell based apps, or bring a new life to them ☘️. Dockerizing your shellscripts make them available to EVERYONE!
- Users, folks who don't like installing random shell scripts from the internet & want a nice controllable isolated environment for them i.e. containers.

## When not to use `shedock`, or when not to write a `Dockerfile` for your script?

- If you are depending on the host machine's resources, like `notify-send`, `xdg-open`, or anything UI. In these cases, your scripts are deeply tied to the system you use every day, it's hard to replicate that in a containerized environment.

## Its `2024`, why are we still writing shell scripts?

- They are fun to write + They work, deal with it.

## Inspiration

I got inspired by [my own article]() which I wrote while writing a Dockerfile for [ugit (a shellscript based tool)](). I learned cool new stuff which I realized can be materialized into this tool.

To run on Mac, you need to set the following environment variables:

## Installation

```bash
export DOCKER_HOST="unix:///Users/$USER/.docker/run/docker.sock"
export DOCKER_API_VERSION=1.43
```