package stores

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/foxinuni/citas/core/models"
)

type CitaStoreFilter struct {
	Date  time.Time
	Limit int
	Page  int
}

type CitaStore interface {
	GetAll(filter CitaStoreFilter) ([]models.Cita, error)
	GetById(id string) (models.Cita, error)
	Create(cita *models.Cita) error
	Update(cita *models.Cita) error
	Delete(cita *models.Cita) error
}

type FsPath string

var (
	ErrCitaNotFound = errors.New("failed to find cita with the specified ID")
	ErrInvalidId    = errors.New("invalid ID format")
	ErrCitaExists   = errors.New("cita already exists on this date")
)

type InFsCitaStore struct {
	path string
}

func NewInFsCitaStore(relpath FsPath) (CitaStore, error) {
	// Se obtine la ruta absoluta del directorio
	path, err := filepath.Abs(string(relpath))
	if err != nil {
		return nil, err
	}

	// Se verifica si el directorio no existe, se crea
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return &InFsCitaStore{path: path}, nil
}

/*
 * DateToFolderName convierte una fecha a un nombre de carpeta,
 * en formato YYYYMMDD. Esto se usa para almacenar las citas en carpetas
 */
func (s *InFsCitaStore) DateToFolderName(date time.Time) string {
	return date.Format("20060102")
}

/*
 * Los IDs son generadas a partir de la ubicación de la cita en el sistema de archivos.
 * Esta es la ruta relativa de la cita en el sistema de archivos, codificada en base64.
 * GenerateIdForCita genera un ID a partir de una cita.
 */
func (s *InFsCitaStore) GenerateIdForCita(cita *models.Cita) string {
	folderName := s.DateToFolderName(cita.Fecha)
	cedula := cita.Persona.Cedula

	decoded := filepath.Join(folderName, cedula+".json")

	// encode base64
	encoded := base64.StdEncoding.EncodeToString([]byte(decoded))
	return encoded
}

/*
 * GetFilePathForCitaId obtiene la ruta absoluta de un archivo de cita a partir de un ID.
 * El ID es decodificado de base64 y se concatena con la ruta base del sistema de archivos.
 */
func (s *InFsCitaStore) GetFilePathForCitaId(id string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return "", err
	}

	return filepath.Join(s.path, string(decoded)), nil
}

func (s *InFsCitaStore) ParseCitaFromFile(filepath string) (models.Cita, error) {
	reader, err := os.Open(filepath)
	if err != nil {
		return models.Cita{}, ErrCitaNotFound
	}
	defer reader.Close()

	cita := models.Cita{}
	if err := json.NewDecoder(reader).Decode(&cita); err != nil {
		return models.Cita{}, err
	}

	return cita, nil
}

func (s *InFsCitaStore) GetAll(filter CitaStoreFilter) ([]models.Cita, error) {
	folderName := s.DateToFolderName(filter.Date)
	folderPath := filepath.Join(s.path, folderName)

	// Se verifica si el directorio no existe
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return []models.Cita{}, nil
	}

	// Se obtiene la lista de archivos del directorio
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	// Se lee cada archivo y se convierte a una cita
	citas := []models.Cita{}
	for _, file := range files {
		// Se ignoran los archivos que no sean .json
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Se lee el archivo y se convierte a una cita
		cita, err := s.ParseCitaFromFile(filepath.Join(folderPath, file.Name()))
		if err != nil {
			return nil, err
		}

		citas = append(citas, cita)
	}

	return citas, nil
}

func (s *InFsCitaStore) GetById(id string) (models.Cita, error) {
	filepath, err := s.GetFilePathForCitaId(id)
	if err != nil {
		return models.Cita{}, ErrInvalidId
	}

	return s.ParseCitaFromFile(filepath)
}

func (s *InFsCitaStore) Create(cita *models.Cita) error {
	cita.ID = s.GenerateIdForCita(cita)
	folderPath := filepath.Join(s.path, s.DateToFolderName(cita.Fecha))
	
	filePath, err := s.GetFilePathForCitaId(cita.ID)
	if err != nil {
		return err
	}

	/*
	 * Volveria esto su propia funcion, pero creo que ya tengo muchas.
	 * Realmente que la creacion solo se hace una vez, asi que no es muy necesario.
	 */
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return ErrCitaExists
	}

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return err
	}

	writer, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer writer.Close()

	if err := json.NewEncoder(writer).Encode(cita); err != nil {
		return err
	}

	return nil
}

func (s *InFsCitaStore) Update(cita *models.Cita) error {
	panic("not implemented")
}

func (s *InFsCitaStore) Delete(cita *models.Cita) error {
	panic("not implemented")
}
