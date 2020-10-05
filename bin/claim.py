#!/usr/bin/env python3

import os
import json
from sys import stdout

from argparse import ArgumentParser
from jsonschema import validate

def main():
    '''Validate an input claim's structure against the JSON Schema for claims.
       Print the claim to stdout if it is valid.
    '''
    aparser = ArgumentParser(description=main.__doc__)
    aparser.add_argument(
        '-s', '--schema',
        default=os.path.join(
            os.path.dirname(__file__),
            '..', 'share', 'claim.json',
        ),
        help="the JSON Schema to use to validate the claim's structure",
    )
    aparser.add_argument(
        'filename',
        help="the filename containing the claim or '-' to read from stdin",
    )
    args = aparser.parse_args()
    with open(args.schema) as fid:
        print(args.schema)
        schema = json.load(fid)
    with open('/dev/stdin' if args.filename == '-' else args.filename) as fid:
        claim = json.load(fid)
    validate(claim, schema)
    json.dump(claim, stdout)

if __name__ == '__main__':
    main()
