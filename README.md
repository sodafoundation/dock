# SODA DOCK

[![Go Report Card](https://goreportcard.com/badge/github.com/sodafoundation/dock?branch=master)](https://goreportcard.com/report/github.com/sodafoundation/dock)
[![Build Status](https://travis-ci.org/sodafoundation/dock.svg?branch=master)](https://travis-ci.org/sodafoundation/dock)
[![codecov.io](https://codecov.io/github/sodafoundation/dock/coverage.svg?branch=master)](https://codecov.io/github/sodafoundation/dock?branch=master)
[![Releases](https://img.shields.io/github/release/sodafoundation/dock/all.svg?style=flat-square)](https://github.com/sodafoundation/dock/releases)
[![LICENSE](https://img.shields.io/github/license/sodafoundation/dock.svg?style=flat-square)](https://github.com/sodafoundation/dock/blob/master/LICENSE)

<img src="https://sodafoundation.io/wp-content/uploads/2020/01/SODA_logo_outline_color_800x800.png" width="200" height="200">

## Introduction

SODA Dock is an open source implementation for the unified interface to connect heterogeneous storage backends. So dock is a docking station for heterogeneous storage backends! This is where all the different storage vendors drivers for various storage backend models get attached.

It is part of SODA Terra (SDS Controller). There are other two repositories part of SODA Terra viz., [API](https://github.com/sodafoundation/api) and [Controller](https://github.com/sodafoundation/controller)


We strive to make most of the protocols and backends supported as close as ‘plug n play’. Currently, each storage backend needs a thin, easy to develop SODA Driver Plugin to connect the storage backend to the DOCK. The SODA Driver Plugin and Storage vendor driver together can be called SODA Driver for xxx vendor yy model storage. SODA Driver can support one or more or multiple classes of storage backends.

SODA DOCK can interface directly to SODA API or via Controller. We recommend through the controller for a complete end to end solution, as it can provide the metadata management, handling multiple dock etc. For the api to dock direct interfacing, currently the user needs to do the necessary changes.

Plan to have multiple instance, multi driver docks to support multi-cluster, multi-platform or multi-cloud scenarios in future.

This is one of the SODA Core Projects and is maintained by SODA Foundation directly. We recommend adding all the storage vendor drivers under this project to build a huge repository for the storage vendor support. However the soda driver plugins can be maintained anywhere, and if it is compliant with SODA API, it can be part of SODA Project Landscape.

Earlier part of github.com/sodafoundation/soda Or github.com/opensds/opensds

## Documentation

[https://docs.sodafoundation.io](https://docs.sodafoundation.io/)

## Quick Start - To Use/Experience

[https://docs.sodafoundation.io](https://docs.sodafoundation.io/)

## Quick Start - To Develop

[https://docs.sodafoundation.io](https://docs.sodafoundation.io/)

## Latest Releases

[https://github.com/sodafoundation/dock/releases](https://github.com/sodafoundation/dock/releases)

## Support and Issues

[https://github.com/sodafoundation/dock/issues](https://github.com/sodafoundation/dock/issues)

## Project Community

[https://sodafoundation.io/slack/](https://sodafoundation.io/slack/)

## How to contribute to this project?

Join [https://sodafoundation.io/slack/](https://sodafoundation.io/slack/) and share your interest in the ‘general’ channel

Checkout [https://github.com/sodafoundation/dock/issues](https://github.com/sodafoundation/dock/issues) labelled with ‘good first issue’ or ‘help needed’ or ‘help wanted’ or ‘StartMyContribution’ or ‘SMC’

## Project Roadmap

Envision to have huge support for all the industry storage vendor driver support under dock with a standardized and unified storage backend interface.

[https://docs.sodafoundation.io](https://docs.sodafoundation.io/)

## Join SODA Foundation

Website : [https://sodafoundation.io](https://sodafoundation.io/)

Slack  : [https://sodafoundation.io/slack/](https://sodafoundation.io/slack/)

Twitter  : [@sodafoundation](https://twitter.com/sodafoundation)

Mailinglist  : [https://lists.sodafoundation.io](https://lists.sodafoundation.io/)
