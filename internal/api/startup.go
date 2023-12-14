package startup

import (
	"html/template"
	repo "lab3/internal/infra/impl"
	service "lab3/internal/service/impl"
	"net/http"
	"os"
	"path/filepath"
)

const uploadFolder = "./uploads"
const templateDir = "./templates"

func Run() {

	repository, err := repo.NewMongoDBRepository("files")
	if err != nil {
		panic(err)
	}

	fileService := service.NewFileService(repository, uploadFolder)

	http.HandleFunc("/", indexHandler(fileService))
	http.HandleFunc("/upload", uploadHandler(fileService))
	http.HandleFunc("/download/", downloadHandler(fileService))

	if _, err := os.Stat(uploadFolder); os.IsNotExist(err) {
		err := os.Mkdir(uploadFolder, os.ModePerm)
		if err != nil {
			return
		}
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func indexHandler(fileService *service.FileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index.html", nil)
	}
}

func uploadHandler(fileService *service.FileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "No file provided", http.StatusBadRequest)
				return
			}

			defer file.Close()

			fileService.UploadFile(file, header.Filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			renderTemplate(w, "upload.html", nil)
		}
	}
}

func downloadHandler(fileService *service.FileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileReferenceID := r.RequestURI[len("/download/?fileReferenceID="):]
		filePath, err := fileService.DownloadFile(fileReferenceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, filePath)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := filepath.Join(templateDir, tmpl)
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
