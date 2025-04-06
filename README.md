Utworzenie tokenu klasycznego  w github z uprawnieniami 
avatar -> settings -> developer settings -> personal acces tokens -> classic
Przejście do WSL, generowanie kluczy ssh
crimson@DESKTOP-PIV2R31:~$ ssh-keygen -t ed25519 -C "gabriel.piatek.biznes@gmail.com" -f ~/.ssh/ssh_wsl
crimson@DESKTOP-PIV2R31:~$ eval "$(ssh-agent -s)"
Agent pid 876
crimson@DESKTOP-PIV2R31:~$
crimson@DESKTOP-PIV2R31:~$ ssh-add ~/.ssh/ssh_wsl
Identity added: /home/crimson/.ssh/ssh_wsl (gabriel.piatek.biznes@gmail.com)
Dodanie klucza ssh do uwierzytelnienia na github 
avatar -> settings -> ssh and gpg keys -> new ssh key
sprawdzenie tokenu 
crimson@DESKTOP-PIV2R31:~$ ssh -T git@github.com
Hi CrimsonGabriel! You've successfully authenticated, but GitHub does not provide shell access.
Utworzenie repo w katalogu wykonania zadania z lab5
crimson@DESKTOP-PIV2R31:~$ cd ~/lab5
crimson@DESKTOP-PIV2R31:~/lab5$ git init
crimson@DESKTOP-PIV2R31:~/lab5$ git remote add origin git@github.com:CrimsonGabriel/pawcho6.git

Dodanie plików dodatkowych
touch .hidden
->tutaj klucz ssh :)
touch .gitignore
.hidden.txt
*.pem
*.key.hidden.txt

Modyfikacja Dockerfile Dodać buildkita i argumentów do budowania z secretami
# syntax=docker/dockerfile:1.6

# === Etap 1: Budowanie aplikacji Go ===
FROM golang:alpine AS builder

WORKDIR /app

COPY server.go ./
COPY index.html ./

RUN go mod init myapp

ARG VERSION=1.0.0
ENV VERSION=${VERSION}

# Wykorzystanie sekretów z BuildKit
RUN --mount=type=secret,id=mysecret \
    go build -o aprogram -ldflags "-X main.Version=${VERSION}"

# === Etap 2: Konfiguracja Nginx jako reverse proxy ===
FROM nginx:alpine

# Usunięcie domyślnej konfiguracji Nginx
RUN rm /etc/nginx/conf.d/default.conf

# Kopiowanie własnej konfiguracji Nginx
COPY nginx.conf /etc/nginx/conf.d/

# Kopiowanie aplikacji Go do kontenera
COPY --from=builder /app/aprogram /usr/bin/aprogram
COPY --from=builder /app/index.html /usr/share/nginx/html/

# Ustawienie uprawnień do uruchamiania aplikacji
RUN chmod +x /usr/bin/aprogram

# Otwieramy port 80 dla Nginx
EXPOSE 80

# Healthcheck - sprawdzanie działania serwera
HEALTHCHECK --interval=10s --timeout=3s \
  CMD curl -f http://localhost/ || exit 1

# Uruchomienie zarówno aplikacji Go, jak i Nginx
CMD ["/bin/sh", "-c", "/usr/bin/aprogram & nginx -g 'daemon off;'"]

Dodanie plików do repo 
crimson@DESKTOP-PIV2R31:~/lab5$ git add .
jakiś commicik 
crimson@DESKTOP-PIV2R31:~/lab5$ git push origin main

Zbudowanie obrazu

crimson@DESKTOP-PIV2R31:~/lab5$ DOCKER_BUILDKIT=1 docker build --secret id=mysecret,src=.hidden.txt -t lab6 .
[+] Building 14.3s (23/23) FINISHED                                                                                                                                                                                           docker:default
 => [internal] load build definition from Dockerfile                                                                                                                                                                                    0.0s
 => => transferring dockerfile: 1.16kB                                                                                                                                                                                                  0.0s
 => resolve image config for docker-image://docker.io/docker/dockerfile:1.6                                                                                                                                                             1.4s
 => [auth] docker/dockerfile:pull token for registry-1.docker.io                                                                                                                                                                        0.0s
 => CACHED docker-image://docker.io/docker/dockerfile:1.6@sha256:ac85f380a63b13dfcefa89046420e1781752bab202122f8f50032edf31be0021                                                                                                       0.0s
 => => resolve docker.io/docker/dockerfile:1.6@sha256:ac85f380a63b13dfcefa89046420e1781752bab202122f8f50032edf31be0021                                                                                                                  0.0s
 => [internal] load metadata for docker.io/library/golang:alpine                                                                                                                                                                        0.8s
 => [internal] load metadata for docker.io/library/nginx:alpine                                                                                                                                                                         0.7s
 => [auth] library/nginx:pull token for registry-1.docker.io                                                                                                                                                                            0.0s
 => [auth] library/golang:pull token for registry-1.docker.io                                                                                                                                                                           0.0s
 => [internal] load .dockerignore                                                                                                                                                                                                       0.0s
 => => transferring context: 2B                                                                                                                                                                                                         0.0s
 => CACHED [builder 1/6] FROM docker.io/library/golang:alpine@sha256:7772cb5322baa875edd74705556d08f0eeca7b9c4b5367754ce3f2f00041ccee                                                                                                   0.0s
 => => resolve docker.io/library/golang:alpine@sha256:7772cb5322baa875edd74705556d08f0eeca7b9c4b5367754ce3f2f00041ccee                                                                                                                  0.0s
 => [stage-1 1/6] FROM docker.io/library/nginx:alpine@sha256:4ff102c5d78d254a6f0da062b3cf39eaf07f01eec0927fd21e219d0af8bc0591                                                                                                           0.0s
 => => resolve docker.io/library/nginx:alpine@sha256:4ff102c5d78d254a6f0da062b3cf39eaf07f01eec0927fd21e219d0af8bc0591                                                                                                                   0.0s
 => [internal] load build context                                                                                                                                                                                                       0.0s
 => => transferring context: 2.07kB                                                                                                                                                                                                     0.0s
 => [builder 2/6] WORKDIR /app                                                                                                                                                                                                          0.0s
 => [builder 3/6] COPY server.go ./                                                                                                                                                                                                     0.0s
 => [builder 4/6] COPY index.html ./                                                                                                                                                                                                    0.0s
 => [builder 5/6] RUN go mod init myapp                                                                                                                                                                                                 0.4s
 => [builder 6/6] RUN --mount=type=secret,id=mysecret     go build -o aprogram -ldflags "-X main.Version=1.0.0"                                                                                                                        11.0s
 => CACHED [stage-1 2/6] RUN rm /etc/nginx/conf.d/default.conf                                                                                                                                                                          0.0s
 => CACHED [stage-1 3/6] COPY nginx.conf /etc/nginx/conf.d/                                                                                                                                                                             0.0s
 => CACHED [stage-1 4/6] COPY --from=builder /app/aprogram /usr/bin/aprogram                                                                                                                                                            0.0s
 => CACHED [stage-1 5/6] COPY --from=builder /app/index.html /usr/share/nginx/html/                                                                                                                                                     0.0s
 => CACHED [stage-1 6/6] RUN chmod +x /usr/bin/aprogram                                                                                                                                                                                 0.0s
 => exporting to image                                                                                                                                                                                                                  0.1s
 => => exporting layers                                                                                                                                                                                                                 0.0s
 => => exporting manifest sha256:3633455fa57b3d6aef2b957ec20b3b97c33c6551dcebb9cb8a7b53dcc32bb009                                                                                                                                       0.0s
 => => exporting config sha256:cc58eaad62565af389cbca9e28a8887367c86343fc194c08a3746403fa00e489                                                                                                                                         0.0s
 => => exporting attestation manifest sha256:5a717ffdc1b21826946e30a223066816f59ea3b3a06ef6fdfcffe499cda8e611                                                                                                                           0.0s
 => => exporting manifest list sha256:71eb58d2330bcdf96deba289699c9b0fc08c94967e48bff96a39e4aa848abccd                                                                                                                                  0.0s
 => => naming to docker.io/library/lab6:latest                                                                                                                                                                                          0.0s
 => => unpacking to docker.io/library/lab6:latest    

 Otagowanie zbudowanego obrazu i spushowanie go jako paczkę do gita
crimson@DESKTOP-PIV2R31:~/lab5$ docker tag lab6 ghcr.io/crimsongabriel/pawcho6:lab6
crimson@DESKTOP-PIV2R31:~/lab5$ docker push ghcr.io/crimsongabriel/pawcho6:lab6
The push refers to repository [ghcr.io/crimsongabriel/pawcho6]
617894c4e5c8: Pushed
f18232174bc9: Mounted from crimsongabriel/myapp_nginx
144c6700c3a3: Mounted from crimsongabriel/myapp_nginx
24bca82fd4ac: Mounted from crimsongabriel/myapp_nginx
6d79cc6084d4: Mounted from crimsongabriel/myapp_nginx
ccc35e35d420: Mounted from crimsongabriel/myapp_nginx
0c7e4c092ab7: Mounted from crimsongabriel/myapp_nginx
8d27c072a58f: Mounted from crimsongabriel/myapp_nginx
4f4fb700ef54: Mounted from crimsongabriel/myapp_nginx
ab3286a73463: Mounted from crimsongabriel/myapp_nginx
984583bcf083: Mounted from crimsongabriel/myapp_nginx
43f2ec460bdf: Mounted from crimsongabriel/myapp_nginx
2856cad21856: Mounted from crimsongabriel/myapp_nginx
ed1b8c11cc1d: Mounted from crimsongabriel/myapp_nginx
lab6: digest: sha256:71eb58d2330bcdf96deba289699c9b0fc08c94967e48bff96a39e4aa848abccd size: 856

Sprawdzanie
![image](https://github.com/user-attachments/assets/3d81b7a6-4a59-4330-9330-2d245dafc20b)
![image](https://github.com/user-attachments/assets/62aa8b25-7aba-4dbe-b6bb-67b3e1ee34cc)
