# shedock

Generate the most ethically optimized & smallest `Dockerfile` for your shell scripts!

> [!IMPORTANT]
> shedock is under active development and is changing everyday. Please consider `Watching` the repository for _releases_.


## Features (Goals) ✨

### Minimal & Optimized

Shedock will generate the most ethically minimal `Dockerfile` for your shell script. This means.

- Automatically detecting and adding your script & shell dependencies to the final Docker image.
- Add only the stuff _absolutely_ required to run your script.
- This ultimately results in a smaller image size.

### You own it

- Shedock will generate a heavily documented Dockerfile, so that you can understand what's going on.
- This way shedock encourages users to learn and maintain the `Dockerfile` on their own (we are planning to cover this bridge as well).

### Insights

- We would be lying to ourselves if we say all edge cases our covered. shedock WILL fail on some weird cases.
- But we try to figure out what those scenarios are and generate tips for you. So that you can take the best decision until shedock becomes capable to fix it.

### Easy to use

- Shedock doesn't have any unnecessary flags, just install & run `shedock /path/to/script.sh` and it will generate the Dockerfile for you.

## Who is `shedock` built for?

- Authors, folks who want to distribute their shell based apps, or bring a new life to them ☘️. Dockerizing your shellscripts make them available to EVERYONE!
- Users, folks who don't like installing random shell scripts from the internet & want a nice controllable isolated environment for them i.e. containers.
- Folks who want to consider `Docker` as a packaging format for their shell scripts.

## When not to use `shedock`, or when not to write a `Dockerfile` for your script?

- If you are depending on the host machine's resources, like `notify-send`, `xdg-open`, or anything UI. In these cases, your scripts are deeply tied to the system you use every day, it's hard to replicate that in a containerized environment.
- These scripts are still cool, but you should consider shipping them via package managers.

## Inspiration

I got inspired by [my own article]() which I wrote while writing a Dockerfile for [ugit (a shellscript based tool)](). I learned cool new stuff which I realized can be materialized into this tool.

To run on Mac, you need to set the following environment variables:

## Building

### Pre-requisties

1. Docker (running)
2. Go (>=1.21.3)
3. The Internet

```bash
export DOCKER_HOST="unix:///Users/$USER/.docker/run/docker.sock"
export DOCKER_API_VERSION=1.43
```

## FAQs

### Its `2024`, why are we still writing shell scripts?

- `They are fun to write` + `They work`, deal with it.

### What about other tools like `docker-squash` & `docker-slim`?

- They are great, they have a big community behind them actively building and fixing stuff. Give them a try before using shedock.
- shedock is not built having competition in mind, but rather, "education" we want devs to know what exactly is required to run their script. Not hiding stuff behind some weird magic.

### Why not built this for all tech stacks, why only shell-scripts?

- A: because when you are building for everyone, you are building for no-one.
- B: shedock's author is biased towards writing and sharing shell scripts.

### I don't think you are building this right

- Great, we have something in common 🙃. I am figuring stuff on the go. If you think something can be improved, [start a new discussion]() and leave me some helpful tips.
