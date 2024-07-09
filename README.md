# Digital Signing with GoLang

## Installation Go

Make sure Go is installed on your computer before running this project.
To install Go, visit the [official Go website](https://golang.org/) and follow the installation instructions.

## Running a Project

Here are the steps to get your Go project up and running:

#### 1. Clone Repository

```bash
git clone https://github.com/hizidha/digital-signing-go.git
cd digital-signing-go-master
```

#### 2. Install Dependencies
Initialize Go modules and tidy up dependencies:
```bash
go mod init digital-signing-go
go mod tidy
```

#### 3. Enable MongoDB Database
Ensure you have MongoDBCompass installed and running. To install PostgreSQL, visit the [official MongoDB website](https://www.mongodb.com/docs/manual/installation/) and follow the installation instructions. Commanly was running at [localhost:27017](http://localhost:27017).

#### 4. Complete All Data Requirements in ``.env``
Complete all variables according to your preferences.
```bash
PORT=your_port
HOSTNAME=your_hostname
MONGO_URI=your_mongodb
MONGO_DB=your_database_name
MONGO_COLLECTION=your_collection_name
```

#### 5. Execute the Project
Run the following command to start the server:
```bash
go run .
```

#### 6. Access the Link Shortener
Open your web browser and navigate to:
```bash
http://localhost:{your_port}/
```