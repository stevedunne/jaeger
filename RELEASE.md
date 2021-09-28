# Jaeger Backend Release Process

1. Create a PR "Preparing release X.Y.Z" against master or maintenance branch ([example](https://github.com/jaegertracing/jaeger/pull/543/files)) by updating CHANGELOG.md to include:
    * A new section with the header `<X.Y.Z> (YYYY-MM-DD)`
    * A curated list of notable changes and links to PRs. Do not simply dump git log, select the changes that affect the users. To obtain the list of all changes run `make changelog` or use `scripts/release-notes.py`.
    * The section can be split into sub-section if necessary, e.g. UI Changes, Backend Changes, Bug Fixes, etc.
    * If the submodules have new releases, please also upgrade the submodule versions then commit, for example:
        ```
        cd jaeger-ui
        git ls-remote --tags origin
        git fetch
        git checkout {new_version} //e.g. v1.5.0
        ```
2. Add all merged pull requests to the milestone for the release and create a new milestone for a next release e.g. `Release 1.16`.
3. After the PR is merged, create a release on Github:
    * Title "Release X.Y.Z"
    * Tag `vX.Y.Z` (note the `v` prefix) and choose appropriate branch
    * Copy the new CHANGELOG.md section into the release notes
4. The release tag will trigger a build of the docker images. Since forks don't have jaegertracingbot dockerhub token, they can never publish images to jaegertracing organisation.
   1. Check the images are available on [Docker Hub](https://hub.docker.com/r/jaegertracing/).
   2. For monitoring and troubleshooting, refer to the [jaegertracing/jaeger GithubActions tab](https://github.com/jaegertracing/jaeger/actions).
5. [Publish documentation](https://github.com/jaegertracing/documentation/blob/master/RELEASE.md) for the new version in [jaegertracing.io](https://www.jaegertracing.io/docs/latest).
   1. Check [jaegertracing.io](https://www.jaegertracing.io/docs/latest) redirects to the new documentation release version URL.
   2. For monitoring and troubleshooting, refer to the [jaegertracing/documentation GithubActions tab](https://github.com/jaegertracing/documentation/actions).
6. Announce the release on the [mailing list](https://groups.google.com/g/jaeger-tracing), [slack](https://cloud-native.slack.com/archives/CGG7NFUJ3), and [twitter](https://twitter.com/JaegerTracing?lang=en).

Maintenance branches should follow naming convention: `release-major.minor` (e.g.`release-1.8`).

## Release managers

A Release Manager is the person responsible for ensuring that a new version of Jaeger is released. This person will coordinate the required changes, including to the related components such as UI, IDL, and jaeger-lib and will address any problems that might happen during the release, making sure that the documentation above is correct.

In order to ensure that knowledge about releasing Jaeger is spread among maintainers, we rotate the role of Release Manager among maintainers.

Here are the release managers for future versions with the tentative release dates. The release dates are the first Wednesday of the month, and we might skip a release if not enough changes happened since the previous release. In such case, the next tentative release date is the first Wednesday of the subsequent month.

| Version   | Release Manager  | Tentative release date |
|-----------|------------------|------------------------|
| 1.26.0    | @yurishkuro      |  1 September 2021      |
| 1.27.0    | @joe-elliott     |  6 October   2021      |
| 1.28.0    | @albertteoh      |  3 November  2021      |
| 1.29.0    | @jpkrohling      |  1 December  2021      |
| 1.30.0    | @pavolloffay     |  5 January   2022      |
| 1.31.0    | @vprithvi        |  2 February  2022      |
