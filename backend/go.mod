module github.com/diakovliev/mesap/backend

go 1.18

require (
	github.com/diakovliev/mesap/backend/controllers v0.0.1 // indirect
	github.com/diakovliev/mesap/backend/fake_database v0.0.1 // indirect
	github.com/diakovliev/mesap/backend/ifaces v0.0.1 // indirect
	github.com/diakovliev/mesap/backend/models v0.0.1 // indirect
	github.com/go-chi/chi v1.5.4 // indirect
	github.com/go-chi/chi/v5 v5.0.7 // indirect
)

replace github.com/diakovliev/mesap/backend/models v0.0.1 => ./models

replace github.com/diakovliev/mesap/backend/ifaces v0.0.1 => ./ifaces

replace github.com/diakovliev/mesap/backend/fake_database v0.0.1 => ./fake_database

replace github.com/diakovliev/mesap/backend/controllers v0.0.1 => ./controllers
