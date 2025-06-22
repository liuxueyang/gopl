package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

var dir_name *string = flag.String("dir", "", "project name")
var bin *bool = flag.Bool("bin", false, "create binary file")

func main() {
	flag.Parse()

	if *dir_name == "" {
		panic("Please specify a project name with -dir flag")
	}

	err := os.MkdirAll(*dir_name, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", *dir_name, err)
		os.Exit(1)
	}

	if *bin {
		err = createBinaryFile(*dir_name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating binary file: %v\n", err)
			os.Exit(1)
		}
	} else {
		libName := *dir_name
		err = createLibraryFile(*dir_name, libName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating library file: %v\n", err)
			os.Exit(1)
		}
	}

	// change to the project directory
	err = os.Chdir(*dir_name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error changing to directory %s: %v\n", *dir_name, err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "mod", "init", *dir_name)
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Go module: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Project %s created successfully\n", *dir_name)
}

const binaryFile = `package main

func main() {

}
`

const libraryFile = `package {{.LibName}}

`

var binaryTemplate = template.Must(template.New("binary").Parse(binaryFile))
var libraryTemplate = template.Must(template.New("library").Parse(libraryFile))

func createBinaryFile(dir string) error {
	file, err := os.Create(dir + "/main.go")
	if err != nil {
		return fmt.Errorf("error creating main.go: %w", err)
	}
	defer file.Close()

	err = binaryTemplate.Execute(file, nil)
	if err != nil {
		return fmt.Errorf("error executing template for main.go: %w", err)
	}

	return nil
}

func createLibraryFile(dir, libName string) error {
	file, err := os.Create(fmt.Sprintf("%s/%s.go", dir, libName))
	if err != nil {
		return fmt.Errorf("error creating %s.go: %w", libName, err)
	}
	defer file.Close()

	data := struct {
		LibName string
	}{
		LibName: libName,
	}

	err = libraryTemplate.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template for %s.go: %w", libName, err)
	}

	return nil
}
