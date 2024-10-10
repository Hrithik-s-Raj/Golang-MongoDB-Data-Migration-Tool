
# MongoDB Database Migration Script

This Go script migrates data from a source MongoDB database to a target MongoDB database.

## Configuration

Edit the `sourceURI`, `targetURI`, `sourceDBName`, and `targetDBName` variables in the script:

```go
sourceURI := "mongodb://source-uri"
targetURI := "mongodb://target-uri"
sourceDBName := "source_db"
targetDBName := "target_db"
```

## Installation

1. Install the MongoDB Go driver:
   ```bash
   go get go.mongodb.org/mongo-driver/mongo
   ```

2. Run the script:
   ```bash
   go run main.go
   ```

## Troubleshooting

### Context Deadline Exceeded Error

- **Increase Timeout**: Adjust the timeout for large datasets.
- **Reduce Batch Size**: Process data in smaller chunks.
- **Check Network**: Ensure stable connections between the app and MongoDB.

## Author

- **Hrithik S Raj** - https://github.com/Hrithik-s-Raj

## License

Licensed under the [MIT License](LICENSE).
