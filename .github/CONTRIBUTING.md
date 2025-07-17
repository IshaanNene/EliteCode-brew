# Contributing to Elitecode

First off, thank you for considering contributing to Elitecode! It's people like you that make Elitecode such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for Elitecode. Following these guidelines helps maintainers and the community understand your report, reproduce the behavior, and find related reports.

Before creating bug reports, please check [this list](#before-submitting-a-bug-report) as you might find out that you don't need to create one. When you are creating a bug report, please [include as many details as possible](#how-do-i-submit-a-good-bug-report).

#### Before Submitting A Bug Report

* Check the [documentation](../README.md) for a list of common questions and problems.
* Check if the problem has already been reported in the [issue tracker](https://github.com/yourusername/elitecode/issues).
* Ensure you're using the latest version of Elitecode.

#### How Do I Submit A Good Bug Report?

Bugs are tracked as [GitHub issues](https://github.com/yourusername/elitecode/issues). Create an issue and provide the following information:

* Use a clear and descriptive title
* Describe the exact steps which reproduce the problem
* Provide specific examples to demonstrate the steps
* Describe the behavior you observed after following the steps
* Explain which behavior you expected to see instead and why
* Include screenshots and animated GIFs if possible
* Include your environment details (OS, Go version, etc.)

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for Elitecode, including completely new features and minor improvements to existing functionality.

#### Before Submitting An Enhancement Suggestion

* Check if the enhancement has already been suggested in the [issue tracker](https://github.com/yourusername/elitecode/issues).
* Check if the enhancement is compatible with the project's goals and scope.

#### How Do I Submit A Good Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://github.com/yourusername/elitecode/issues). Create an issue and provide the following information:

* Use a clear and descriptive title
* Provide a step-by-step description of the suggested enhancement
* Provide specific examples to demonstrate the steps
* Describe the current behavior and explain which behavior you expected to see instead
* Explain why this enhancement would be useful to most Elitecode users
* List some other tools or applications where this enhancement exists, if applicable

### Pull Requests

* Fill in the required template
* Do not include issue numbers in the PR title
* Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
* Include thoughtfully-worded, well-structured tests
* Document new code based on the [Documentation Styleguide](#documentation-styleguide)
* End all files with a newline

## Development Process

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
   * Follow the [style guides](#style-guides)
   * Add or update tests as needed
   * Add or update documentation as needed
4. Run tests and linters
   ```bash
   make test
   make lint
   ```
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Style Guides

### Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line
* Consider starting the commit message with an applicable emoji:
    * 🎨 `:art:` when improving the format/structure of the code
    * 🐎 `:racehorse:` when improving performance
    * 🚱 `:non-potable_water:` when plugging memory leaks
    * 📝 `:memo:` when writing docs
    * 🐛 `:bug:` when fixing a bug
    * 🔥 `:fire:` when removing code or files
    * 💚 `:green_heart:` when fixing the CI build
    * ✅ `:white_check_mark:` when adding tests
    * 🔒 `:lock:` when dealing with security
    * ⬆️ `:arrow_up:` when upgrading dependencies
    * ⬇️ `:arrow_down:` when downgrading dependencies

### Go Style Guide

* Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
* Use `gofmt` to format your code
* Document all exported functions, types, and variables
* Keep functions focused and small
* Use meaningful variable names
* Handle errors appropriately
* Write tests for new code

### Documentation Styleguide

* Use [Markdown](https://guides.github.com/features/mastering-markdown/)
* Reference functions, classes, and modules in backticks
* Use section headers for organization
* Include code examples when relevant
* Keep documentation up to date with code changes

## Additional Notes

### Issue and Pull Request Labels

This section lists the labels we use to help us track and manage issues and pull requests.

* `bug` - Issues that are bugs
* `documentation` - Issues about documentation
* `duplicate` - Issues that are duplicates of other issues
* `enhancement` - Issues that are feature requests
* `good first issue` - Good for newcomers
* `help wanted` - Extra attention is needed
* `invalid` - Issues that aren't valid
* `question` - Further information is requested
* `wontfix` - Issues that won't be worked on

## Recognition

Contributors who submit a PR that gets merged will be added to the [Contributors](../README.md#contributors) section of the README. 