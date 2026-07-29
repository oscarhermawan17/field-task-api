go run ./cmd/api ==> hasil build akan berada di tempat lain, agar compilasi berikutnya lebih cepat.

    go build ./cmd/api
        hasil build berada di folder root.

    go build -o bin/field-task-api.exe ./cmd/api
        -o berarti output. menentukan sendiri nama file dan lokasi hasil build: