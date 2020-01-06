# What is this?

This repository contains an example of using [bazel](http://bazel.io)
for frontend web development.

The goal is to provide a well commented, end to end example covering the full
development workflow: from starting a development server, to generating a
packaged, minified, optimized web root to be used in a backend application and
having good facilities for testing.

In terms of technology, the example is based on:

- [yarn](http://yarnpkg.com) for javascript package management.
- [react](http://reactjs.org) as the frontend framework.
- [rollup](http://rollupjs.org) for the packaging, mostly because it is well supported out of the box with bazel.
- [typescript](http://typescriptlang.org) as the frontend language.
- [sass](http://sass-lang.com) well, for styling.
- [golang](http://golang.org) for backend development, but the document provides enough pointers to use other languages easily.

with bazel, we used the packages:

- [rules_nodejs](https://bazelbuild.github.io/rules_nodejs/) for the frontend development rules.
- [rules_sass](https://github.com/bazelbuild/rules_sass/blob/master/README.md) for the SASS integration.
- [ibazel](https://github.com/bazelbuild/bazel-watcher) to watch changes in my files, and automatically rebuild targets.

further, we used:

- [material-ui](https://material-ui.com/) to get more graphic elements.

# Using this repository

Rather than just running `git clone` and starting editing the repository
to adapt it to your own need, we actually recommend you go through the
steps described here to create your own repository. It'll be easier for you to pick the pieces
you need, and adapt them for your goals.

Start from the [next section](#getting-started), and follow along section after section.

If you just want to try out the build system in this repository, install the dependencies
indicated in the [getting started section](#getting-started), and:

- run `bazel build frontend` to create a released version of the example frontend code,
  which will be output in `dist/bin/frontend/release.js`. Alternatively, you can use
  `bazel build frontend:release`.
- run `bazel run frontend/devserver` to start a test web server listening on :5432.

To run a test server and rebuild and reload the page content every time it is changed, you can use `ibazel run frontend/devserver` instead.

# Getting started

In terms of evnvironment, you need to make sure to:

- Have bazel installed. [Follow instructions here](https://docs.bazel.build/versions/master/install.html) if you have never done so.
  Check with `bazel --version` to verify that you installed it properly, and it is in your PATH. We used bazel `2.0.0` to write this document.
- Have a working version of yarn installed. [Follow instructions here if you don't](https://yarnpkg.com/lang/en/docs/install/).
  Verify that tools installed by yarn can be run out of the box. This generally means that the directory shown with `yarn global bin` is in your PATH (`echo $PATH` to verify).
  Check with `yarn --version` to see that the installation worked. We used yarn `1.21.1` to write this document.

# Basic frontend repository

First, I created a `bazel workspace` by running:

    yarn create @bazel --packageManager=yarn --typescript scratch

The parameters are easy to explain:

- `--packageManager=yarn` tells the @bazel creation command that we want to use yarn moving forward.
- `--typescript` indicates we want to use typescript
- `scratch` is the name we want to give to the workspace.

Running this command will:

- create a directory `scratch`
- in this directory, you will see:
  - a `BUILD.bazel` file - the equivalent of a Makefile, CMake, SCons, build.xml for ants.
  - a `WORKSPACE.bazel` file - the configuration of bazel for this directory and subdirectories, listing
    the dependencies of your project.
  - a `tsconfig.json` - the configuration file for the typescript compiler.
  - a `package.json` file - usual list of javascript dependencies of your project

The `WORKSPACE` file and `package.json` will have the name we picked, `scratch`, as the name
of the workspace and the name of the package. The name of the workspace is particularly important,
as it used to define the paths for a few things you will see later.

## BUILD.bazel file

If you open `BUILD.bazel` in your favorite editor, you will see a simple rule like:

    exports_files(["tsconfig.json"], visibility = ["//:__subpackages__"])

just saying that `tsconfig.json` is made available to the entire project.

## WORKSPACE.bazel file

If you open the `WORKSPACE.bazel` file, you will see a few a bunch of:

    http_archive(url = "http://...", name = "whatever"...)
    load("@whatever/rule/file.bzl", "rule_to_import")

    rule_to_import(...)

which is all fairly intuitive:

- `http_archive` downloads a dependency via http, and makes it available as `@name_that_was_specified`, `@whatever` in the example above.
- `load` opens a .bzl file downloaded with `http_archive`, and imports the rule specified, `rule_to_import`.
- `rule_to_import` invokes the rule.

In short, the rules in the created WORKSPACE.bazel file:

- install the tools needed by bazel to manage the javascript/typescript/yarn dependencies (eg, yarn, npm, ts compiler, ...).
- make yarn available, and ensure that yarn dependencies are installed before the build is performed.

# Testing the workspace

Enter the created directory, `cd scratch`, and run:

    bazel clean
    bazel run @nodejs//:yarn -- --version

The first command will clean and reset the workspace. **I always recommend running `bazel clean` or even `bazel clean --expunge`**
at first, and whenever you get into weird errors.

The second command will tell bazel to download yarn, run it, and show the version. `bazel` tries to be hermetic, and for its builds,
rather than use the system commands, it will use the **commands and dependencies downloaded from the WORKSPACE definition**.

For example, `yarn --version` on my system shows `1.21.1`, but `bazel run @nodejs//:yarn -- --version` shows version `1.19.1`,
whatever the owners of the nodejs rules last tested and qualified.

Whenever I touch my system, I use `yarn`. Whenever I touch the workspace, I prefer to use `bazel run @nodejs//:yarn`, although the difference should be minimal.

# Installing a few more tools

Next, I installed `ibazel`, with:

    bazel run @nodejs//:yarn -- add -D @bazel/ibazel
    yarn global add @bazel/ibazel

The first to download and track it as a dev dependency in our workspace,
the second to make sure I have a version of ibazel I can run from my shell.

Further, I've verified that the yarn installation directory is in my
path by ensuring that `$PATH` contains:

    yarn global bin
    /home/username/.nvm/versions/node/v12.14.0/bin

Next, I installed `react` and rollup:

    bazel run @nodejs//:yarn -- add react
    bazel run @nodejs//:yarn -- add react-dom

    bazel run @nodejs//:yarn -- yarn add --dev @bazel/rollup
    bazel run @nodejs//:yarn -- yarn add --dev rollup

Finally, I configured SASS as per [instructions here](https://github.com/bazelbuild/rules_sass/blob/master/README.md#setup),
which also explains how to compile sass libraries.

# React "Hello World"

To get a react hello world running, just look in the `frontend/` directory.

Start from `main.tsx`, and you will see a pretty standard react hello world example:

    import * as React from "react";
    import * as ReactDOM from "react-dom";

    ReactDOM.render(<h1>Hello World</h1>, document.getElementById("root"));

Next look at the `BUILD.bazel` file. You will basically see 3 targets:

- The definition of a library, called `app`, that contains the `main.tsx` file.
  Running `bazel build frontend:app` will cause the typescript compiler to parse
  all the dependencies, and output syntax errors and similar.
- The definition of a rollup bundle, called `release`. Running `bazel build frontend:release`
  will first build the app, invoking the typescript compiler, and then run rollup
  with the `rollup.config.js` file, and generate the `release.js` file.

If you further look in the `frontend/devserver` directory, you will see an
`index.html` file, just loading the dependencies and the app, and a `BUILD.bazel` file
defining a dependency on the previous `app` target.

Running `bazel run frontend/devserver` will result in a web server starting on
port `5432` and reachable via http://127.0.0.1:5432/ running the "Hello World" example.

### END

If you get strange errors related to yarn.lock file or packages.lock
or end up in situations where the behavior is inconsistent with what you
would expect, run:

    bazel clean
    bazel clean --expunge # only if bazel clean is not enough.

This will clean some of the internal bazel state, and start all the
next builds from scratch.

Followed instructions here to setup a ts_devserver:
https://bazelbuild.github.io/rules_nodejs/TypeScript.html#serving-typescript-for-development

From a directory outside my repository installed ibazel globally:

yarn add -D @bazel/ibazel
yarn global add @bazel/ibazel

Verified that the path used by npm is in my `$PATH`:

    $ yarn global bin
    /home/ccontavalli/.nvm/versions/node/v12.14.0/bin

bazel run @nodejs//:yarn -- add react

yarn add react
yarn add react-dom

yarn add --dev @bazel/rollup
yarn add --dev rollup

alternatives:
https://github.com/zenclabs/bazel-javascript
https://github.com/zenclabs/bazel-javascript
https://github.com/zenclabs/bazel-javascript
https://github.com/zenclabs/bazel-javascript
https://github.com/zenclabs/bazel-javascript
https://github.com/zenclabs/bazel-javascript
https://github.com/zenclabs/bazel-javascript
