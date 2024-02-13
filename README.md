# shedock

Generate the most _ethically_ optimized `Dockerfile` 🐳 for your shell scripts!

> [!IMPORTANT]
> shedock is under active development and is changing everyday. Please consider `Watching` the repository for _releases_ to ge notified whenever a stable build is out!


## Features (a.k.a Goals) ✨

### 🧳 Minimal & Optimized

Shedock will generate the most ethically minimal `Dockerfile` for your shell script. This means.

- Automatically detecting and adding your script & shell dependencies to the final Docker image.
- Add only the stuff _absolutely_ required to run your script and get rid of unnecessary things whic ultimately results in a smaller image size.

### ℹ️ You own it

- Shedock will generate a heavily documented Dockerfile so that you can understand what's going on.
- This way shedock encourages users to learn and maintain the `Dockerfile` on their own (we are planning to cover this bridge as well).

### 🧐 Insights

- We would be lying to ourselves if we say all edge cases are covered. shedock WILL fail in some weird cases. But we try to figure out what those scenarios are & generate tips for you. So that you can make the best decision until shedock becomes capable of fixing it.
- This also reduces the need to learn a new tool (except Docker). Everything needed will be present in the Dockerfile generated.


## Who is `shedock` built for?

- Authors, folks who want to distribute their shell-based apps, or bring a new life to them ☘️. Dockerizing your shell scripts makes them available to EVERYONE!
- Users, folks who don't like installing random shell scripts from the internet & want a nice controllable isolated environment for them i.e. containers.
- Folks who want to consider `Docker` as a delivery format for their shell scripts.

## When not to use `shedock`, or when not to write a `Dockerfile` for your script?

- If you are depending on the host machine's resources, like `notify-send`, `xdg-open`, or anything UI. In these cases, your scripts are deeply tied to the system you use every day, it's hard to replicate that in a containerized environment.
- These scripts are still cool, but you should consider shipping them via package managers.

## Inspiration

I got inspired by [my own article](https://bhupesh.me/publishing-my-first-ever-dockerfile-optimization-ugit/) which I wrote while writing a Dockerfile for [ugit (a shellscript based tool)](https://github.com/Bhupesh-V/ugit). I learned cool new stuff which then I realized can be materialized into this tool.


## Building

### Pre-requisties

1. Docker (running)
2. Go (>=1.21.3)
3. The Internet

To run on Mac, you need to set the following environment variables:

```bash
export DOCKER_HOST="unix:///Users/$USER/.docker/run/docker.sock"
export DOCKER_API_VERSION=1.43
```

## Shell Comptability Chart

List & status of shells supported by shedock:


|   Shell    | Comptability Status | Notes |
| :--------: | :-----------------: | ----- |
|    bash    |          ✅          |       |
|    zsh     |          ⚪️          |       |
|     sh     |          ⚪️          |       |
|    fish    |          ⚪️          |       |
|    oil     |          ⚪️          |       |
|    csh     |          ⚪️          |       |
|    tcsh    |          ⚪️          |       |
|    ksh     |          ⚪️          |       |
|    dash    |          ⚪️          |       |
|    mksh    |          ⚪️          |       |
|    oksh    |          ⚪️          |       |
|    osh     |          ⚪️          |       |
|    elv     |          ⚪️          |       |
| powershell |          ⚪️          |       |


- ✅ means that the shell is supported by shedock.
- ⚪️ means that it's a WIP.
- ❌ means that the shell is not supported by shedock.

## FAQs

<details>
  <summary><h3>It's <code>202N</code>, why are we still writing shell scripts?</h3></summary>
- `They are fun to write` + `They work`, deal with it.
</details>
<details>
  <summary><h3><code>Docker</code> is not a package manager, why are we using it to package shell scripts?</h3></summary>
- It's literally built to share your work across different systems. `brew` is not popular with Linux. Flatpaks, AppImages don't work on Mac, and the new Windows Terminal has its package manager now. How much time do you want to spend just packaging stuff compared to people utilising your work?
</details>
<details>
  <summary><h3>What about other tools like <code>docker-squash</code> & <code>docker-slim</code>?</h3></summary>

- They are great, they have a big community behind them actively building and fixing stuff. Give them a try before using shedock.
- shedock is built to educate devs, we want folks to know what exactly is required to run their script. Not hiding stuff behind some weird magic.
</details>
<details>
<summary><h3>Why not build this for all tech stacks, why only shell scripts?</h3></summary>

- A: because when you are building for everyone, you are building for no one.
- B: the author is biased towards writing and sharing shell scripts 🤓.
- C: the author doesn't have the mental energy to build & test it across 100s of tech stacks.
</details>
<details>
<summary><h3>I don't think you are building this right</h3></summary>

- Great, we have something in common 🙃. I am figuring out stuff on the go. If you think something can be improved, [start a new discussion](https://github.com/shedock/shedock/discussions) and leave me some helpful tips.
</details>
