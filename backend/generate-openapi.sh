if ! (which npm &> /dev/null)
then
  echo "Node.js not installed. install Node.js and retry"
  exit 1
fi

if ! (which oapi-codegen &> /dev/null)
then
  echo "oapi-codegen not found, trying to install it..."
  go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
fi

if ! (which redocly &> /dev/null)
then
  echo "redocly/cli not found, trying to install it..."
  npm install -g @redocly/cli
fi

find internal/controller/http/v1/ internal/entity/ -type f -name '*.gen.*' | xargs rm -f
redocly bundle openapi/openapi-config.yaml -o openapi/global-openapi.yaml
oapi-codegen -config openapi/server-openapi.cfg.yaml openapi/global-openapi.yaml
oapi-codegen -config openapi/entity-openapi.cfg.yaml openapi/global-openapi.yaml

echo "Open Api generated successfully."