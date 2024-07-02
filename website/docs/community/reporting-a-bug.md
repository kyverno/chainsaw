# Bug Reports

Chainsaw, developed by Kyverno, is an actively maintained project that we constantly strive to improve. With a project of this size and complexity, bugs may occur. If you think you have discovered a bug, you can help us by submitting an issue in our public [issue tracker], following this guide.

[issue tracker]: https://github.com/kyverno/chainsaw/issues

## Before Creating an Issue

With numerous users, issues are created regularly. The maintainers of this project strive to address bugs promptly. By following this guide, you will know exactly what information we need to help you quickly.

__Please do the following before creating an issue:__

### Upgrade to Latest Version

Chances are that the bug you discovered was already fixed in a subsequent version. Before reporting an issue, ensure that you're running the [latest version] of Chainsaw. Consult our [upgrade guide] to learn how to upgrade to the latest version.

!!! warning "Bug fixes are not backported"

    Please understand that only bugs that occur in the latest version of Chainsaw will be addressed. Also, to reduce duplicate efforts, fixes cannot be backported to earlier versions.

<!-- [latest version]: ../changelog/index.md
[upgrade guide]: ../upgrade.md -->

### Remove Customizations

If you're using customizations like additional configurations, remove them before reporting a bug. We can't offer official support for bugs that might hide in your overrides, so make sure to omit custom settings from your configuration files.

__Don't be shy to ask on our [discussion board] for help if you run into problems.__

[discussion board]: https://github.com/kyverno/chainsaw/discussions

### Search for Solutions

At this stage, we know that the problem persists in the latest version and is not caused by any of your customizations. However, the problem might result from a small typo or a syntactical error in a configuration file.

Before creating a bug report, save time for us and yourself by doing some research:

1. [Search our documentation] for relevant sections related to your problem. Ensure everything is configured correctly.
2. [Search our issue tracker] as another user might have already reported the same problem.
3. [Search our discussion board] to see if other users are facing similar issues and find possible solutions.

__Keep track of all search terms and relevant links; you'll need them in the bug report.__

[Search our documentation]: ?q=
[issue tracker]: https://github.com/kyverno/chainsaw/issues
[discussion board]: https://github.com/kyverno/chainsaw/discussions

---

If you still haven't found a solution to your problem, create an issue. It's now likely that you've encountered something new. Read the following section to learn how to create a complete and helpful bug report.

## Issue Template

We have created a new issue template to make the bug reporting process as simple as possible and more efficient for our community and us. It consists of the following parts:

- [Title]
- [Context] <small>optional</small>
- [Bug Description]
- [Related Links]
- [Reproduction]
- [Steps to Reproduce]
- [Browser] <small>optional</small>
- [Checklist]

[Title]: #title
[Context]: #context
[Bug Description]: #bug-description
[Related Links]: #related-links
[Reproduction]: #reproduction
[Steps to Reproduce]: #steps-to-reproduce
[Browser]: #browser
[Checklist]: #checklist

### Title

A good title is short and descriptive. It should be a one-sentence executive summary of the issue, so the impact and severity of the bug can be inferred from the title.

| <!-- --> | Example  |
| -------- | -------- |
| :material-check:{ style="color: #4DB6AC" } __Clear__ | Chainsaw `apply` command fails with specific CRD
| :material-close:{ style="color: #EF5350" } __Wordy__ | The `apply` command in Chainsaw fails when used with a certain Custom Resource Definition
| :material-close:{ style="color: #EF5350" } __Unclear__ | Command does not work
| :material-close:{ style="color: #EF5350" } __Useless__ | Help

### Context <small>optional</small> {#context}

Before describing the bug, you can provide additional context to help us understand what you were trying to achieve. Explain the circumstances under which you're using Chainsaw, and what you think might be relevant. Don't describe the bug here.

### Bug Description

Provide a clear, focused, specific, and concise summary of the bug you encountered. Explain why you think this is a bug that should be reported to Chainsaw, and not to one of its dependencies. Follow these principles:

- __Explain the what, not the how__ – don't explain how to reproduce the bug here, we're getting there. Focus on articulating the problem and its impact.
- __Keep it short and concise__ – if the bug can be precisely explained in one or two sentences, perfect. Don't inflate it.
- __One bug at a time__ – if you encounter several unrelated bugs, create separate issues for them.

### Related Links

Share links to relevant sections of our documentation and any related issues or discussions. This helps us improve our documentation and understand the problem better.

### Reproduction

A minimal reproduction is essential for a well-written bug report, as it allows us to recreate the conditions necessary to inspect the bug. Follow the guide to create a reproduction:

[:material-bug: Create reproduction][Create reproduction]{ .md-button .md-button--primary }

After creating the reproduction, you should have a `.zip` file, ideally not larger than 1 MB. Drag and drop the `.zip` file into the issue field, which will automatically upload it to GitHub.

!!! warning "Don't share links to repositories"

    While linking to a repository is a common practice, we currently don't support this. The reproduction, created using the built-in info plugin, contains all necessary environment information.

<!-- [Create reproduction]: ../guides/creating-a-reproduction.md -->

### Steps to Reproduce

List specific steps to follow when running your reproduction to observe the bug. Keep the steps concise and ensure nothing is left out. Use simple language and focus on continuity.

### Browser <small>optional</small> {#browser}

If the bug only occurs in specific browsers, let us know which ones are affected. This field is optional, as it is only relevant for bugs that do not involve a crash when previewing or building your site.

!!! incognito "Incognito Mode"

    Verify that the bug is not caused by a browser extension by switching to incognito mode. If the bug disappears, it is likely caused by an extension.

### Checklist

Before submitting, ensure you have:

- Followed this guide thoroughly
- Provided all necessary information
- Created a minimal reproduction

Thanks for following the guide and creating a high-quality bug report. We will take it from here.

