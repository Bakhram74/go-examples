if ! (which mockgen &> /dev/null)
then
  echo "mockgen not found, trying to install it..."
  go install go.uber.org/mock/mockgen@latest
fi
mockgen -source ./internal/usecase/interfaces.go -package usecase_test > ./internal/usecase/mocks_test.go

echo "Mocks generated successfully"
