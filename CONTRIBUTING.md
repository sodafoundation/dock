# Soda

[![Go Report Card](https://goreportcard.com/badge/github.com/soda/soda?branch=master)](https://goreportcard.com/report/github.com/soda/soda)
[![Build Status](https://travis-ci.org/soda/soda.svg?branch=master)](https://travis-ci.org/soda/soda)
[![Coverage Status](https://coveralls.io/repos/github/soda/soda/badge.svg?branch=master)](https://coveralls.io/github/soda/soda?branch=master)

<img src="https://www.soda.io/wp-content/uploads/sites/18/2016/11/logo_soda.png" width="100">

## How to contribute

soda is Apache 2.0 licensed and accepts contributions via GitHub pull requests. This document outlines some of the conventions on commit message formatting, contact points for developers and other resources to make getting your contribution into soda easier.

## Email and chat

- Email: [soda-dev](https://lists.soda.io/mailman/listinfo)
- Slack: #[soda](https://soda.slack.com)

Before you start, NOTICE that ```master``` branch is the relatively stable version
provided for customers and users. So all code modifications SHOULD be submitted to
`development` branch.

## Getting started

- Fork the repository on GitHub.
- Read the README.md and INSTALL.md for project information and build instructions.

For those who just get in touch with this project recently, here is a proposed contributing [tutorial](https://github.com/leonwanghui/installation-note/blob/master/soda_fork_contribute_tutorial.md).

## Contribution Workflow

### Code style

The coding style suggested by the Golang community is used in soda. See the [doc](https://github.com/golang/go/wiki/CodeReviewComments) for more details.

Please follow this style to make soda easy to review, maintain and develop.

### Report issues

A great way to contribute to the project is to send a detailed report when you encounter an issue. We always appreciate a well-written, thorough bug report, and will thank you for it!

When reporting issues, refer to this format:

- What version of env (soda, os, golang etc) are you using?
- Is this a BUG REPORT or FEATURE REQUEST?
- What happened?
- What you expected to happen?
- How to reproduce it?(as minimally and precisely as possible)

### Propose PRs

- Raise your idea as an [issue](https://github.com/soda/soda/issues)
- If it is a new feature that needs lots of design details, a design proposal should also be submitted [here](https://github.com/soda/design-specs/pulls).
- After reaching consensus in the issue discussions and design proposal reviews, complete the development on the forked repo and submit a PR.
  Here are the [PRs](https://github.com/soda/soda/pulls?q=is%3Apr+is%3Aclosed) that are already closed.
- If a PR is submitted by one of the core members, it has to be merged by a different core member.
- After PR is sufficiently discussed, it will get merged, abondoned or rejected depending on the outcome of the discussion.

Thank you for your contribution !
