<div align="center">
    <h1>Pasta 🍝</h1>
    <h3>Copy-paste files between devices</h3>
</div>

## Usage
Start the server by running:

![example-pasta](https://user-images.githubusercontent.com/19193227/221366059-62c122b2-1e78-4f3b-93f2-6657ed46dcc1.png)

Then, browse to the link on the other device (or scan the QR code to do that), and upload the files you want.

For more information:
```sh
pasta --help
```

## Installation
From source:
```sh
git clone https://github.com/idan22moral/pasta
cd pasta
go install
```
From release:
```sh
wget https://github.com/idan22moral/pasta/releases/download/v1.0.0/pasta-v1.0.0-linux.tar.gz
tar -xzvf pasta-v1.0.0-linux.tar.gz
mv linux/pasta /usr/local/bin/pasta
```
