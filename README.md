# Product Service @ Erajaya

TECH STACK : GO, POSTGRESQL, REDIS 

## Dibuat berdasarkan prinsip Clean Architecture 

Penerapan Clean Architecture pada aplikasi ini betujuan untuk menghasilkan code tertata, rapih, dan sesuai dengan scope masing masing. Saya membaginya menjadi 4 Layer, yaitu

    1. Model / Entity

        Pendefinisian entitas yang ada dalam aplikasi

    2. Repository

        Bagian dimana aplikasi melakukan operasi ke db, redis atau yang lainya.

    3. Service / Usecase

        Tempat dimana Business Logic ditulis
    
    4. Handler / Controller 

        Pengecekan request yang masuk dan handling yang akan di lakukan


## Cara menjalankan aplikasi

    1. Buat sebuah database postgresql baru yang bernama "product_service"

    2. Ubah credensial database dan redis pada file .env, gunakan "localhost" jika service db / redis berjalan pada sistem lokal dan maka jalankan aplikasi dengan menggunakan perintah "make run dev", jika anda menjalankan aplikasi dengan "docker compose up" dan service db / redis anda berjalan di lokal sistem juga maka ganti "localhost" dengan "host.docker.internal"

    3. pastikan package "go-migare" telah terinstall di komputer anda , jalankan "make migrate up"

    4. Jika migration berhasil, maka anda bisa menjalankan aplikasi dengan "make run dev" atau "docker compose up" yang akan langsung membuild dan menjalankan aplikasi diatas docker.

    5. Port yang digunakan oleh aplikasi ini adalah port 8088