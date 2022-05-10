module github.com/diakovliev/mesap/backend/controllers

go 1.18

require (
	github.com/diakovliev/mesap/backend/fake_database v0.0.1 // indirect
	github.com/diakovliev/mesap/backend/ifaces v0.0.1 // indirect
	github.com/diakovliev/mesap/backend/models v0.0.1 // indirect
	github.com/go-chi/chi/v5 v5.0.7 // indirect
	github.com/kong/go-srp v0.0.0-20191210190804-cde1efa3c083 // indirect
	golang.org/x/crypto v0.0.0-20200109152110-61a87790db17 // indirect
	golang.org/x/sys v0.0.0-20190412213103-97732733099d // indirect
)

replace github.com/diakovliev/mesap/backend/models v0.0.1 => ../models

replace github.com/diakovliev/mesap/backend/ifaces v0.0.1 => ../ifaces

replace github.com/diakovliev/mesap/backend/fake_database v0.0.1 => ../fake_database
