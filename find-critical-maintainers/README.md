<!--
© 2025 Functional Software, Inc. d/b/a Sentry
SPDX-License-Identifier: Apache-2.0
-->

# find-critical-maintainers

This tool finds the most depended upon maintainers across multiple Open Source ecosystems, using data from
[ecosyste.ms][ecosyste.ms].

## Explanation

Here's how this broadly works:

1. ecosyste.ms collects data about packages across multiple ecosystems (npm, PyPI and so on), including package download
   data.
2. We find critical packages for each ecosystem — the packages that, when their downloads are summed up, they together
   make up 80% of the downloads for that ecosystem. See the [origin of this calculation][calc-origin] and the
   [calculation source code][calc-code].
3. We find the maintainers that maintain those packages. Sometimes ecosystems can tell us who maintains a package. But
   sometimes it's more complicated.
4. We try to find each package's repository, and go through its commit statistics. We find the committers with the
   highest commit count who, when summed up, make up at least 30% of the repo's commits. We call these “significant
   committers”, and include them in the maintainer list.
5. We can also use GitHub issues to identify maintainers. GitHub issue comments can bear specific badges, such as
   “Collaborator” or “Member”, which signify that the user who left that commit has been granted certain permissions on
   that repository. [We add Members, Owners and Collaborators][member-assoc] to the maintainer list for that package.

## Limitations

This method has multiple limitations:

* It only handles source-level dependencies, such as in manifests like `package.json`. It does not handle binary
  dependencies, so misses out the maintainers of Linux, PostgreSQL, build tools and so on.
* It cannot identify human maintainers when packages are maintained by an organisation account, such as
  [rubygems/awscloud](https://rubygems.org/profiles/awscloud).
* It overrepresents small ecosystems, such as NuGet.
* Only some ecosystems are supported. Only ecosystems that publish per-package download statistics can be supported. For
  example, it does not seem possible to get download statistics for Go packages.
* Data is incomplete. Some critical packages have no committer information and no commit data, so maintainers cannot be
  identified.
* Some heuristics are used, which can lead to incorrect data. For example, maintainers are deduplicated by merging
  maintainer data where two ostensibly different maintainers have the same username, name or email address. This should
  mostly yield correct results, but has the potential to merge unrelated maintainers with identical names, for example.

## TODOs

* Implement local cache for ecosyste.ms endpoints. This would allow us to rerun revisions of the code without hitting
  the network every time. The cache could be cleared whenever fresh data is desired.
* Move this tool into its own repository. Maybe with a spiffier name?
* Write a case study blog post for [the ecosyste.ms blog][eco-blog].
* Write a case study for [the ecosyste.ms docs website][eco-docs].
* Write a blog post for [the Open Source Pledge blog][osp-blog].

## Useful ecosyste.ms resources

* Get all critical packages<br>
  https://packages.ecosyste.ms/api/v1/packages/critical?page=1&per_page=1000
* Get all critical packages for a registry<br>
  https://packages.ecosyste.ms/api/v1/registries/rubygems.org/packages?critical=true&page=1&per_page=1000
* Get all registries<br>
  https://packages.ecosyste.ms/api/v1/registries/
* Look up a project's summary by repository URL (particularly, see committers and issues)<br>
  https://summary.ecosyste.ms/api/v1/projects/lookup?url=https://github.com/pythonadelaide/karmabot
* Get data about GitHub Sponsors accounts and their sponsors<br>
  https://sponsors.ecosyste.ms/docs/index.html
* Get high-level info about a user (eg funding links, total stars etc)<br>
  https://repos.ecosyste.ms/api/v1/hosts/Github/owners/andrew

## Contact

This tool is maintained by [Vlad-Stefan Harbuz](https://vlad.website) — contact him with questions and suggestions.

[calc-code]: https://github.com/ecosyste-ms/packages/blob/81886822f53f32b31a61d65505a62e5cb354038c/app/models/registry.rb#L415
[calc-origin]: https://github.com/chadwhitacre/openpath/issues/20#issuecomment-1929436690
[eco-blog]: https://blog.ecosyste.ms/
[eco-docs]: https://docs.ecosyste.ms/]
[ecosyste.ms]: https://ecosyste.ms
[member-assoc]: https://github.com/ecosyste-ms/issues/blob/9678d9cba89f054f2114969f96d751c371ed0c33/app/models/issue.rb#L22
[osp-blog]: https://opensourcepledge.com/blog/
