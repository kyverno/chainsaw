---
date: 2023-10-12
slug: wip
categories:
  - announcements
authors:
  - eddycharly
---

# :tada::tada::tada: It's happening - Finally ! :tada::tada::tada:

<p align="center">
  <img src="https://i.pinimg.com/originals/26/fe/ed/26feed1d744536fbe4bd1c52f8053564.gif" />
</p>

Hello everyone!

We finally started writing our own testing tool !

Let's join forces and make better, stronger **Open Source** and **Community Driven** tool.

<!-- more -->

## Why we made it ?

While developing [Kyverno](https://kyverno.io) we needed to run end to end tests to make sure our admission controller worked as expected.

[Kyverno](https://kyverno.io) can validate, mutate and generate resources based on policies installed in a cluster and a typical test is:

1. Create a policy
1. Create a resource
1. Check that Kyverno acted as expected
1. Cleanup and move to the next test

We started with another tool called [Kuttl](https://kuttl.dev). While [Kuttl](https://kuttl.dev) is great we identified some limitations and forked the tool to add the features we needed.

At some point we needed more flexibility than what [Kuttl](https://kuttl.dev) offered and we designed a new assertion model.

This was simpler to start from scratch than continuing making changes in our [Kuttl](https://kuttl.dev) fork.

## What's next ?

This is still WIP and needs a lot of work before we can consider it ready but things are moving fast.

We would love to build a community driven tool and welcome all contributors.

Feel free to fork this repository and start submitting pull requests :pray: