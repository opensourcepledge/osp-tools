#!/usr/bin/env python3

# © 2024 Vlad-Stefan Harbuz <vlad@vladh.net>
#
# SPDX-License-Identifier: Apache-2.0

import argparse
import os
import sys

from gql import gql, Client
from gql.transport.aiohttp import AIOHTTPTransport


USER_SPONSOR_QUERY = gql("""
query getUserSponsors($login: String!, $after: String) {
    user(login: $login) {
        sponsors(first: 100, after: $after) {
            nodes {
                ... on Organization {
                    login
                }
            }
            pageInfo {
                endCursor
                hasNextPage
            }
        }
    }
}
""")

ORG_SPONSORING_QUERY = gql("""
query getOrgSponsoring($login: String!, $after: String) {
    organization(login: $login) {
        sponsoring(first: 100, after: $after) {
            nodes {
                ... on User {
                    login
                }
            }
            pageInfo {
                endCursor
                hasNextPage
            }
        }
    }
}
""")

ORG_SPONSORING_COUNT_QUERY = gql("""
query getOrgSponsoringCount($login: String!) {
    organization(login: $login) {
        sponsoring(first: 1) {
            totalCount
        }
    }
}
""")


"""
Prints to stderr.
"""
def eprint(*args, **kwargs):
    print(*args, file=sys.stderr, **kwargs)


"""
Gets a `gql.Client`.
"""
def get_gql_client():
    token = os.environ.get('GH_TOKEN', '')
    transport = AIOHTTPTransport(url="https://api.github.com/graphql",
        headers={'Authorization': f'bearer {token}'})
    return Client(transport=transport, fetch_schema_from_transport=True)


"""
Gets all organizations which sponsor user with login `name`.
"""
def get_user_sponsors(client, name):
    eprint(f"get_user_sponsors(_, {name})")
    sponsors = set()

    after = None
    while True:
        page_results = client.execute(USER_SPONSOR_QUERY,
            variable_values={'login': name, 'after': after})
        for sponsor in page_results['user']['sponsors']['nodes']:
            if sponsor:
                sponsors.add(sponsor['login'])
        after = page_results['user']['sponsors']['pageInfo']['endCursor']
        if not page_results['user']['sponsors']['pageInfo']['hasNextPage']:
            break

    return sponsors


"""
Gets all users sponsored by an organization with login `name`.
"""
def get_org_sponsored_users(client, name):
    eprint(f"get_org_sponsored_users(_, {name})")
    sponsored_users = set()

    after = None
    while True:
        page_results = client.execute(ORG_SPONSORING_QUERY,
            variable_values={'login': name, 'after': after})
        for sponsored_user in page_results['organization']['sponsoring']['nodes']:
            if sponsored_user:
                sponsored_users.add(sponsored_user['login'])
        after = page_results['organization']['sponsoring']['pageInfo']['endCursor']
        if not page_results['organization']['sponsoring']['pageInfo']['hasNextPage']:
            break

    return sponsored_users


"""
Gets the number of users sponsored by the organization with login `name`.
"""
def get_org_sponsored_user_count(client, name):
    eprint(f"get_org_sponsored_user_count(_, {name})")
    sponsored_users = set()

    page_results = client.execute(ORG_SPONSORING_COUNT_QUERY,
        variable_values={'login': name})

    return page_results['organization']['sponsoring']['totalCount']


"""
Gets all GitHub sponsors of `degree` starting from organization with login
`name`.

For example, all sponsors of degree 2 in Sentry's network are all organizations
who sponsor users sponsored by organizations which sponsor users which Sentry
sponsors.
"""
def get_sponsor_network(client, start_org, degree):
    seen_users = set()
    seen_sponsors = set()
    sponsors = set([start_org])

    for idx_iteration in range(0, degree):
        sponsored_users = set()

        new_sponsors = sponsors.difference(seen_sponsors)
        for idx, new_sponsor in enumerate(new_sponsors):
            eprint(f"{idx + 1}/{len(new_sponsors)}")
            try:
                new_sponsored_users = get_org_sponsored_users(client, new_sponsor)
                sponsored_users = sponsored_users.union(new_sponsored_users)
            except Exception as e:
                eprint(f"Couldn't get sponsored users for org {new_sponsor}. Ignoring.", e)

        seen_sponsors = seen_sponsors.union(sponsors)

        new_users = sponsored_users.difference(seen_users)
        for idx, new_user in enumerate(new_users):
            eprint(f"{idx + 1}/{len(new_users)}")
            try:
                user_sponsors = get_user_sponsors(client, new_user)
                sponsors = sponsors.union(user_sponsors)
            except Exception as e:
                eprint(f"Couldn't get sponsors for user {new_user}. Ignoring.", e)

        seen_users = seen_users.union(sponsored_users)

    return sponsors


def main():
    parser = argparse.ArgumentParser("gh_sponsor_finder")
    parser.add_argument("--start-org",
        help="The organization to start the search from",
        type=str,
        required=True)
    args = parser.parse_args()

    client = get_gql_client()

    sponsors = get_sponsor_network(client, args.start_org, degree=2)
    for idx, sponsor in enumerate(sponsors):
        eprint(f"{idx + 1}/{len(sponsors)}")
        count = get_org_sponsored_user_count(client, sponsor)
        print(f'* [{sponsor}](https://github.com/{sponsor}), sponsoring {count}')


if __name__ == '__main__':
    main()
