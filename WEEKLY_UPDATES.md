# Weekly Updates

## 26.09.2023

Hello!

We’ve received some feedback that it’s currently hard to track the status of OpenTofu and what is currently being worked on. We hear you, and we’ll do better. We’re slowly moving discussions from our internal slack to public GitHub issues, to make them more inclusive. In addition, we’ll now be publishing an update like this once a week to give you a breakdown of where we’re at, what we’re currently working on, how you can best help, and what’s up next.

- Current Status
  - Right now we’re primarily working towards the first alpha release. There are two remaining blockers here:
    - Finishing the renaming to OpenTofu. There's only [a single PR](https://github.com/opentofu/opentofu/pull/576) left here.
    - Preparing the alpha registry replacement to be ready for use. The main current issue here is [#97](https://github.com/opentofu/registry/issues/97), described in a bit more detail below.
  - The [current registry design](https://github.com/opentofu/registry) is a glorified GitHub redirector.
    - This is by no means the final design! It’s meant to be a “get it working” design for the alpha release. Soon (~once the alpha release is out) we’ll start public discussions regarding the “stable” design for the registry, and we’re not ruling out any options here yet.
    - Right now the main issue is caching. Since the registry is a GitHub redirector, it is affected by GitHub rate limits. The rate limit is 5k requests per hour. We already do extensive caching here, but we can still do better and want to do better prior to the alpha release. [This issue contains more details](https://github.com/opentofu/registry/issues/97) about the work being done here.
    - Also, even though the OpenTofu alpha will [skip signature validation](https://github.com/opentofu/opentofu/issues/266) **if the key is not available**, we’re collecting provider signing keys for the registry to serve. So, if you’re a provider author, please make sure to create a Pull Request to the registry repo and add your public key. [Here's the latest example of such a PR](https://github.com/opentofu/registry/pull/95).
    - We’ve also made a mirror of all official HashiCorp providers so that they’re hosted on the OpenTofu GitHub, like any other 3rd-party provider. This included getting all historic versions to build (which was quite a challenge!) but we’re done with it and the end result is that out of 2260 provider versions we only have 16 failures, which are generally old and broken versions, so you should be able to easily work around this.
      - Long-term we’re planning to host these most-used providers on Fastly.
- Up next
  - Early next week we’re planning to make available an alpha release of OpenTofu, with a fully drop-in working registry for providers and modules.
  - With the above, we’ll also kick off the discussion around the stable registry.
    - We’ll publish an initial list of requirements for the registry in the [original registry issue](https://github.com/opentofu/opentofu/issues/258). The list itself will be open to changes based on further discussion.
    - Based on those requirements, we will be creating RFC’s for possible solutions. If you’d like to propose a solution, feel free to post an RFC too!
    - After having discussions there and possibly doing some PoC’s, the technical steering committee will make the final call which approach to pursue.
- How can I help?
  - Right now the best way to help is to create issues, discuss on issues, and spread the word about OpenTofu.
  - There are some occasional minor issues which are accepted and open to external contribution, esp. ones outside the release-blocking path. We’re also happy to accept any minor refactors or linter fixes. [Please see the contributing guide for more details](https://github.com/opentofu/opentofu/blob/main/CONTRIBUTING.md).
  - We have multiple engineers available full-time in the core team, so we’re generally trying to own any issues that are release blockers - this way we can make sure we get to the release as soon as possible.
  - The amount of pending-decision-labeled issues on the repository might be a bit off-putting. The reason for that is that right now we’re prioritizing the alpha and stable release. Only after we have a stable release in place do we aim to start actually accepting enhancement proposals and getting them implemented/merged. Still, we encourage you to open those issues and discuss them!
    - Issues and Pull Requests with enhancements or major changes will generally be frozen until we have the first stable release out. We will introduce a milestone to mark them as such more clearly.

Please let us know if you have any feedback on what we could improve, either with these updates or more generally! We're available on Slack, via GitHub issues, or even in the pull request creating this very file.
