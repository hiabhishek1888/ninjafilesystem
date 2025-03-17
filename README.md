# Peer-to-Peer Distributed File Storage System
## Overview
This project is a **Peer-to-Peer (P2P) Distributed File Storage System**, encrypted, decentralized file sharing(store, send, recieve and delete) among connected remote nodes (peers). 
It is built with **Go**, using its various features like **Concurrency**, **Goroutines**, **channels**, **mutex/locks**, **waitgroups**, **gob encoding** and **TCP connections** for efficient data transmission.
The system architecture is a basic **Pub-Sub & Message Queue** pattern, where a device can broadcast the data to its multiple connected peer (remote nodes) and consume data to read/store the data independently in real-time.

## Features
- **Decentralized Storage**: Files are securly stored across multiple peers without a central authority, but the file content is encrypted for connected peer and it can be only read by the file owner only.
- **Secure File Transfer**: Client side encryption using Uses AES for data security, during transit and at rest as well.
- **GOB Serialization**: Efficient data encoding and decoding.
- **Concurrency Support**: Utilizes Goroutines and Channels for real-time message passing.
- **Fault Tolerance**: Ensures data availability even if some peers disconnect.

## Features to be implemented:
- **Cli support**
- **Peer Discovery**: Dynamically connects to other peers in the network.
- **Snapshot (Backup)**: Current implementation overwrite the modifications. Needs to introduce versioning to support backups and data recovery. Since data can be large, we will implement a way to store only the changes in newer versions
- **Chunking**: Multipart Upload for large files
- **Fully Secure Data transfer:**: Current client side encryption is not suitable for large file encryption and do not prevent MITM attacks. Chunking and TLS needs to be implemented.

## Architecture
- **Transport Layer**: Handles TCP connections and peers.
- **Server**: Handles all kind of operation for a peer i.e. StartConnection, bootstrapConnectionWithPeers, StoreData, GetData, encryption, and few Goroutines.
- **Storage**: Handles the disk based operations to read, store and delete files over OS.
- **Message Queue**: Uses channels to manage data flow to Goroutines.




*This project is built by taking help from youtube and internet and it is not exactly copied. Also, there are multiple issues which needs to be fixed, currently working on few, for more you can refer todo.txt file*
