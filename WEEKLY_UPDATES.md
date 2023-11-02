# Weekly Updates

## 2023-11-02

Hey!

This week was still mostly spent on discussing and PoC'in stable registry designs, as well as work on the S3 backend.

- Current Status
  - OpenTofu Day will be happening as a KubeCon + CloudNativeCon Europe CNCF-hosted Co-located Event, learn more [here](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/opentofu-day/#registration-details)!
  - The technical steering committee is meeting today to pick a stable registry design.
    - There are [three registry RFCs right now](https://github.com/opentofu/opentofu/issues?q=is%3Aopen+is%3Aissue+label%3Arfc+label%3Af-registry)
  - [Extensive work is happening](https://github.com/opentofu/opentofu/issues/700) around updating the S3 state backend configuration to be in line with the AWS provider configuration. Only two pull requests left here and we'll be done.
  - [The proper fix](https://github.com/opentofu/opentofu/pull/773) for the [interesting bug report](https://github.com/opentofu/opentofu/issues/763) mentioned last week in the update has been finished and merged.
- Up next
  - The technical steering committee will pick a stable registry design later today, and the core engineering team will focus most of its time on that in the upcoming weeks.
- How can I help?
  - Since the alpha release is out, right now the best way to help is to take the alpha for a test drive and see if there are any bugs / issues.
  - You can also submit and discuss RFCs related to the stable registry.
  - Other than that, the best way to help is to create issues, discuss on issues, and spread the word about OpenTofu.
  - There are some occasional minor issues which are accepted and open to external contribution, esp. ones outside the release-blocking path. We’re also happy to accept any minor refactors or linter fixes. [Please see the contributing guide for more details](https://github.com/opentofu/opentofu/blob/main/CONTRIBUTING.md).
  - We have multiple engineers available full-time in the core team, so we’re generally trying to own any issues that are release blockers - this way we can make sure we get to the release as soon as possible.
  - The amount of pending-decision-labeled issues on the repository might be a bit off-putting. The reason for that is that right now we’re prioritizing the alpha and stable release. Only after we have a stable release in place do we aim to start actually accepting enhancement proposals and getting them implemented/merged. Still, we encourage you to open those issues and discuss them!
    - Issues and Pull Requests with enhancements or major changes will generally be frozen until we have the first stable release out. We will introduce a milestone to mark them as such more clearly.

Please let us know if you have any feedback on what we could improve, either with these updates or more generally - this very document was a result of such feedback! We're available on Slack, via GitHub issues, or even in the pull request updating this file.

## 2023-10-24

Hey!

This week was again spent on discussing and PoC'in stable registry designs, as well as work on the S3 backend.

- Current Status
  - There are [two registry RFCs right now](https://github.com/opentofu/opentofu/issues?q=is%3Aopen+is%3Aissue+label%3Arfc+label%3Af-registry), and both will most likely have a PoC by the end of the week.
  - [Extensive work is happening](https://github.com/opentofu/opentofu/issues/700) around updating the S3 state backend configuration to be in line with the AWS provider configuration.
  - OpenTofu is available on Homebrew Core, so you can just `brew install opentofu` to use it!
  - A [bug report](https://github.com/opentofu/opentofu/issues/763) has surfaced an interesting issue when migrating from Terraform to OpenTofu with an existing state file. An initial fix [has already been implemented](https://github.com/opentofu/opentofu/pull/773) and will be further polished and released soon.
- Up next
  - The steering committee will pick a registry RFC in the week of 30th Oct - 3rd Nov.
    - Intense implementation work will start on the 6th of November.
    - This should let us have a stable registry in place around mid-November - mid-December.
- How can I help?
  - Since the alpha release is out, right now the best way to help is to take the alpha for a test drive and see if there are any bugs / issues.
  - You can also submit and discuss RFCs related to the stable registry.
  - Other than that, the best way to help is to create issues, discuss on issues, and spread the word about OpenTofu.
  - There are some occasional minor issues which are accepted and open to external contribution, esp. ones outside the release-blocking path. We’re also happy to accept any minor refactors or linter fixes. [Please see the contributing guide for more details](https://github.com/opentofu/opentofu/blob/main/CONTRIBUTING.md).
  - We have multiple engineers available full-time in the core team, so we’re generally trying to own any issues that are release blockers - this way we can make sure we get to the release as soon as possible.
  - The amount of pending-decision-labeled issues on the repository might be a bit off-putting. The reason for that is that right now we’re prioritizing the alpha and stable release. Only after we have a stable release in place do we aim to start actually accepting enhancement proposals and getting them implemented/merged. Still, we encourage you to open those issues and discuss them!
    - Issues and Pull Requests with enhancements or major changes will generally be frozen until we have the first stable release out. We will introduce a milestone to mark them as such more clearly.

Please let us know if you have any feedback on what we could improve, either with these updates or more generally - this very document was a result of such feedback! We're available on Slack, via GitHub issues, or even in the pull request updating this file.

## 2023-10-17

Hello there!

This past week was mostly spent on creating stable registry designs, which is a prerequisite to building a stable provider/module registry and releasing a stable version of OpenTofu.

- Current Status
  - The [registry replacement issue](https://github.com/opentofu/opentofu/issues/258) has now been updated to include a list of stable registry requirements.
    - There is now time until the 27th October to submit and discuss RFCs.
    - You can find relevant RFCs by using [GitHub Issue search](https://github.com/opentofu/opentofu/issues?q=is%3Aopen+is%3Aissue+label%3Arfc+label%3Af-registry).
    - After the 27th October the Technical Steering Committee will pick a single RFC which will then be implemented.
  - A [setup-opentofu](https://github.com/opentofu/setup-opentofu) GitHub Action is now available, it's a great way to take OpenTofu for a spin on GitHub Actions. Thanks a lot to @kislerdm for contributing this!
  - Work on TF 1.6 compatibility has been done.
  - We've closed a lot of issues that were stable release blockers, you can take a look at the current stable milestone status [here](https://github.com/opentofu/opentofu/milestone/3).
- Up next
  - The main TF 1.6 compatibility issue left, on which work has already been started and will continue for a while is the new S3 state backend configuration schema. You can track this issue [here](https://github.com/opentofu/opentofu/issues/700).
  - Until the 27th October deadline, work and discussion on RFCs will happen, as well as the implementation of PoCs for some of these RFCs.
- How can I help?
  - Since the alpha release is out, right now the best way to help is to take the alpha for a test drive and see if there are any bugs / issues.
  - You can also submit and discuss RFCs related to the stable registry.
  - Other than that, the best way to help is to create issues, discuss on issues, and spread the word about OpenTofu.
  - There are some occasional minor issues which are accepted and open to external contribution, esp. ones outside the release-blocking path. We’re also happy to accept any minor refactors or linter fixes. [Please see the contributing guide for more details](https://github.com/opentofu/opentofu/blob/main/CONTRIBUTING.md).
  - We have multiple engineers available full-time in the core team, so we’re generally trying to own any issues that are release blockers - this way we can make sure we get to the release as soon as possible.
  - The amount of pending-decision-labeled issues on the repository might be a bit off-putting. The reason for that is that right now we’re prioritizing the alpha and stable release. Only after we have a stable release in place do we aim to start actually accepting enhancement proposals and getting them implemented/merged. Still, we encourage you to open those issues and discuss them!
    - Issues and Pull Requests with enhancements or major changes will generally be frozen until we have the first stable release out. We will introduce a milestone to mark them as such more clearly.

Please let us know if you have any feedback on what we could improve, either with these updates or more generally - this very document was a result of such feedback! We're available on Slack, via GitHub issues, or even in the pull request updating this file.

## 2023-10-11

Hey!

Last week we've released the [first alpha release](https://github.com/opentofu/opentofu/releases/tag/v1.6.0-alpha1) (and [another one](https://github.com/opentofu/opentofu/releases/tag/v1.6.0-alpha2) the next day). This has been a huge milestone for the OpenTofu community, and means that you can actually download OpenTofu and play with it. It also marks the beginning of the stable registry design discussions.

Just to give you some metrics about the release:
- The alpha releases of OpenTofu have over 2k downloads.
- Some registry stats
  - Uptime: 100% 
  - Max RPS: 105 
  - Total Requests: 78869 
  - Average latency: 19.15ms 
  - Total unique providers used: 194

We'll also have [our first office hours](https://opentofu.org/blog/opentofus-new-office-hours) today! Make sure you don't miss that.

- Current Status
  - First Alpha Release
    - We've managed to get the alpha release out, including [a working registry](https://github.com/opentofu/registry).
  - Second Alpha Release
    - A [bug was reported](https://github.com/opentofu/opentofu/issues/655) to OpenTofu and later also to Terraform. We've deemed it high priority and possibly impacting the ability of the community to test OpenTofu with existing codebases, so we've worked on fixing it asap. The root cause was in the [hcl library](https://github.com/hashicorp/hcl/issues/629). We've [made a fix](https://github.com/opentofu/opentofu/pull/659) in OpenTofu, released a new alpha, and afterward [the fix got accepted to the hcl library](https://github.com/hashicorp/hcl/pull/630).
- Up next
  - Stable Registry Design (see [the tracking issue](https://github.com/opentofu/opentofu/issues/258))
    - Even though there's a registry included with the alpha release, it's not stable. It was mostly a "get it working" design, and we will now be working on a stable design. 
    - Right now we're working on a list of requirements, which will most likely be shared tomorrow in the tracking issue above.
    - Then, there will be an approximate 2 weeks to submit and discuss RFCs, after which the steering committee will choose the final design, which we'll begin implementing right away.
    - Please discuss in the issue and incoming RFCs, we want to hear all your voices! Please also feel free to submit your own RFCs.
- How can I help?
  - Since the alpha release is out, right now the best way to help is to take the alpha for a test drive and see if there are any bugs / issues.
  - Other than that, the best way to help is to create issues, discuss on issues, and spread the word about OpenTofu.
  - There are some occasional minor issues which are accepted and open to external contribution, esp. ones outside the release-blocking path. We’re also happy to accept any minor refactors or linter fixes. [Please see the contributing guide for more details](https://github.com/opentofu/opentofu/blob/main/CONTRIBUTING.md).
  - We have multiple engineers available full-time in the core team, so we’re generally trying to own any issues that are release blockers - this way we can make sure we get to the release as soon as possible.
  - The amount of pending-decision-labeled issues on the repository might be a bit off-putting. The reason for that is that right now we’re prioritizing the alpha and stable release. Only after we have a stable release in place do we aim to start actually accepting enhancement proposals and getting them implemented/merged. Still, we encourage you to open those issues and discuss them!
    - Issues and Pull Requests with enhancements or major changes will generally be frozen until we have the first stable release out. We will introduce a milestone to mark them as such more clearly.

Please let us know if you have any feedback on what we could improve, either with these updates or more generally - this very document was a result of such feedback! We're available on Slack, via GitHub issues, or even in the pull request updating this file.

## 2023-10-02

Hey!

Another week, another update! This last week has been all about the alpha release, which is planned to be released this week, ideally tomorrow (Tuesday).

- Current Status
  - Alpha release is blocked only by the initial registry work
    - [There's a PR open](https://github.com/opentofu/opentofu/pull/379) for using the OpenTofu registry as the default in OpenTofu.
    - There's been extensive work on caching, discussed in [a GitHub issue](https://github.com/opentofu/registry/issues/97) and [a Slack channel](https://opentofucommunity.slack.com/archives/C05UH44L449).
      - [One PR](https://github.com/opentofu/registry/pull/103) introduced a better version listing cache in DynamoDB, to reuse the cache across provider listing and provider downloading.
      - [Another PR](https://github.com/opentofu/registry/pull/105) introduced http-level caching for all requests to GitHub, which generally avoids making the same request twice in quick succession.
    - There's a bug in the provider mirroring scripts which makes them not properly trigger GitHub Actions. We're working on a fix.
    - [There's a bug](https://github.com/opentofu/registry/issues/121) with module downloading, for which both tags in the form of `x.y.z` and `vx.y.z` are supported, while we only support one of them.
- Up next
  - The alpha release is planned for this week, using our alpha registry.
  - As described last week, once the alpha release is out, we'll be kicking off discussions about the stable registry design.
    - One such issue [has already been opened](https://github.com/opentofu/opentofu/issues/619) and explores the problem space of handling signatures and public keys.
- How can I help?
  - Right now the best way to help is to create issues, discuss on issues, and spread the word about OpenTofu. Additionally, the alpha release is coming out this week. We'd appreciate you taking it for a test drive and letting us know about any issues you encounter.
  - There are some occasional minor issues which are accepted and open to external contribution, esp. ones outside the release-blocking path. We’re also happy to accept any minor refactors or linter fixes. [Please see the contributing guide for more details](https://github.com/opentofu/opentofu/blob/main/CONTRIBUTING.md).
  - We have multiple engineers available full-time in the core team, so we’re generally trying to own any issues that are release blockers - this way we can make sure we get to the release as soon as possible.
  - The amount of pending-decision-labeled issues on the repository might be a bit off-putting. The reason for that is that right now we’re prioritizing the alpha and stable release. Only after we have a stable release in place do we aim to start actually accepting enhancement proposals and getting them implemented/merged. Still, we encourage you to open those issues and discuss them!
    - Issues and Pull Requests with enhancements or major changes will generally be frozen until we have the first stable release out. We will introduce a milestone to mark them as such more clearly.

Please let us know if you have any feedback on what we could improve, either with these updates or more generally! We're available on Slack, via GitHub issues, or even in the pull request creating this very file.

## 2023-09-26

Hello!

We’ve received some feedback that it’s currently hard to track the status of OpenTofu and what is currently being worked on. We hear you, and we’ll do better. We’re slowly moving discussions from our internal slack to public GitHub issues, to make them more inclusive. In addition, we’ll now be publishing an update like this once a week to give you a breakdown of where we’re at, what we’re currently working on, how you can best help, and what’s up next.

- Current Status
  - Right now we’re primarily working towards the first alpha release. There are two remaining blockers here:
    - Finishing the renaming to OpenTofu. There's only [a single PR](https://github.com/opentofu/opentofu/pull/576) left here.
    - Preparing the alpha registry replacement to be ready for use. The main current issue here is [#97](https://github.com/opentofu/registry/issues/97), described in a bit more detail below.
  - The [current registry design](https://github.com/opentofu/registry) is a glorified GitHub redirector.
    - This is by no means the final design! It’s meant to be a “get it working” design for the alpha release. Soon (~once the alpha release is out) we’ll start public discussions regarding the “stable” design for the registry, and we’re not ruling out any options here yet.
    - Right now the main issue is caching. Since the registry is a GitHub redirector, it is affected by GitHub rate limits. The rate limit is 5k requests per hour. We already do extensive caching here, but we can still do better and want to do better prior to the alpha release. [This issue contains more details](https://github.com/opentofu/registry/issues/97) about the work being done here.
    - Also, even though OpenTofu alpha will [skip signature validation](https://github.com/opentofu/opentofu/issues/266) **when the key is not available**, we’re collecting provider signing keys for the registry to serve. So, if you’re a provider author, please make sure to create a Pull Request to the registry repo and add your public key. [Here's the latest example of such a PR](https://github.com/opentofu/registry/pull/95).
    - We’ve also made a mirror of all official HashiCorp providers so that they’re hosted on the OpenTofu GitHub, like any other 3rd-party provider. This included getting all historic versions to build (which was quite a challenge!) but we’re done with it and the end result is that out of 2260 provider versions we only have 16 failures, which are generally old and broken versions, so you should be able to easily work around this.
      - The reason for this is that the registry being a GitHub redirector needs all artifacts on GitHub, while HashiCorp hosts these providers directly on their registry and doesn't include them in GitHub releases.
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
  - Once we make the alpha release available next week, the best way to help will be by test-driving that release and creating GitHub issues for any problems you find.

Please let us know if you have any feedback on what we could improve, either with these updates or more generally! We're available on Slack, via GitHub issues, or even in the pull request creating this very file.
