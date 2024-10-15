package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

const (
	// Путь к файлу куда падают билеты от поставщиков
	pathToFile = "C:\\FTP\\FTPForTickets"
	// Путь куда сохранять файлы
	pathToArchive = "C:\\FTP\\Archive"
	// Путь к файлу с конфигом подключения к FTP
	configPath = "configFTP.yaml"
)

// FTPConfig Структура для хранения данных локальных FTP
type FTPConfig struct {
	Name    string
	Address string
	User    string
	Pass    string
}

// SuppliersFTPConfig Структура для хранения данных FTP/SFTP поставщиков
type SuppliersFTPConfig struct {
	Name        string
	Address     string
	TypeConnect string
	User        string
	Pass        string
	SaveIn      string
	Path        []string
}

// Config Структура для хранения данных о всех FTP которые используются
type Config struct {
	FTPConfigs         map[string]FTPConfig          `yaml:"localFtpConfigs"`
	SuppliersFTPConfig map[string]SuppliersFTPConfig `yaml:"suppliersFtpConfigs"`
}

// foldersSuppliers Массив с папками поставщиков
var foldersSuppliers = [...]string{
	"Aerotur", //Мы прыгаем на их FTP
	"LocalUser\\LinerAmadeus",
	"LocalUser\\LinerB2BB",
	"LocalUser\\LinerSSOD",
	"LocalUser\\MyAgent",
	"OLTC_02aka_ppr02",
	"OLTC_02aka_ppr05",
	"OLTC_27moa_ppr34",
	"OLTC_52msa_ppr01",
	"OLTC_71mos_ppr25",
	"OLTC_87mov_ppr05",
	"OLTC_agn_02mok_main",
	"OLTC_agn_05msv_main",
	"OLTC_agn_u6173_main",
	"S7Smart", //Мы прыгаем на их SFTP На 1С эта папка называется S7
	"Sabre",
	"LocalUser\\TTBooking",
	"VipServiceAmadeus",   //Мы прыгаем на их FTP
	"VipServicePortbilet", //Мы прыгаем на их FTP
}

// foldersSuppliers Массив с папками поставщиков
var foldersSave = [...]string{
	"Aerotur", //Мы прыгаем на их FTP
	"LinerAmadeus",
	"LinerB2B",
	"LinerSSOD",
	"TTT",
	"OLTC_02aka_ppr02",
	"OLTC_02aka_ppr05",
	"OLTC_27moa_ppr34",
	"OLTC_52msa_ppr01",
	"OLTC_71mos_ppr25",
	"OLTC_87mov_ppr05",
	"OLTC_agn_02mok_main",
	"OLTC_agn_05msv_main",
	"OLTC_agn_u6173_main",
	"S7Smart",
	"Sabre",
	"TTBooking",
	"VipServiceAmadeus",   //Мы прыгаем на их FTP
	"VipServicePortbilet", //Мы прыгаем на их FTP
}

// Создаётся переменную под структуру Config, в которой будет храниться конфиг
var config Config

func main() {
	log.Println("Программа запущена")

	logFile := createLogFile()
	// Закрываем файл при завершении программы
	defer logFile.Close()

	// Настраиваем MultiWriter для вывода в консоль и в файл
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Устанавливаем логирование в консоль и файл одновременно
	log.SetOutput(multiWriter)

	//Инициализируем конфиг из файла
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %s", err)
	}

	for {
		err := parserSupplier(config.SuppliersFTPConfig, pathToFile)
		if err != nil {
			log.Print(err)
		}

		//Цикл который берёт из массива каждый элемент (название папки)
		for i, folderSupplier := range foldersSuppliers {
			folderSave := foldersSave[i]
			//Запуск копирования файлов в архив и отправку на локальные FTP из папки от поставщиков
			err := processSupplierFiles(folderSupplier, folderSave, config.FTPConfigs)
			if err != nil {
				log.Println(err)
			}
		}
		//Пауза
		second := 10
		//log.Printf("Пауза в %d секунд", second)
		time.Sleep(time.Duration(second) * time.Second)
	}
}

// createLogFile Создаёт папку "logs" в случае её отсутствия и создаёт в ней лог файл формата "", куда будут сохраняться все логи
func createLogFile() *os.File {
	// Получаем текущую дату и время
	now := time.Now()

	// Задаем имя папки для логов
	logDir := "logs"

	// Проверяем, существует ли папка, если нет — создаем
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			log.Fatal("Ошибка при создании папки для логов: ", err)
		}
	}

	// Формируем имя файла с логами
	fileName := fmt.Sprintf("start_%d-%d-%d_%d.%d.%d.log", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())

	// Полный путь до файла
	filePath := filepath.Join(logDir, fileName)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Ошибка при создании файла для записи логов: ", err)
	}

	// Устанавливаем этот файл как место вывода для логов
	log.SetOutput(file)
	return file
}

// loadConfig Функция инициализирующая конфиг из файла с конфигом
func loadConfig(path string) (*Config, error) {
	//Читает файл
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Не удалось прочитать файл конфигурации: %w", err)
	}

	//Парсит Прочитанный файл в config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("Не удалось распарсить YAML: %w", err)
	}

	return &config, nil
}

// parserSupplier Забирает файлы с FTP сервера поставщика
func parserSupplier(suppliersFtpConfigs map[string]SuppliersFTPConfig, pathToFile string) error {
	//В конфиге несколько поставщиков, для этого нужен цикл
	for _, config := range suppliersFtpConfigs {
		switch config.TypeConnect {
		case "ftp":
			err := parserSupplierFTP(config, pathToFile)
			if err != nil {
				return err
			}
		case "sftp":
			err := parserSupplierSFTP(config, pathToFile)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Неверно указан тип подключения к поставщику в конфигурационном фалйе")
		}
	}
	return nil
}

// parserSupplier Забирает файлы с FTP сервера поставщика
func parserSupplierFTP(config SuppliersFTPConfig, pathToFile string) error {
	//Подключаюсь на FTP сервер
	client, err := ftp.Dial(config.Address)
	if err != nil {
		return fmt.Errorf("Ошибка при подключение к FTP серверу поставщика %s.\nОшибка: %s", config.Name, err)
	}
	defer client.Quit()

	//Вводим данные Логин:Пароль
	err = client.Login(config.User, config.Pass)
	if err != nil {
		return fmt.Errorf("Ошибка при вводе логина/пароля на FTP: %s.\nОшибка: %s", config.Name, err)
	}

	//В конфиге может быть несколько папок, которые нужно проверить. Для этого перебираем массив путей
	for _, path := range config.Path {
		//Переходим в нужную нам директорию (папку указанную в конфиге)
		err = client.ChangeDir(path)
		if err != nil {
			return fmt.Errorf("Ошибка при смнене директории на: %s на FTP: %s\nОшибка: %s", config.Name, path, err)
		}

		//Сканируем папку на наличие файлов
		files, _ := client.List("")

		//Проходимся по полученному массиву для скачивания файлов в локальную директорию в случае наличии их на удалённом сервере
		for _, file := range files {
			//Проверка, если это папка, то пропускает её
			if file.Type == ftp.EntryTypeFolder {
				break
			}

			//Переменная с полным путём к удалённому файлу
			filePath := fmt.Sprintf("%s%s", path, file.Name)

			//Получаем копию файла с удалённого сервера
			retr, err := client.Retr(filePath)
			if err != nil {
				return fmt.Errorf("Ошибка при получение даных файла: %s на удалённом FTP сервере: %s\nОшибка: %s", filePath, config.Name, err)
			}
			defer retr.Close()

			//Создание переменной с полным путём к создаваемому локальному файлу
			pathToSupplierFile := fmt.Sprintf("%s\\%s\\%s", pathToFile, config.SaveIn, file.Name)

			//Создание пустого локального файла
			localFile, err := os.Create(pathToSupplierFile)
			if err != nil {
				return fmt.Errorf("Ошибка при создании локального файла: %s полученного с удалённого FTP сервера: %s\nОшибка: %s", pathToSupplierFile, config.Name, err)
			}
			defer localFile.Close()

			//Копирование данных из удалённого файла в локальный файл
			_, err = io.Copy(localFile, retr)
			if err != nil {
				return fmt.Errorf("Ошибка при копирование данных файла с удалённого FTP сервера: %s на локальный\nОшибка: %s", config.Name, err)
			}
			//Удаления файла на FTP поставщика
			_ = client.Delete(filePath)
		}
	}
	return nil
}

// parserSupplier Забирает файлы с FTP сервера поставщика
func parserSupplierSFTP(config SuppliersFTPConfig, pathToFile string) error {
	//Создаём SSH конфиг
	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //Игнорировать проверку ключа хоста
	}

	//Подключаемся на SFTP сервер
	sshClient, err := ssh.Dial("tcp", config.Address, sshConfig)
	if err != nil {
		return fmt.Errorf("Ошибка при подключение на SFTP сервер поставщика %s.\nОшибка: %s", config.Name, err)
	}
	defer sshClient.Close()

	//Подкалючаемя к самой SFTP
	sftpClient, err := sftp.NewClient(sshClient)

	//В конфиге может быть несколько папок, которые нужно проверить. Для этого перебираем массив путей
	for _, path := range config.Path {
		//Сканируем указанную нами директорию, на наличие файлов (папку указанную в конфиге)
		files, err := sftpClient.ReadDir(path)
		if err != nil {
			return fmt.Errorf("Ошибка при сканирование директории на: %s на SFTP: %s\nОшибка: %s", config.Name, path, err)
		}

		//Проходимся по полученному массиву для скачивания файлов в локальную директорию в случае наличии их на удалённом сервере
		for _, file := range files {
			//Переменная с названием файла
			fileName := string(file.Name())
			//Переменная с полным путём к удалённому файлу
			filePath := fmt.Sprintf("%s%s", path, fileName)

			//Получаем копию файла с удалённого сервера
			srcFile, err := sftpClient.Open(filePath)
			if err != nil {
				return fmt.Errorf("Ошибка при получение даных файла: %s на удалённом SFTP сервере: %s\nОшибка: %s", filePath, config.Name, err)
			}
			defer srcFile.Close()

			//Создание переменной с полным путём к создаваемому локальному файлу
			pathToSupplierFile := fmt.Sprintf("%s\\%s\\%s", pathToFile, config.SaveIn, fileName)

			//Создание пустого локального файла
			localFile, err := os.Create(pathToSupplierFile)
			if err != nil {
				return fmt.Errorf("Ошибка при создании локального файла: %s полученного с удалённого FTP сервера: %s\nОшибка: %s", pathToSupplierFile, config.Name, err)
			}
			defer localFile.Close()

			//Копирование данных из удалённого файла в локальный файл
			_, err = io.Copy(localFile, srcFile)
			if err != nil {
				return fmt.Errorf("Ошибка при копирование данных файла с удалённого FTP сервера: %s на локальный\nОшибка: %s", config.Name, err)
			}

			//Удаления файла на SFTP поставщика
			_ = sftpClient.Remove(filePath)
		}
	}
	return nil
}

// processSupplierFiles Функция по копированию файлов в архив и отправкой на FTP сервер
func processSupplierFiles(folderSupplier, folderSave string, ftpConfigs map[string]FTPConfig) error {
	//Создание переменных с путём по поставщикам, куда должны сохраняться файлы
	fullPathToFile := fmt.Sprintf("%s\\%s", pathToFile, folderSupplier)
	fullPathToArchive := fmt.Sprintf("%s\\%s", pathToArchive, folderSave)

	//Читаю папку на наличие файлов, которые получаю от поставщиков
	files, err := os.ReadDir(fullPathToFile)
	if err != nil {
		return fmt.Errorf("Ошибка при сканирование папки: %s на наличие файлов.\nОшибка: %s", folderSupplier, err)
	}

	//Если нашёл файл, то заходим в цикл
	for _, file := range files {
		//Создание переменных с путём и названием файла
		fullPathToFileAndFileName := fmt.Sprintf("%s\\%s", fullPathToFile, file.Name())
		fullPathToArchiveAndFileName := fmt.Sprintf("%s\\%s", fullPathToArchive, file.Name())

		//Копирует файлы в архив
		err := copyFile(fullPathToFileAndFileName, fullPathToArchiveAndFileName)
		if err != nil {
			return fmt.Errorf("Ошибка при копировании файла!\n%s", err)
		}

		//Отправляет файлы на FTP сервера, данные берёт из конфига
		for name, config := range ftpConfigs {
			if err := sendToFTP(config, file.Name(), folderSave, fullPathToFileAndFileName); err != nil {
				return fmt.Errorf("Ошибка при отправке файла на FTP %s: %s", name, err)
			}
		}

		//Удаляем оригинальный файл
		err = os.Remove(fullPathToFileAndFileName)
		if err != nil {
			return fmt.Errorf("Ошибка при удаление файла: %s\nОшибка: %s", fullPathToFileAndFileName, err)
		}
	}
	return nil
}

// copyFile Копирует файл в указанную папку
func copyFile(pathToFileAndFileName, pathToSaveAndFileName string) error {
	//Открываю файл который буду копировать
	sourceFile, err := os.Open(pathToFileAndFileName)
	if err != nil {
		return fmt.Errorf("Ошибка при открытии оригинального файла в директории: %s | для копирования данных.\nОшибка: %s", pathToFileAndFileName, err)
	}
	defer sourceFile.Close()

	//Создаю файл в который будут скопированы данные
	newFile, err := os.Create(pathToSaveAndFileName)
	if err != nil {
		return fmt.Errorf("Ошибка при создании нового пустого файла в директории: %s\nОшибка: %s", pathToSaveAndFileName, err)
	}
	defer newFile.Close()

	//Копираю данные из ориг. файла в новый
	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		return fmt.Errorf("Ошибка при копирование данных из оригинального файла: %s в новый файл : %s\nОшибка: %s", pathToFileAndFileName, pathToSaveAndFileName, err)
	}
	log.Printf("Файл скопирован в архив: %s", pathToSaveAndFileName)
	return nil
}

// sendToFTP ftpHost, ftpPort, ftpUser, ftpPass - креды для подключения к FTP серверу; supplierName Папка поставщика; pathToFileAndFileName Путь к файлу;
func sendToFTP(config FTPConfig, fileName, supplierName, pathToFileAndFileName string) error {
	//Подключаюсь на FTP сервер
	client, err := ftp.Dial(config.Address)
	if err != nil {
		return fmt.Errorf("Ошибка при подключение к FTP серверу.\nОшибка: %s", err)
	}
	defer client.Quit()

	//Вводим данные Логин:Пароль
	err = client.Login(config.User, config.Pass)
	if err != nil {
		return fmt.Errorf("Ошибка при вводе логина/пароля.\nОшибка: %s", err)
	}

	//Переходим в нужную нам директорию (папку с названием поставщика)
	err = client.ChangeDir(supplierName)
	if err != nil {
		return fmt.Errorf("Ошибка при смнене директории на: %s\nОшибка: %s", supplierName, err)
	}

	//Открываем файл который будет скопирован
	sourceFile, err := os.Open(pathToFileAndFileName)
	if err != nil {
		return fmt.Errorf("Ошибка при открытии оригинального файла в директории: %s | для копирования данных.\nОшибка: %s", pathToFileAndFileName, err)
	}
	defer sourceFile.Close()

	//Создаём на удалённом сервере файл с данными
	err = client.Stor(fileName, sourceFile)
	if err != nil {
		return fmt.Errorf("Ошибка при создании и переносе данных на удалённый сервер.\nОшибка: %s", err)
	}

	log.Printf("Файл: %s успешно передан на FTP: %s в директорию %s", fileName, config.Name, supplierName)
	return nil
}
