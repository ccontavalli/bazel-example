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

This devserver is created with the [ts_devserver rule](https://bazelbuild.github.io/rules_nodejs/TypeScript.html#serving-typescript-for-development)
look at the `BUILD.bazel` file and the link above for more details, but it's pretty
trivial to setup.

# Creating a golang backend

## Setup

For the `bazel` setup for `golang` I started from [this blog post](https://brendanjryan.com/golang/bazel/2018/05/12/building-go-applications-with-bazel.html),
by Brendan Ryan (thanks).

In short:

1. Add the `golang` tooling to the `WORKSPACE.bazel` file. The blog post above provides
   a blob to cut and paste similar to the one below. DO NOT COPY IT. Instead, go on the
   [bazel golang repository](https://github.com/bazelbuild/rules_go/releases), and copy the recommended latest version.


       # DO NOT COPY THIS FROM HERE. Go to https://github.com/bazelbuild/rules_go/releases AND COPY THE LATEST VERSION.
       http_archive(
           name = "io_bazel_rules_go",
           urls = [
               "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/v0.20.3/rules_go-v0.20.3.tar.gz",
               "https://github.com/bazelbuild/rules_go/releases/download/v0.20.3/rules_go-v0.20.3.tar.gz",
           ],
           sha256 = "e88471aea3a3a4f19ec1310a55ba94772d087e9ce46e41ae38ecebe17935de7b",
       )

       load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

       go_rules_dependencies()
       go_register_toolchains()

2. Next, configure `gazelle`. It is a tool that allows to automatically maintain
   `BUILD.bazel` files based on the content of go files. You don't want to do this manually.
   Follow the [instructions here](https://github.com/bazelbuild/bazel-gazelle#running-gazelle-with-bazel), which
   will result in adding to `WORKSPACE.bazel`:


       http_archive(
           name = "bazel_gazelle",
           urls = [
               "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.19.1/bazel-gazelle-v0.19.1.tar.gz",
               "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.19.1/bazel-gazelle-v0.19.1.tar.gz",
           ],
           sha256 = "86c6d481b3f7aedc1d60c1c211c6f76da282ae197c3b3160f54bd3a8f847896f",
       )

       load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
       gazelle_dependencies()

3. Add `gazelle` as a buildable tool in the top level `BUILD.bazel` file. The important part here
   is to costumize the comment just before the `gazelle(...)` rule, as it tells gazelle how your
   own golang code will be referenced in import statements. To be clear, let's say in the `backend`
   directory you have a `lib` subdirectory. To import that code, you would normally have something like `import "github.com/whatever/backend/lib"`.
   The `#prefix` you have to provide gazelle - which can be any arbitrary string - is this `github.com/whatever`.
   Basically, it tells gazelle "really, anything my go code tries to import that starts with github.com/whatever comes
   from this same directory, no need to download it. For example, for this repository, we'd have:


       load("@bazel_gazelle//:def.bzl", "gazelle")
       # gazelle:prefix github.com/ccontavalli/bazel-example
       gazelle(name = "gazelle")

4. Finally, test the setup by running:

   bazel run //:gazelle

Given that gazelle is itself written in go, being able to run it means
that bazel was able to compile a golang binary, and run it.

**IMPORTANT**: as of January 5th, 2020, gazelle will fail to run if the
workspace file is named `WORKSPACE.bazel` (default used by `yarn create @bazel`)
instead of just `WORKSPACE`. See [the issue filed on github](https://github.com/bazelbuild/bazel-gazelle/issues/678).
The fix is simple: rename WORKSPACE.bazel in WORKSPACE, with `mv WORKSPACE.bazel WORKSPACE`.

## Writing and running some go code

It's really easy: just write go code as usual, run `gazelle`, and call it a day.

If you look in this repository:

1. I started by creating `backend/main.go`, like any other golang main.
2. I added an example library, `backend/lib/library.go`, exporting just the `MyLog` function.
3. Run gazelle, with `bazel run //:gazelle`. This created the `BUILD.bazel` files for me.
4. Minimally edited the BUILD.bazel files to have sensible names for the BUILD rules.
   For example, in `backend/lib/BUILD.bazel` I gave the generated rule name the name `lib`,
   while in `backend/BUILD.bazel`, I renamed the `go_binary` target to be named `backend`.
5. As you edit and update your code, just re-run `gazelle` as necessary to update your
   BUILD files. It will even automatically add tests for you.

With the minimal work here, I can now:

- run `bazel build backend` or `bazel run backend` to build and start my backend, which
  will start listening on [http://127.0.0.1:5433/](http://127.0.0.1:5433/).
- run `bazel build backend/lib` to just build the library.
- run `bazel test ...:all` and get the golang tests run as well as all the other tests.

## External dependencies in your go code

gazelle and bazel can be used in combination with [go modules](https://blog.golang.org/using-go-modules),
which is generally recommended, although they work slightly differently. To start using them:

1. `go mod init github.com/ccontavalli/bazel-example` will create a `go.mod` file.
2. Every time you need an external dependency, run `go get github.com/whatever` to install it.
   This will update the `go.mod` file and the `go.sum` file.
3. Run `bazel run //:gazelle -- update-repos -from_file=go.mod` to have the dependency tracked by bazel.
4. Keep coding as usual, repeat as necessary.

The difference here is that all the dependencies your code need will be downloaded and tracked by
bazel. Any external dependency not tracked will cause a build failure, easy to fix - regardless of
what is already installed on your machine.

Further, bazel will track the exact version and hash of each dependency, making your build
hermetic and easy to run in the cloud.

## Including assets in our backend

Let's say our golang backend needs to be able to serve our static frontend file, plus a number
of other assets. What I did in this repository:

1. Generate all the assets in the `backend/assets` directory, using `BUILD.bazel` rules.
2. Include those dependencies in the generated golang binary by using a `go_embed_data` rule.

[go_embed_data](https://github.com/bazelbuild/rules_go/blob/master/go/extras.rst#id3) is pretty darn simple: takes a list of input files - whatever they are - and
turns them into a go file you can just `import` from your code.
through a variable in that library.

For example, with:

    go_embed_data(
        name = "data",
        package = "assets",
        srcs = [
          "index.html",
          "favicon.ico",
        ],
    )

a `.go` file containing `index.html` and `favicon.ico` is created. To use it, we
need to include it in a library, for example, with:

    go_library(
        name = "assets",
        srcs = [":data"],
        importpath = "github.com/ccontavalli/bazel-example/backend/assets",
        visibility = ["//visibility:public"],
    )

From the go code, we can access the content of those files with:

    import "github.com/ccontavalli/backend/assets"
    ...
    log.Printf("index.html: %s", string(assets.Data["index.html"]))
    ...

And we can build the code after running the usual `bazel run //:gazelle` as
we were before with `bazel build backend` and `bazel run backend`.

In our go code, `assets.Data` will be a `map` using the path of the input file as key, and
returning the content as a `[]byte`, array of bytes.

## Including our react application as an asset

Including the `release.js` file generated from `frontend/BUILD.bazel` is relatively
simple. In our `backend/assets/BUILD.bazel`, all we have to do is modify the
`go_embed_data` to have:

    go_embed_data(
        name = "data",
        package = "assets",
        srcs = [
            "favicon.ico",
            "index.html",
            "//frontend:release",
        ],
        visibility = ["//visibility:public"],
    )

If you now run `bazel run frontend`, you will notice the longer build time,
and the fact that frontend dependencies are built and linked in.

And of course, the frontend will start listening and serving on [http://127.0.0.1:5432/](http://127.0.0.1:5432/).

## Including more assets

One thing that is annoying with `go_embed_data` is that the files are imported with weird
paths. Basically, `favicon.ico` and `index.html` will be accessible just by name, while `//frontend:release`
will be accessible by the full path, like `dist/bin/frontend/release.js`.

To mitigate the necessary path guessing and mangling, you can, well, do something else.
Not use `go_embed_data`, maybe create a docker image containing the files in a directory, or a .tar.gz, all easy to do with bazel.

I personally like embedding assets directly into the binary, so it is self contained,
easy to deploy and run even without docker. Another approach you can use is to have
bazel copy all the assets in one place, and only the use `go_embed_data`.

One way to do so is by using the `pkg_web` rule, provided with the nodejs ruleset. You
can check the `backend/assets/BUILD.bazel` file, but in short, you will end up with
something like:

    # Move all the assets in one place.
    pkg_web(
        name = "data",
        additional_root_paths = [
            "frontend"
        ],
        srcs = [
            "favicon.ico",
            "index.html",
            "//frontend:release",
        ]
    )

    # Generate a .go file containing all the assets as
    # go strings in the Data map.
    go_embed_data(
        name = "embedded",
        package = "assets",
        srcs = [":data"],
        visibility = ["//visibility:public"],
    )

    # Create a .go library containing the source file generated
    # by the go_embed_data target.
    go_library(
        name = "assets",
        srcs = [":embedded"],
        importpath = "github.com/ccontavalli/bazel-example/backend/assets",
        visibility = ["//visibility:public"],
    )

## Seeing the hello world in your new backend

# Importing material-ui

TODO

# Using SASS

TODO

# Using protocol buffers instead of json

TODO

# Writing and running unit tests

TODO

# Pulling it all together

# Conclusions

If you get strange errors related to yarn.lock file or packages.lock
or end up in situations where the behavior is inconsistent with what you
would expect, run:

    bazel clean
    bazel clean --expunge # only if bazel clean is not enough.

This will clean some of the internal bazel state, and start all the
next builds from scratch.
