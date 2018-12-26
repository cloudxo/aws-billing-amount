# aws-billing-amount

## ローカルで動作確認をするとき

```go
func main() {
	// Run local
	Run() // <- コメントを外す

	// Run lambda
	// lambda.Start(Run)
}
```

## AWS Lambda に反映する

### 0. コメントを戻す

### 1. ビルド

```bash
$ GOOS=linux GOARCH=amd64 go build -o aws-billing-amount
```

### 2. ZIP 化

```
$ zip aws-billing-amount.zip ./aws-billing-amount 
```

### 3. アップロード

![lambda_management_console](https://user-images.githubusercontent.com/5449002/50445360-934feb00-0951-11e9-9195-583301d8103b.png)
