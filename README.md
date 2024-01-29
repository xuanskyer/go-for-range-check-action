# go-for-range-check-action

This action will recursively traverse the golang files in your specified directory, scan the functions and methods inside,

parse the nesting level of the for statement (including the for range statement) in each function and method,

and then determine whether there is a list of functions and methods that exceed the specified nesting level

based on the specified valid nesting parameters

# Usage

<!-- start usage -->

```yaml
- uses: xuanskyer/go-for-range-check-action@v1
  with:
    # Maximum for loop nesting level
    # Default: 3
    max-for-range-level: '3'

    # scan target directory
    target-directory: ''

    # ignore directory
    ignore-directory: ''
```

<!-- end usage -->

## demo

```yaml
name: Test Go For Range Check Action
on:
  push:
    branches:
      - main
jobs:
  test:
    runs-on: ubuntu-latest
    defaults:
      run:
        shell: bash
    steps:
      - name: Checkout repository
        uses: actions/checkout@v2

      - name: Go For Range Check Action
        uses: ./
        with:
          max-for-range-level: 3
          target-directory: 'biz/'  # set scan target directory
          ignore-directory: '[\"vendor\"]'
```

# License
