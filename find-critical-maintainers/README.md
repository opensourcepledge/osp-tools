<!--
© 2025 Functional Software, Inc. d/b/a Sentry
SPDX-License-Identifier: Apache-2.0
-->

# find-critical-maintainers

This tool finds the most depended upon maintainers across multiple Open Source ecosystems, using data from
[ecosyste.ms][ecosyste.ms].

Here's how this broadly works:

1. ecosyste.ms collects data about packages across multiple ecosystems (npm, PyPI and so on), including package download
   data.
2. We find critical packages for each ecosystem — the packages that, when their downloads are summed up, they together
   make up 80% of the downloads for that ecosystem. See the [origin of this calculation][calc-origin] and the
   [calculation source code][calc-code].
3. We find the maintainers that maintain those packages.

This method has multiple limitations:

* It only handles source-level dependencies, such as in manifests like `package.json`. It does not handle binary
  dependencies, so misses out the maintainers of Linux, PostgreSQL, build tools and so on.
* It cannot identify human maintainers when packages are maintained by an organisation account, such as
  [rubygems/awscloud](https://rubygems.org/profiles/awscloud).
* It overrepresents small ecosystems, such as NuGet.
* Only some ecosystems are supported. Only ecosystems that publish per-package download statistics can be supported. For
  example, it does not seem possible to get download statistics for Go packages.

There are some improvements still to make:

* Consider commit authors (eg using ecosyste.ms summary service).

Useful ecosyste.ms resources:

* Get all critical packages<br>
  https://packages.ecosyste.ms/api/v1/packages/critical?page=1&per_page=1000
* Get all critical packages for a registry<br>
  https://packages.ecosyste.ms/api/v1/registries/rubygems.org/packages?critical=true&page=1&per_page=1000
* Get all registries<br>
  https://packages.ecosyste.ms/api/v1/registries/
* Look up a project's summary by repository URL (particularly, see committers)<br>
  https://summary.ecosyste.ms/api/v1/projects/lookup?url=https://github.com/pythonadelaide/karmabot
* Get data about GitHub Sponsors accounts and their sponsors<br>
  https://sponsors.ecosyste.ms/docs/index.html
* Get high-level info about a user (eg funding links, total stars etc)<br>
  https://repos.ecosyste.ms/api/v1/hosts/Github/owners/andrew

[calc-code]: https://github.com/ecosyste-ms/packages/blob/81886822f53f32b31a61d65505a62e5cb354038c/app/models/registry.rb#L415
[calc-origin]: https://github.com/chadwhitacre/openpath/issues/20#issuecomment-1929436690
[ecosyste.ms]: https://ecosyste.ms
