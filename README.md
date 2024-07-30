# Rhizom pg-null

A simple package for use attributes as a SQL null

## Building from source

To build from source, you will need the following prerequisites:

- Go 1.14 or greater;
- Git

### Downloading the code

First, clone the project:

```bash
git clone git@github.com:plateausnetwork/pg-null.git /your/directory/of/choice/rhizom
cd /your/directory/of/choice/rhizom
```

### Using as a library

```go
import (
  null "github.com/plateausnetwork/pg-null"
)

type myStruct struct {
  Name        null.String `json:"name,omitempty"`
	Description null.String `json:"description,omitempty"`
}

```

## License

For more details about our license model, please take a look at the [LICENSE](LICENSE) file.

**2020**, Rhizom Platform.
