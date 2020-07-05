#!/bin/env bash

set -euo pipefail
IFS=$'\n\t'

# TEST 1. request the Genesis block
curl -s "localhost:8000/block/0" | jq
echo

# TEST 2. Make request of ownership
curl -s -X POST \
  -H 'Content-Type: application/json' \
  -d '{"address":"1CAvNmCrxSRympnSoVxYKLuXdDthyB74xu"}' \
  localhost:8000/requestValidation
echo

# TEST 3. Submit your Star
curl -s -X POST -H 'Content-Type: application/json' localhost:8000/submitStar -d @- <<\EOF | jq
  { "address": "1CAvNmCrxSRympnSoVxYKLuXdDthyB74xu",
    "signature": "H8MtnshsXv4Aw1VVGZKAyhKLyya9ebYyMnLgTW13B7aAILqQNiaHox28vsLok39Zf36msVEFWQoAj7stPSJ6yIQ=",
    "message": "1CAvNmCrxSRympnSoVxYKLuXdDthyB74xu:1590921500:starRegistry",
    "star": {
      "dec": "68° 52' 56.9",
      "ra": "16h 29m 1.0s",
      "story": "Testing the story 5"
    }
  }
EOF
echo

# TEST 4. Retrieve Stars owned by me
curl -s localhost:8000/blocks/1CAvNmCrxSRympnSoVxYKLuXdDthyB74xu | jq
echo

# TEST 5. Get block by hash
curl -s localhost:8000/block/hash/b06ad471a19ef484b8d26fc4bc9255aca274239e8395a188d8acfed1f97d0206 | jq
echo

