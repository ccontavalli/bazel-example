Documentation entry point: https://bazelbuild.github.io/rules_nodejs/

Initialized workspace with:
yarn create @bazel --packageManager=yarn --typescript scratch

Configured SASS as per instructions here:
https://github.com/bazelbuild/rules_sass/blob/master/README.md

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
