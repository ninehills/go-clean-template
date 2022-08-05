#!/usr/bin/env python3
"""Custom inteface inject.
"""

content = ""

with open('internal/dao/querier.go', 'r') as f:
    for line in f.readlines():
        content += line
        if "UpdateUser" in line:
            content += '	QueryUser(ctx context.Context, arg QueryUserParams) ([]User, int64, error)\n'

with open('internal/dao/querier.go', 'w') as f:
    f.write(content)
