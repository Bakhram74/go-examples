if ! (which mockgen &> /dev/null)
then
  echo "mockgen not found, trying to install it..."
  go install go.uber.org/mock/mockgen@latest
fi
mkdir -p ./pkg/accessor/accessor_test
mockgen -source ./pkg/accessor/accessor.go -package accessortest > ./pkg/accessor/accessor_test/mocks.go

echo "Mocks generated successfully"
