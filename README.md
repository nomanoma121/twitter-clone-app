# Twitter Clone App

## 技術スタック

- [Go](https://golang.org/)
- [React](https://reactjs.org/)
- [Docker](https://www.docker.com/)
- [MySQL](https://www.mysql.com/)
- [Nginx](https://www.nginx.com/)
- [JWT](https://jwt.io/)

## 開発方法

Client側の依存関係をインストール (ローカルで開発するための型定義参照のため)

```bash
cd client
npm install
```

Server側の依存関係をインストール

```bash
cd server
go mod download
```

`.env`ファイルを作成

```bash
cp .env.example .env
```

適宜`.env`ファイルを編集してください。

Dockerコンテナを起動

```bash
docker-compose -f docker-compose.dev.yml up --build
```

### Server側の更新反映

Server側は、`server`ディレクトリがボリュームとしてマウントされているため、Server側のファイルを更新してコンテナを再起動するだけで更新が反映されます。

```bash
docker-compose -f docker-compose.dev.yml restart
```

ただし、マイグレーションファイルを更新した場合などDBのスキーマが更新される場合は、ボリュームの再作成を推奨します。
(この時DBのデータは全て消えるので注意してください。開発時はseederを作るなどしてダミーデータを生成すると良いでしょう)

```bash
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up --build
```

### Client側の更新反映

基本的には、Viteのホットリロードが有効になっているため、Client側のファイルを更新するだけで更新が反映されます。
ライブラリの追加などでnode_modulesの更新が予期される場合は、virtual volumeの`node_modules`を削除する必要があります。
(この時DBのデータは全て消えるので注意してください。開発時はseederを作るなどしてダミーデータを生成すると良いでしょう)

```bash
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up --build
```

これでDockerコンテナ内の`node_modules`が更新されます。

## デプロイ

1. クラウドで借りれるEC2やGCPのCompute EngineなどのVPSサービスを契約します。
2. Dockerをインストールします。
3. このリポジトリをクローンします。
4. `.env`ファイルを作成します。(`.env.example`をコピーして、中身を編集してください)
5. `docker compose up -d --build`を実行します。

これで、Webアプリケーションが`:80`ポートで公開されます。
