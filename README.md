Initialized workspace with:
npm init @bazel --packageManager=yarn --typescript scratch

Configured SASS as per instructions here:
https://github.com/bazelbuild/rules_sass/blob/master/README.md

Followed instructions here to setup a ts_devserver:
https://bazelbuild.github.io/rules_nodejs/TypeScript.html#serving-typescript-for-development

From a directory outside my repository installed ibazel globally:

npm install -g @bazel/ibazel

Verified that the path used by npm is in my `$PATH`:

    $ npm bin -g
    /home/ccontavalli/.nvm/versions/node/v12.14.0/bin
