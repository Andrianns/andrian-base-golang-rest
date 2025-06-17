package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func main() {
	pkg := flag.String("pkg", "", "Nama package/project yang akan dibuat")
	flag.Parse()

	if *pkg == "" {
		fmt.Println("Gunakan: --pkg=nama-package")
		os.Exit(1)
	}

	fmt.Println("Membuat project:", *pkg)
	dirs := []string{
		"app/config",
		"app/client",
		"app/repository",
		"app/controller",
		"app/routes",
		"app/models",
	}

	for _, dir := range dirs {
		os.MkdirAll(filepath.Join(*pkg, dir), os.ModePerm)
	}

	generateFromTemplate("templates/main.go.txt", filepath.Join(*pkg, "main.go"), *pkg)
	generateFromTemplate("templates/app/controller/controller.go.txt", filepath.Join(*pkg, "app", "controller", "controller.go"), *pkg)
	generateFromTemplate("templates/app/config/config.go.txt", filepath.Join(*pkg, "app", "config", "config.go"), *pkg)
	generateFromTemplate("templates/app/config/connection.go.txt", filepath.Join(*pkg, "app", "config", "connection.go"), *pkg)
	generateFromTemplate("templates/app/client/google_drive_client.go.txt", filepath.Join(*pkg, "app", "client", "google_drive_client.go"), *pkg)
	generateFromTemplate("templates/app/controller/controller.go.txt", filepath.Join(*pkg, "app", "controller", "controller.go"), *pkg)
	generateFromTemplate("templates/app/controller/document_controller.go.txt", filepath.Join(*pkg, "app", "controller", "document_controller.go"), *pkg)
	generateFromTemplate("templates/app/models/user_model.go.txt", filepath.Join(*pkg, "app", "models", "user_model.go"), *pkg)
	generateFromTemplate("templates/app/repository/user_repository.go.txt", filepath.Join(*pkg, "app", "repository", "user_repository.go"), *pkg)
	generateFromTemplate("templates/app/routes/route.go.txt", filepath.Join(*pkg, "app", "routes", "route.go"), *pkg)

	goMod := fmt.Sprintf("module %s\n\ngo 1.24.3", *pkg)
	err := os.WriteFile(filepath.Join(*pkg, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		log.Fatalf("Gagal menulis go.mod: %v", err)
	}

	copyFile("templates/env-example.txt", filepath.Join(*pkg, "env-example.txt"))
	copyFile("templates/env-example.txt", filepath.Join(*pkg, ".env"))
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = *pkg
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Gagal menjalankan go mod tidy: %v", err)
	}
	fmt.Println("Sukses membuat base project:", *pkg)
}

func generateFromTemplate(src, dst, pkg string) {
	content, err := os.ReadFile(src)
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("file").Parse(string(content)))
	f, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := struct {
		Package string
	}{
		Package: pkg,
	}

	t.Execute(f, data)
}

func copyFile(src, dst string) {
	fmt.Println(src, dst)
	content, err := os.ReadFile(src)
	if err != nil {
		log.Fatalf("Gagal membaca file %s: %v", src, err)
	}
	err = os.WriteFile(dst, content, 0644)
	if err != nil {
		log.Fatalf("Gagal menulis file %s: %v", dst, err)
	}
}
