# SODA DOCK

[![Go Report Card](https://goreportcard.com/badge/github.com/sodafoundation/dock?branch=master)](https://goreportcard.com/report/github.com/sodafoundation/dock)
[![Build Status](https://travis-ci.org/sodafoundation/dock.svg?branch=master)](https://travis-ci.org/sodafoundation/dock)
[![codecov.io](https://codecov.io/github/sodafoundation/dock/coverage.svg?branch=master)](https://codecov.io/github/sodafoundation/dock?branch=master)
[![Releases](https://img.shields.io/github/release/sodafoundation/dock/all.svg?style=flat-square)](https://github.com/sodafoundation/dock/releases)
[![LICENSE](https://img.shields.io/github/license/sodafoundation/dock.svg?style=flat-square)](https://github.com/sodafoundation/dock/blob/master/LICENSE)

![SODA Logo](https://sodafoundation.io/wp-content/uploads/2020/01/SODA_logo_outline_color_800x800.png)

## Overview

SODA Dock is an essential component of the SODA Terra, a software-defined storage (SDS) controller. It acts as a docking station for various storage backends, enabling a unified interface to connect heterogeneous storage systems.

## What is SODA Dock?

SODA Dock is where different storage vendor drivers for various backend models attach, providing a seamless integration point for storage solutions. The aim is to support a vast range of protocols and backends, ensuring most of them are as close to 'plug n play' as possible.

Each storage backend requires a thin SODA Driver Plugin for connection. This plugin, combined with the vendor driver, can be termed as the SODA Driver for a particular vendor model. A SODA Driver can cater to one or more classes of storage backends.

## Integration Points

The DOCK can interface directly with the SODA API or through the Controller. However, for a comprehensive solution, integration through the controller is recommended as it facilitates metadata management and handles multiple docks. For direct interfacing between the API and the DOCK, users must implement the necessary changes.

## Future Goals

Plans are in place to introduce multi-instance, multi-driver docks to support diverse scenarios, including multi-cluster, multi-platform, and multi-cloud.

## Affiliation

SODA Dock is one of the core projects maintained directly by the SODA Foundation. We encourage storage vendors to add their drivers under this project, creating a vast repository for storage support. However, SODA driver plugins can be hosted elsewhere, and if compliant with the SODA API, they can be integrated into the SODA Project Landscape.

Earlier, this was a part of [github.com/sodafoundation/opensds](https://github.com/sodafoundation/opensds) or [github.com/opensds/opensds](https://github.com/opensds/opensds).

## Useful Links

- **Documentation**: [SODA Foundation Documentation](https://docs.sodafoundation.io/)
  
- **Quick Start (Usage)**: [Getting Started Guide](https://docs.sodafoundation.io/)

- **Quick Start (Development)**: [Developer Guide](https://docs.sodafoundation.io/)

- **Latest Releases**: [Releases on GitHub](https://github.com/sodafoundation/dock/releases)

- **Support & Issues**: [GitHub Issues](https://github.com/sodafoundation/dock/issues)

- **Community**: [Join the SODA Slack](https://sodafoundation.io/slack/)

## How to Contribute?

1. Join the [SODA Slack Channel](https://sodafoundation.io/slack/) and share your interests in the ‘general’ channel.
2. Browse the [SODA Dock GitHub Issues](https://github.com/sodafoundation/dock/issues) labeled with tags like ‘good first issue’, ‘help needed’, ‘StartMyContribution’, etc.

## Roadmap

We envision a vast support base for all storage vendor drivers under the dock, creating a standardized and unified storage backend interface. More about the roadmap can be found [here](https://docs.sodafoundation.io/).

## Connect with SODA Foundation

- **Website**: [SODA Foundation](https://sodafoundation.io/)
  
- **Slack**: [Join the Conversation](https://sodafoundation.io/slack/)
  
- **Twitter**: [Follow @sodafoundation](https://twitter.com/sodafoundation)
  
- **Mailing List**: [SODA Foundation Mailing List](https://lists.sodafoundation.io/)