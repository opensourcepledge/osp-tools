<!--
© 2024 Vlad-Stefan Harbuz <vlad@vlad.website>

SPDX-License-Identifier: CC-BY-SA-4.0
-->

# gh-sponsor-finder

Gets all GitHub sponsor organizations of a certain degree starting from an
organization.

For example, all sponsor organizations of degree 2 in Sentry's network are all
sponsors who sponsor users sponsored by organizations which sponsor users which
Sentry sponsors.

This works by starting at an organization, getting the users it sponsors,
getting all organizations that sponsor those users, getting all users those
organizations sponsor, and so on.

## Usage

```
GH_TOKEN="github_..." ./main.py --start-org getsentry
```
