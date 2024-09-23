init:
  #!/usr/bin/env bash
  git clone --filter=blob:none --no-checkout --depth 1 --sparse git@github.com:beyond-all-reason/Beyond-All-Reason.git bar-repo
  cd bar-repo
  git sparse-checkout set --no-cone "units" "language/en" "luaui/configs" # "unitpics"
  git sparse-checkout list
  git checkout

find-by-name name:
  cat bar-repo/language/en/units.json | jq '.units.names | to_entries | .[] | select(.value=="{{name}}") | .key'\

file key:
  find bar-repo/units -name "{{key}}.lua";

data key:
  cat `just file {{key}}`

run:
  go run *.go
