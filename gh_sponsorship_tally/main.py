#!/usr/env/bin python3

import argparse
import os
import sys

from gql import gql, Client
from gql.transport.aiohttp import AIOHTTPTransport


USER_SPONSORSHIP_VALUES_QUERY = gql("""
query getUserSponsorshipValues($login: String!, $after: String) {
    user(login: $login) {
        lifetimeReceivedSponsorshipValues(first: 100, after: $after) {
            nodes {
                amountInCents
                formattedAmount
                sponsor {
                    ... on User { login }
                    ... on Organization { login }
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
def get_user_sponsorship_values(client, name):
    eprint(f"get_user_sponsorship_values(_, {name})")

    after = None
    while True:
        page_results = client.execute(USER_SPONSORSHIP_VALUES_QUERY,
            variable_values={'login': name, 'after': after})
        print(page_results)
        after = page_results['user']['lifetimeReceivedSponsorshipValues']['pageInfo']['endCursor']
        if not page_results['user']['lifetimeReceivedSponsorshipValues']['pageInfo']['hasNextPage']:
            break


def main():
    parser = argparse.ArgumentParser("gh_sponsor_finder")
    parser.add_argument("--org",
        help="The organization to get sponsorship amounts for",
        type=str,
        required=True)
    args = parser.parse_args()

    client = get_gql_client()

    get_user_sponsorship_values(client, args.org)


if __name__ == '__main__':
    main()
