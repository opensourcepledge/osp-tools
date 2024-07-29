#!/usr/env/bin python3

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


def get_gql_client():
    token = os.environ.get('GH_TOKEN', '')
    transport = AIOHTTPTransport(url="https://api.github.com/graphql",
        headers={'Authorization': f'bearer {token}'})
    return Client(transport=transport, fetch_schema_from_transport=True)


def get_user_sponsors(client, name):
    sponsors = []

    after = None
    while True:
        pageResults = client.execute(USER_SPONSOR_QUERY,
            variable_values={'login': name, 'after': after})
        for sponsor in pageResults['user']['sponsors']['nodes']:
            if sponsor:
                sponsors.append(sponsor)
        after = pageResults['user']['sponsors']['pageInfo']['endCursor']
        if not pageResults['user']['sponsors']['pageInfo']['hasNextPage']:
            break

    return sponsors


def main():
    parser = argparse.ArgumentParser("sponsor_finder")
    parser.add_argument("--username",
        help="The user to get sponsors for",
        type=str,
        required=True)
    args = parser.parse_args()

    client = get_gql_client()
    sponsors = get_user_sponsors(client, args.username)
    for sponsor in sponsors:
        print(f"https://github.com/{sponsor['login']}")


if __name__ == '__main__':
    main()
