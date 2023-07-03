export PORT=${PG_PORT:-5433}

if ! which postgres > /dev/null; then
  version=$(ls -1 /usr/lib/postgresql/ | tail -n 1)
  if [[ "$version" == "" ]]; then
    echo "can't find postgres executable path."
    exit 1
  fi

  export PATH="/usr/lib/postgresql/$version/bin:$PATH"
fi
