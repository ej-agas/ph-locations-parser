# PH Locations Parser

[![Star](https://img.shields.io/github/stars/ej-agas/ph-locations-parser.svg?style=flat-square)](https://github.com/ej-agas/ph-locations-parser/stargazers) [![License](https://img.shields.io/github/license/ej-agas/ph-locations-parser.svg?style=flat-square)](https://github.com/ej-agas/ph-locations/blob/main/LICENSE) [![Release](https://img.shields.io/github/release/ej-agas/ph-locations-parser.svg?style=flat-square)](https://github.com/ej-agas/ph-locations-parser/releases)

CLI application that parses Philippine Standard Geographic Code (PSGC) publication and stores it to ph-locations' database

<p>
    <img src="https://github.com/ej-agas/ph-locations-parser/blob/main/assets/ph-locations-parser.gif" width="100%" alt="PH Locations Parser Example">
</p>


## 📖 Usage

### Help
```shell
ph-locations-parser help
```

### Parsing the publication file
```shell
ph-locations-parser parse file.xlsx --host localhost --port 5432 --db foo_db --user ph_locations_user --password
```
### Flags
**--host** database host (default value: 127.0.0.1) \
**--port** database port (default value: 5432) \
**--db** database name (default value: ph_locations_db) \
**--user** database user (default value: ph_locations_user) \
**--password** database password (will prompt for password input)

