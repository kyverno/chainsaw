# Pull Requests

You can contribute to Chainsaw by making a [pull request] that will be reviewed by maintainers and integrated into the main repository when the changes made are approved. You can contribute bug fixes, documentation changes, or new functionalities.

[pull request]: https://docs.github.com/en/pull-requests

!!! note "Considering a pull request"

    Before deciding to spend effort on making changes and creating a pull request, please discuss what you intend to do. If you are responding to what you think might be a bug, please issue a [bug report] first. If you intend to work on documentation, create a [documentation issue]. If you want to work on a new feature, please create a [change request].

    Keep in mind the guidance given and let people advise you. It might be that there are easier solutions to the problem you perceive and want to address. It might be that what you want to achieve can already be done by configuration or [customization].

[bug report]: reporting-a-bug.md
[documentation issue]: reporting-a-docs-issue.md
[change request]: requesting-a-change.md
<!-- [customization]: ../customization.md -->

## Learning about pull requests

Pull requests are a concept layered on top of Git by services that provide Git hosting. Before you consider making a pull request, you should familiarize yourself with the documentation on GitHub, the service we are using. The following articles are of particular importance:

1. [Forking a repository]
2. [Creating a pull request from a fork]
3. [Creating a pull request]

Note that they provide tailored documentation for different operating systems and different ways of interacting with GitHub. We do our best in the documentation here to describe the process as it applies to Chainsaw but cannot cover all possible combinations of tools and ways of doing things. It is also important that you understand the concept of a pull-request in general before continuing.

[Forking a repository]: https://docs.github.com/en/get-started/quickstart/fork-a-repo
[Creating a pull request from a fork]: https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork
[Creating a pull request]: https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request

## Pull request process

In the following, we describe the general process for making pull requests. The aim here is to provide the 30k ft overview before describing details later on.

### Preparing changes and draft PR

The diagram below describes what typically happens to repositories in the process or preparing a pull request. We will be discussing the review-revise process below. It is important that you understand the overall process first before you worry about specific commands. This is why we cover this first before providing instructions below.

``` mermaid
sequenceDiagram
  autonumber

  participant chainsaw
  participant PR
  participant fork
  participant local

  chainsaw ->> fork: fork on GitHub
  fork ->> local: clone to local
  local ->> local: branch
  loop prepare
    loop push
      loop edit
        local ->> local: commit
      end
      local ->> fork: push
    end
    chainsaw ->> fork: merge in any changes
    fork ->>+ PR: create draft PR
    PR ->> PR: review your changes
  end
```

1. Fork the Repository: Fork the Chainsaw repository on GitHub to create your own copy.
2. Clone to Local: Clone your fork to your local machine.
3. Create a Branch: Create a topic branch for your changes.
4. Set Up Development Environment: Follow the instructions to set up a development environment.
5. Iterate and Commit: Make incremental changes and commit them with meaningful messages.
6. Push Regularly: Push your commits to your fork regularly.
7. Merge Changes from Upstream: Regularly merge changes from the original Chainsaw repository to avoid conflicts.
8. Create a Draft Pull Request: Once satisfied with your changes, create a draft pull request for early feedback.
9. Review and Revise: Review your work critically, address feedback, and refine your changes.

### Finalizing

Once you are happy with your changes, you can move to the next step, finalizing
your pull request and asking for a more formal and detailed review. The diagram
below shows the process:

``` mermaid
sequenceDiagram
  autonumber
  participant chainsaw
  participant PR
  participant fork
  participant local

  activate PR
  PR ->> PR: finalize PR
  loop review
    loop discuss
      PR ->> PR: request review
      PR ->> PR: discussion
      local ->> fork: push further changes
    end
    PR ->> chainsaw: merge (and squash)
    deactivate PR
    fork ->> fork: delete branch
    chainsaw ->> fork: pull
    local ->> local: delete branch
    fork ->> local: pull
  end

```

1. Finalize PR: Signal that your changes are ready for review.
2. Request Review: Ask the maintainer to review your changes.
3. Discuss and Revise: Engage in discussions, make necessary revisions, and iterate.
4. Merge and Squash: Once approved, the maintainer will merge and possibly squash your commits.
5. Clean Up: Delete the branch used for the PR from both your fork and local clone.