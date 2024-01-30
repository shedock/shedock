# shedock

Generate the most optimized & smallest `Dockerfile` for your shell scripts!

## Features ✨

### Optimized

Shedock will generate the most ethically minimal image for your shell script. This means.

- Automatically detecting and adding your script & shell dependencies to the final Docker image. Yes this includes external commands like `curl`, `wget`, `git`, etc.
- Build the image with only the stuff necessary to run your script.

### You own it

- Shedock will generate a heavily documented `Dockerfile` for you, so that you can understand what's going on.
- This way shedock encourages users to learn and maintain the `Dockerfile` on their own (for now).

### Easy

Shedock is easy to use, just run `shedock <script.sh>` and it will generate the `Dockerfile` for you.

## Who is `shedock` built for?

- Authors, folks who want to distribute their shell based apps, or bring a new life to them ☘️. Dockerizing your shellscripts make them available to EVERYONE!
- Users, folks who don't like installing random shell scripts from the internet & want a nice controllable isolated environment for them i.e. containers.

## When not to use `shedock`, or when not to write a `Dockerfile` for your script?

- If you are depending on the host machine's resources, like `notify-send`, `xdg-open`, or anything UI. In these cases your scripts are deeply tied to the system you use everyday, its hard to replicate that in a containerized env.

## Its `2024`, why are we still writing shell scripts?

- They are fun to write + They work, deal with it.

## Inspiration

I got inspired by [my own article]() which I wrote while writing a Dockerfile for [ugit (a shellscript based tool)](). I learned cool new stuff which I realised can be materialzed into this tool.

To run on Mac, you need to set the following environment variables:

## Installation

```bash
export DOCKER_HOST="unix:///Users/$USER/.docker/run/docker.sock"
export DOCKER_API_VERSION=1.43
```