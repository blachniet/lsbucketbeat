# lsbucketbeat

`lsbucketbeat` sends information about files in a directory to Elasticsearch.

## Getting Started

`lsbucketbeat` groups files into **buckets**. A bucket is given a title, points to a directory and has a file pattern associated with it. For example:

```yaml
lsbucketbeat:
  period: 1m
  buckets:
    - title: "My Markdown Files"
      dir: "/Users/johndoe/Documents"
      filePattern: "*.md"
      retryCount: 3
      retryDelay: 1s
```

This configuration defines a single bucket. Every 1 minute, the `/Users/johndoe/Documents` directory is scanned for files matching the `*.md` pattern. A document is created for each file matching this pattern.

```json
{
    "@timestamp": "2017-04-23T13:09:59.922Z",
    "beat": {
        "hostname": "Johns-Air.lan",
        "name": "Johns-Air.lan",
        "version": "5.3.1"
    },
    "bucket": {
        "title": "My Markdown Files"
    },
    "file": {
        "dir": "/Users/johndoe/Documents",
        "modTime": "2017-04-10T22:27:48.000Z",
        "name": "README.md",
        "path": "/Users/johndoe/Documents/README.md",
        "size": 234
    },
    "type": "lsbucketbeat"
}
```

See [docs/fields.asciidoc](docs/fields.asciidoc) for more details on the structure of the events.

You may have noticed the `retryCount` and `retryDelay` configuration properties in our example. These define the parameters for retrying if there are any errors while trying to scan the target directory and files.

You may define as many buckets as you like in your configuration. The `bucket.title` property is the best way to differentiate between buckets in Elasticsearch.
