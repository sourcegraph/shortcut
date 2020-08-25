# shortcut

This program runs on `cs.dev` and lets you quickly search code on Sourcegraph by typing `cs.dev/my query` into your browser's URL bar.

Try it!

- [`cs.dev/lang:python import yaml`](https://cs.dev/lang%3Apython%20import%20yaml)

## Release a new version

```
docker build -t sourcegraph/shortcut . && \
docker push sourcegraph/shortcut
```

For internal Sourcegraph usage, then bump the deployed version by updating the SHA-256 image digest in all files that refer to `sourcegraph/shortcut:latest@sha256:...`.
